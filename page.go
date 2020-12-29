package ui

import (
  "bytes"
  "strconv"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-tr"
)

type PageInfo struct {
  CODE           string                 `json:"code"              yaml:"code"`
  Title          string                 `json:"title"             yaml:"title"`
  Keywords       string                 `json:"keywords"          yaml:"keywords"`             // MAX 250
  Description    string                 `json:"description"       yaml:"description"`          // MAX 140
  Template       string                 `json:"template"          yaml:"template"`
  Cache          bool                   `json:"cache"             yaml:"cache"`
  MultiLanguage  bool                   `json:"multi_language"    yaml:"multi_language"`
}

var mapPages = make(map[string]PageInfo)

func pageCount() int {
  return len(mapPages)
}

func pageNew(info PageInfo) string {
  info.TrDef()
  if glog.V(9) {
    glog.Infof("LOG: Load Page: '%s'", info.CODE)
  }
  mapPages[info.CODE] = info
  return info.CODE
}

func pageGet(code string) *PageInfo {
  i, ok := mapPages[code]
  if ok {
    return &i
  }
  return nil
}

func (info *PageInfo) TrDef() {
  tr.SetDef(info.Title)
  tr.SetDef(info.Description)
}

func (info *PageInfo) Tr(lang string) {
  info.Title, _       = tr.Tr(lang, info.Title)
  info.Description, _ = tr.Tr(lang, info.Description)
  info.Keywords, _    = tr.Tr(lang, info.Keywords)
}

func (info *PageInfo) Copy(in *PageInfo) {
  buf, _ := json.Marshal(in)
  json.Unmarshal(buf, info)
}

func pageFilter(lang string, pages []string) map[string]PageInfo {
  res := make(map[string]PageInfo)
  for _, page_code := range pages {
    m := pageGet(page_code)
    if m != nil {
      m.Tr(lang)
      res[page_code] = (*m)
    }
  }
  return res
}

func pageInit(configPath string) {
  mapPages = make(map[string]PageInfo)
  env.LoadFromYMLFiles(configPath, loadPageYAML)
}

func loadPageYAML(filename string, yamlFile []byte) int {
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return 0
  }
  var mapTmp PageInfo

  err = yaml.Unmarshal(yamlFile, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  pageNew(mapTmp)

  return 1
}

func indexPage(lang string, page_code string, style string, private bool) string {
  return "PAGE#" + lang + "#" + page_code + "#" + style + "#" + strconv.FormatBool(private)
}

func renderPage(lang string, page_code string, style string, private bool, data *map[string]interface{}) (string, bool) {
  p := pageGet(page_code)
  if p == nil {
    glog.Errorf("ERR: PAGE(%s) NOT FOUND ", page_code)
    return "", false
  }
  var page PageInfo
  page.Copy(p)
  page.Tr(lang)
  
  var tpl bytes.Buffer
  
  tmplPage := getTemplate(page.Template, "pages", style, lang)
  if tmplPage == nil {
    glog.Errorf("ERR: LOAD: templateFile(%s)", page_code)
    return "", false
  }
  // For multi language pages
  if page.MultiLanguage {
    tmplPageML := getTemplate(page.Template + "_" + lang, "pages", style, lang)
    if tmplPageML != nil {
      tmplPage = tmplPageML
    }
  }
  
  trMap := makeTrMap(tmplPage, lang)
  
  propPage := map[string]interface{} {
      "PAGE_ID": page.CODE,
      "CODE": page.CODE,
      "TITLE": page.Title,
      "DESCRIPTION": page.Description,
      "KEYWORDS": page.Keywords,
      "LANG": lang,
    }
  unionMap(&propPage, data)
  unionMapStr(&propPage, &trMap)
  
  err := tmplPage.Execute(&tpl, propPage)
  if err != nil {
    glog.Errorf("ERR: templateFile(%s).Execute: '%v'", page_code, err)
    return "", false
  }
  
  if private {
    tmplPage := getPrivateTemplate(page.Template, tpl.String(), style, lang)
    if tmplPage == nil {
      glog.Errorf("ERR: LOAD: templatePrivate(%s)", page_code)
      return "", false
    }
    tpl.Reset()
    err := tmplPage.Execute(&tpl, propPage)
    if err != nil {
      glog.Errorf("ERR: templateFile(%s).Execute: '%v'", page_code, err)
      return "", false
    }
  }
  return tpl.String(), true
}
