package ui

  import (
  "bytes"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-ref"
  "github.com/Lunkov/lib-cache"
  "github.com/Lunkov/lib-tr"
)

type Views struct {
  t               *tr.Tr
  m                cache.ICache
  templates       *Templates
}

func NewViews(cfg *cache.CacheConfig, t *tr.Tr, ts *Templates) *Views {
  return &Views{ m: cache.NewConfig(cfg), t: t, templates: ts }
}

func (v *Views) Count() int64 {
  return v.m.Count()
}

func (v *Views) New(info ViewInfo) string {
  info.TrDef(v.t)
  v.m.Set(info.CODE, info)
  if glog.V(9) {
    glog.Infof("DBG: Load View: '%s'", info.CODE)
  }
  return info.CODE
}

func (v *Views) Get(code string) *ViewInfo {
  var item ViewInfo
  i, ok := v.m.Get(code, &item)
  if ok {
    if ref.GetType(i) == "ViewInfo" {
      item, ok = i.(ViewInfo)
      if ok {
        return &item
      }
    }
  }
  return nil
}

func (v *Views) Filter(lang string, views []string) map[string]ViewInfo {
  res := make(map[string]ViewInfo)
  for _, view_code := range views {
    m := v.Get(view_code)
    if m != nil {
      m.Tr(v.t, lang)
      res[view_code] = (*m)
    }
  }
  return res
}

func (v *Views) Load(configPath string) {
  env.LoadFromFiles(configPath, ".view", v.load)
}

func (v *Views) load(filename string, buf []byte) int {
  var mapTmp ViewInfo

  err := yaml.Unmarshal(buf, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: File(%s): YAML: %v", filename, err)
  }
  v.New(mapTmp)

  return 1
}

func (v *Views) index(lang string, view_code string, style string) string {
  return "VIEW#" + lang + "#" + view_code + "#" + style
}

func (v *Views) Render(view_code string, lang string, style string, data *map[string]interface{}) (string, bool) {
  p := v.Get(view_code)
  if p == nil {
    glog.Errorf("ERR: VIEW(%s) NOT FOUND ", view_code)
    return "", false
  }
  var view ViewInfo
  view.Copy(p)
  view.Tr(v.t, lang)
  
  var tpl bytes.Buffer
  
  tmplView := v.templates.Get(view.Template, "views", style, lang)
  if tmplView == nil {
    glog.Errorf("ERR: LOAD: templateFile(%s)", view_code)
    return "", false
  }
  
  trMap := v.templates.makeTrMap(tmplView, lang)
  propView := map[string]interface{} {
      "VIEW_ID": view.CODE,
      "CODE": view.CODE,
      "TITLE": view.Title,
      "DESCRIPTION": view.Description,
      "LANG": lang,
      "VIEW": view,
    }
  ref.UnionMaps(&propView, data)
  ref.UnionMapsStr(&propView, &trMap)
  
  err := tmplView.Execute(&tpl, &propView)
  if err != nil {
    glog.Errorf("ERR: templateFile(%s).Execute: '%v'", view_code, err)
    return "", false
  }
  return tpl.String(), true
}
