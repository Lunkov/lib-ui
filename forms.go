package ui

import (
  "bytes"
  "strings"
  "strconv"
  "sort"
  "os"
  "io"
  "fmt"
  "net/http"
  "mime/multipart"
  
  "html/template"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-tr"
  "github.com/Lunkov/lib-ref"
  "github.com/Lunkov/lib-cache"
)


type Forms struct {
  t               *tr.Tr
  m                cache.ICache
  templates       *Templates
}

func NewForms(cfg *cache.CacheConfig, t *tr.Tr, ts *Templates) *Forms {
  return &Forms{ m: cache.NewConfig(cfg), t: t, templates: ts }
}

func (f *Forms) Count() int64 {
  return f.m.Count()
}

func (f *Forms) New(info FormInfo) string {
  info.TrDef(f.t)
  if info.Method == "" {
    info.Method = "POST"
  }
  if info.EncType == "" {
    info.EncType = "multipart/form-data"
  }
  // Sorting by Order
  sort.Slice(info.Tabs, func(i, j int) bool { return info.Tabs[i].Order < info.Tabs[j].Order })
  sort.Slice(info.Properties, func(i, j int) bool { return info.Properties[i].Order < info.Properties[j].Order })
  sort.Slice(info.Buttons, func(i, j int) bool { return info.Buttons[i].Order < info.Buttons[j].Order })

  f.m.Set(info.Code, info)

  if glog.V(9) {
    glog.Infof("LOG: Load Form: '%s'", info.Code)
  }
  return info.Code
}

func (f *Forms) Get(code string) *FormInfo {
  var item FormInfo
  i, ok := f.m.Get(code, &item)
  if ok {
    if ref.GetType(i) == "FormInfo" {
      item, ok = i.(FormInfo)
      if ok {
        return &item
      }
    }
  }
  return nil
}

func (f *Forms) Filter(lang string, forms []string) map[string]FormInfo {
  res := make(map[string]FormInfo)
  for _, form_code := range forms {
    m := f.Get(form_code)
    if m != nil {
      m.Tr(f.t, lang)
      res[form_code] = (*m)
    }
  }
  return res
}

func (f *Forms) Load(configPath string) {
  env.LoadFromFiles(configPath, ".form", f.load)
}

func (f *Forms) GetParameters(r *http.Request, storagePath string, subPath string) (map[string]interface{}, bool) {
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
      params[part.FormName()] = f.readValueFromFormData(part)
      if glog.V(9) {
        glog.Infof("DBG: Set parameter '%s': `%v`", part.FormName(), f.readValueFromFormData(part))
      }
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
    
    f.appendFiles(&params, part.FormName(), filename)

    if glog.V(9) {
      glog.Infof("DBG: Set parameter '%s': `%v`", part.FormName(), filename)
    }

  }
  return params, true
}

func (f *Forms) appendFiles(params *map[string]interface{}, key string, filename string) {
  i := strings.Index(key, ":")
  if i < 0 {
    (*params)[key] = filename
  } else {
    name := key[0:i]
    index := key[i+1:]
    in, err := strconv.Atoi(index)
    if err == nil {
      if in == 0 {
        (*params)[name] = filename
      } else {
        if in == 1 {
          // is array. Revert to "name:0"
          (*params)[name + ":0"] = (*params)[name]
          (*params)[key] = filename
          delete((*params), name)
        }
      }
    }
  }
}

func (f *Forms) readValueFromFormData(part *multipart.Part) string {
  buf := new(bytes.Buffer)
  buf.ReadFrom(part)
  return buf.String()
}

func (f *Forms) readFileParam(part *multipart.Part) ([]byte, bool) {
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

func (f *Forms) ExpandProperties(info *FormInfo) {
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

func (f *Forms) load(filename string, buf []byte) int {
  var mapTmp FormInfo

  err := yaml.Unmarshal(buf, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: File(%s): YAML: %v", filename, err)
  }
  f.ExpandProperties(&mapTmp)
  f.New(mapTmp)

  return 1
}

func (f *Forms) index(lang string, form_code string, style string, isModal bool, data *map[string]interface{}) string {
  return "FORM#" + lang + "#" + form_code + "#" + style + "#" + strconv.FormatBool(isModal) + "#" + strconv.FormatBool(data == nil)
}

func (f *Forms) Render(form_code string, lang string, style string, isModal bool, data *map[string]interface{}) string {
  itemForm := f.Get(form_code)
  if itemForm == nil {
    err_str := fmt.Sprintf("ERR: FORM(%s) NOT FOUND ", form_code)
    glog.Errorf(err_str)
    return err_str
  }

  var form FormInfo
  form.Copy(itemForm)
  form.Tr(f.t, lang)
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

	tmpl := f.templates.Get(typeForm, "forms", style, lang)
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
    tmplProp := f.templates.Get("input_" + prop.WidgetType, "forms", style, lang)
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
        str, okstr := ref.ValueToString(v)
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
    
    tmplTabTitle   := f.templates.Get("tab_title",   "forms", style, lang)
    tmplTabContent := f.templates.Get("tab_content", "forms", style, lang)
    
    tmplTab := f.templates.Get(tTab, "forms", style, lang)
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
    tmplBtn := f.templates.Get("button_" + button.Type, "forms", style, lang)
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
  ref.UnionMaps(&propForm, data)
  tpl.Reset()
  err := tmpl.Execute(&tpl, propForm)
  if err != nil {
    err_str := fmt.Sprintf("ERR: templateFile(%s:%s).Execute: '%v'", style, form_code, err)
    glog.Errorf(err_str)
    return err_str
  }
  return tpl.String()
}

