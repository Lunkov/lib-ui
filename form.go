package ui

import (
  "bytes"
  "strconv"
  "sort"
  "encoding/json"
  "os"
  "io"
  "fmt"
  "net/http"
  "mime/multipart"
  
  "io/ioutil"
  "html/template"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-tr"
)

// Form methods: POST or GET.
const (
  POST = "POST"
  GET  = "GET"
)

type InfoItem struct {
  Description  string      `json:"description"   yaml:"description"`
  Enum         []string    `json:"enum"          yaml:"enum"`
}

type InfoProperty struct {
  CODE           string       `json:"code"          yaml:"code"`
  Order          int          `json:"order"         yaml:"order"`
  Title          string       `json:"title"         yaml:"title"`
  Description    string       `json:"description"   yaml:"description"`
  DataType       string       `json:"data_type"     yaml:"data_type"`
  WidgetType     string       `json:"widget"        yaml:"widget"`
  Required       string       `json:"required"      yaml:"required"`
  Format         string       `json:"format"        yaml:"format"`
  Table          string       `json:"table"         yaml:"table"`
  Default        string       `json:"default"       yaml:"default"`
  OneOf          []InfoItem   `json:"oneOf"         yaml:"oneOf"`
  OneOfConst     string       `json:"oneOfConst"    yaml:"oneOfConst"`
}

type InfoButton struct {
  CODE           string       `json:"code"          yaml:"code"`
  Order          int          `json:"order"         yaml:"order"`
  Title          string       `json:"title"         yaml:"title"`
  Type           string       `json:"type"          yaml:"type"`
}

type InfoTab struct {
  CODE         string    `json:"code"         yaml:"code"`
  Order        int       `json:"order"        yaml:"order"`
  Title        string    `json:"title"        yaml:"title"`
  Description  string    `json:"description"  yaml:"description"`
  Fields       []string  `json:"fields"       yaml:"fields"`
}

type ArProperties []InfoProperty
type ArButtons    []InfoButton
type ArTabs       []InfoTab

type FormInfo struct {
  Code          string                     `json:"code"         yaml:"code"`
  EncType       string                     `json:"enctype"      yaml:"enctype"`
  Method        string                     `json:"method"       yaml:"method"`
  Action        string                     `json:"action"       yaml:"action"`
  Version       int32                      `json:"version"      yaml:"version"`
  Date          string                     `json:"date"         yaml:"date"`
  Date_start    string                     `json:"date_start"   yaml:"date_start"`
  Date_finish   string                     `json:"date_finish"  yaml:"date_finish"`
  Title         string                     `json:"title"        yaml:"title"`
  Type          string                     `json:"type"         yaml:"type"`
  Description   string                     `json:"description"  yaml:"description"`
  Tabs          ArTabs                     `json:"tabs"         yaml:"tabs,flow"`
  Properties    ArProperties               `json:"properties"   yaml:"properties,flow"`
  Buttons       ArButtons                  `json:"buttons"      yaml:"buttons,flow"`
}

var mapForms = make(map[string]FormInfo)

func formCount() int {
  return len(mapForms)
}

func formNew(info FormInfo) string {
  info.TrDef()
  if info.Method == "" {
    info.Method = "POST"
  }
  if info.EncType == "" {
    info.EncType = "multipart/form-data"
  }
  if glog.V(9) {
    glog.Infof("LOG: Load Form: '%s'", info.Code)
  }
  // Sorting by Order
  sort.Slice(info.Tabs, func(i, j int) bool { return info.Tabs[i].Order < info.Tabs[j].Order })
  sort.Slice(info.Properties, func(i, j int) bool { return info.Properties[i].Order < info.Properties[j].Order })
  sort.Slice(info.Buttons, func(i, j int) bool { return info.Buttons[i].Order < info.Buttons[j].Order })

  mapForms[info.Code] = info
  return info.Code
}

func formGet(code string) *FormInfo {
  i, ok := mapForms[code]
  if ok {
    return &i
  }
  return nil
}

func (info *FormInfo) TrDef() {
  tr.SetDef(info.Title)
  tr.SetDef(info.Description)
  for _, tab := range info.Tabs {
    tr.SetDef(tab.Title)
    tr.SetDef(tab.Description)
  }
  for _, prop := range info.Properties {
    tr.SetDef(prop.Title)
    tr.SetDef(prop.Description)
  }
  for _, button := range info.Buttons {
    tr.SetDef(button.Title)
  }
}

func (info *FormInfo) Tr(lang string) {
  info.Title, _ = tr.Tr(lang, info.Title)
  info.Description, _ = tr.Tr(lang, info.Description)
  for i, tab := range info.Tabs {
    tab.Title, _ = tr.Tr(lang, tab.Title)
    tab.Description, _ = tr.Tr(lang, tab.Description)
    info.Tabs[i] = tab
  }
  for i, prop := range info.Properties {
    prop.Title, _ = tr.Tr(lang, prop.Title)
    prop.Description, _ = tr.Tr(lang, prop.Description)
    info.Properties[i] = prop
  }
  for i, button := range info.Buttons {
    button.Title, _ = tr.Tr(lang, button.Title)
    info.Buttons[i] = button
  }
}

func (info *FormInfo) Copy(in *FormInfo) {
  buf, _ := json.Marshal(in)
  json.Unmarshal(buf, info)
}

func formFilter(lang string, forms []string) map[string]FormInfo {
  res := make(map[string]FormInfo)
  for _, form_code := range forms {
    m := formGet(form_code)
    if m != nil {
      m.Tr(lang)
      res[form_code] = (*m)
    }
  }
  return res
}

func formInit(configPath string) {
  mapForms = make(map[string]FormInfo)
  memTemplate = make(map[string]*template.Template)
  env.LoadFromYMLFiles(configPath, loadFormYAML)
}

func FormGetParameters(r *http.Request, storagePath string) (map[string]interface{}, bool) {
  params := make(map[string]interface{})
  /* TODO
   * code_form string, 
  form := formGet(code_form)
  if form == nil {
    glog.Errorf("ERR: Form '%s' not found", code_form)
    return params, false
  }*/
  
  reader, err := r.MultipartReader()

  if err != nil {
    glog.Errorf("ERR: MultipartReader URL '%s': `%v`", r.URL.Path, err.Error())
    return params, false
  }
  for {
    part, err := reader.NextPart()
    if err == io.EOF {
      break
    }

    if part.FileName() == "" {
      name := part.FormName()
      value := readValueFromFormData(part)
      params[name] = value
      continue
    }

    if glog.V(9) {
      glog.Infof("DBG: UploadFile URL '%s': `%v`", r.URL.Path, part.FileName())
    }
    
    filename := storagePath + part.FileName()
    dst, err := os.Create(filename)
    defer dst.Close()

    if err != nil {
      glog.Errorf("ERR: Save UploadFile URL '%s': `%v` %s", r.URL.Path, part.FileName(), err.Error())
      return params, false
    }
        
    if _, err := io.Copy(dst, part); err != nil {
      glog.Errorf("ERR: Write UploadFile URL '%s': `%v` %s", r.URL.Path, part.FileName(), err.Error())
      return params, false
    }
    
    params[part.FormName()] = filename

  }
  return params, true
}

func readValueFromFormData(part *multipart.Part) string {
  buf := new(bytes.Buffer)
  buf.ReadFrom(part)
  return buf.String()
}

func readFileParam(part *multipart.Part) ([]byte, bool) {
  content := []byte{}
  uploaded := false
  var err error

  for !uploaded {
    chunk := make([]byte, 4096)
    if _, err = part.Read(chunk); err != nil {
      if err != io.EOF {
        glog.Errorf("ERR: error reading multipart file %v", err.Error())
        return nil, false
      }
      uploaded = true
    }
    content = append(content, chunk...)
  }

  return content, true
}

func formExpandProperties(info *FormInfo) {
      /*
  for i, prop := range info.Properties {
    if prop.OneOfConst != "" {
      // TODO
      c := constant.Get(prop.OneOfConst)
      if c != nil {
        for _, item := range c.Table {
          ip := InfoItem{Description: item.Name}
          ip.Enum = append(ip.Enum, strconv.Itoa(item.Id))
          // log.Printf("FORM: expandProperties(%s)=%d ", item.Name, item.Id)
          prop.OneOf = append(prop.OneOf, ip)
        }
        info.Properties[i] = prop
      }
    }
  }*/
}

func loadFormYAML(filename string, yamlFile []byte) int {
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return 0
  }
  var mapTmp FormInfo

  err = yaml.Unmarshal(yamlFile, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  formExpandProperties(&mapTmp)
  formNew(mapTmp)

  return 1
}

func indexForm(lang string, form_code string, style string, isModal bool, data *map[string]interface{}) string {
  return "FORM#" + lang + "#" + form_code + "#" + style + "#" + strconv.FormatBool(isModal) + "#" + strconv.FormatBool(data == nil)
}

func renderForm(lang string, form_code string, style string, isModal bool, data *map[string]interface{}) string {
  f := formGet(form_code)
  if f == nil {
    err_str := fmt.Sprintf("ERR: FORM(%s) NOT FOUND ", form_code)
    glog.Errorf(err_str)
    return err_str
  }
  var form FormInfo
  form.Copy(f)
  form.Tr(lang)
  typeForm := ""
  if form.Type != "" {
    typeForm = form.Type
  }
  if isModal {
    typeForm = "modal"
  }
  if typeForm == "" {
    typeForm = "form"
  }

	tmpl := getTemplate(typeForm, "forms", style, lang)
	if tmpl == nil {
    err_str := fmt.Sprintf("ERR: templateFile(%s:%s).Load ", style, form_code)
    glog.Errorf(err_str)
    return err_str
	}
  var tpl bytes.Buffer
  
  body := ""
  footer := ""
  
  mapProps := make(map[string]string)
  for _, prop := range form.Properties {
    tmplProp := getTemplate("input_" + prop.WidgetType, "forms", style, lang)
    if tmplProp == nil {
      err_str := fmt.Sprintf("ERR: templateFile(%s:%s): Field(%s): not found widget '%s'", style, form_code, prop.CODE, prop.WidgetType)
      glog.Errorf(err_str)
      mapProps[prop.CODE] = "<span>" + err_str + "</span>"
      continue
    }
    value := ""
    if data != nil {
      v, ok := (*data)[prop.CODE]
      if ok {
        str, okstr := value2String(v)
        if okstr {
          value = str
        }
      }
    }
    propP := map[string]string{
      "FORM_ID": form.Code,
      "CODE": prop.CODE,
      "TITLE": prop.Title, // tr.Tr(lang, prop.Title),
      "VALUE": value,
      "DESCRIPTION": prop.Description, // tr.Tr(lang, prop.Description),
      "LANG": lang,
    }
    tpl.Reset()
    err := tmplProp.Execute(&tpl, propP)
    if err == nil {
      if len(form.Tabs) > 0 {
        mapProps[prop.CODE] = tpl.String()
      } else {
        body = body + tpl.String()
      }
    }
  }

  if len(form.Tabs) > 0 {
    tTab := "tabs"
    if data == nil {
      tTab = "stepper"
    }
    tabTitle := ""
    tabContent := ""
    
    tmplTabTitle   := getTemplate("tab_title",   "forms", style, lang)
    tmplTabContent := getTemplate("tab_content", "forms", style, lang)
    
    tmplTab := getTemplate(tTab, "forms", style, lang)
    for i, tab := range form.Tabs {
      body = ""
      for _, field := range tab.Fields {
        // TODO: Need sort by Order
        item, ok := mapProps[field]
        if ok {
          body = body + item
        } else {
          glog.Errorf("ERR: templateFile(%s:%s).Tab(%d): not found field '%s'", style, form_code, i, field)
        }
      }
      // RENDER TITLES
      active := ""
      if i == 0 {
        active = "active"
      }
      propT := map[string]interface{} {
        "FORM_ID": form.Code,
        "CODE":   tab.CODE,
        "ACTIVE": active,
        "TITLE":  tab.Title, //tr.Tr(lang, tab.Title),
        "LANG": lang,
      }
      tpl.Reset()
      err := tmplTabTitle.Execute(&tpl, propT)
      if err == nil {
        tabTitle += tpl.String()
      }
      // RENDER BODIES
      active = ""
      if i == 0 {
        active = "active show"
      }
      propT = map[string]interface{} {
        "FORM_ID": form.Code,
        "CODE": tab.CODE,
        "ACTIVE": active,
        "TITLE": tab.Title, // tr.Tr(lang, tab.Title),
        "CONTENT": template.HTML(body),
        "LANG": lang,
      }
      tpl.Reset()
      err = tmplTabContent.Execute(&tpl, propT)
      if err == nil {
        tabContent += tpl.String()
      }

    }
    
    propT := map[string]interface{} {
      "FORM_ID": form.Code,
      "TABS_TITLE": template.HTML(tabTitle),
      "TABS_CONTENT": template.HTML(tabContent),
      "LANG": lang,
    }
    tpl.Reset()
    err := tmplTab.Execute(&tpl, propT)
    if err == nil {
      body = body + tpl.String()
    }
  }
  
  for _, button := range form.Buttons {
    tmplBtn := getTemplate("button_" + button.Type, "forms", style, lang)
    if tmplBtn != nil {
      tpl.Reset()
      propB := map[string]interface{} {
        "FORM_ID": form.Code,
        "ID":    button.CODE,
        "TITLE": button.Title, //tr.Tr(lang, button.Title),
        "LANG": lang,
      }
      err := tmplBtn.Execute(&tpl, propB)
      if err == nil {
        footer = footer + tpl.String()
      }
    } else {
      err_str := fmt.Sprintf("ERR: templateFile(%s:%s): not found button '%s'", style, form_code, button.CODE)
      glog.Errorf(err_str)
      footer = footer + "<span>" + err_str + "</span>"
    }
  }
  
  propForm := map[string]interface{} {
    "FORM_ID":      form.Code,
    "FORM_METHOD":  form.Method,
    "FORM_ACTION":  form.Action,
    "FORM_ENCTYPE": form.EncType,
    "FORM_TITLE":        form.Title,
    "FORM_DESCRIPTION":  form.Description,
    "FORM_BODY":   template.HTML(body),
    "FORM_FOOTER": template.HTML(footer),
    "LANG": lang,
  }
  unionMap(&propForm, data)
  tpl.Reset()
  err := tmpl.Execute(&tpl, propForm)
  if err != nil {
    err_str := fmt.Sprintf("ERR: templateFile(%s:%s).Execute: '%v'", style, form_code, err)
    glog.Errorf(err_str)
    return err_str
  }
  return tpl.String()
}

