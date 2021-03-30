package ui

import (
  "bytes"
  "strconv"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-ref"
  "github.com/Lunkov/lib-tr"
  "github.com/Lunkov/lib-cache"
)

type Pages struct {
  t               *tr.Tr
  m                cache.ICache
  templates       *Templates
}

func NewPages(cfg *cache.CacheConfig, t *tr.Tr, ts *Templates) *Pages {
  return &Pages{ m: cache.NewConfig(cfg), t: t, templates: ts }
}

func (p *Pages) Count() int64 {
  return p.m.Count()
}

func (p *Pages) New(info PageInfo) string {
  info.TrDef(p.t)
  p.m.Set(info.CODE, info)
  if glog.V(9) {
    glog.Infof("DBG: Load Page: '%s'", info.CODE)
  }
  return info.CODE
}

func (p *Pages) Get(code string) *PageInfo {
  var item PageInfo
  i, ok := p.m.Get(code, &item)
  if ok {
    if ref.GetType(i) == "PageInfo" {
      item, ok = i.(PageInfo)
      if ok {
        return &item
      }
    }
  }
  return nil
}

func (p *Pages) Filter(lang string, pages []string) map[string]PageInfo {
  res := make(map[string]PageInfo)
  for _, page_code := range pages {
    m := p.Get(page_code)
    if m != nil {
      m.Tr(p.t, lang)
      res[page_code] = (*m)
    }
  }
  return res
}

func (p *Pages) Load(configPath string) {
  env.LoadFromFiles(configPath, ".page", p.load)
}

func (p *Pages) load(filename string, buf []byte) int {
  var mapTmp PageInfo

  err := yaml.Unmarshal(buf, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: File(%s): YAML: %v", filename, err)
  }
  p.New(mapTmp)

  return 1
}

func (p *Pages) index(lang string, page_code string, style string, private bool) string {
  return "PAGE#" + lang + "#" + page_code + "#" + style + "#" + strconv.FormatBool(private)
}

func (p *Pages) Render(page_code string, lang string, style string, private bool, data *map[string]interface{}) (string, bool) {
  pg := p.Get(page_code)
  if pg == nil {
    glog.Errorf("ERR: PAGE(%s) NOT FOUND ", page_code)
    return "", false
  }
  var page PageInfo
  page.Copy(pg)
  page.Tr(p.t, lang)
  
  var tpl bytes.Buffer
  
  tmplPage := p.templates.Get(page.Template, "pages", style, lang)
  if tmplPage == nil {
    glog.Errorf("ERR: LOAD: templateFile(%s)", page_code)
    return "", false
  }
  // For multi language pages
  if page.MultiLanguage {
    tmplPageML := p.templates.Get(page.Template + "_" + lang, "pages", style, lang)
    if tmplPageML != nil {
      tmplPage = tmplPageML
    }
  }
  
  trMap := p.templates.makeTrMap(tmplPage, lang)
  
  propPage := map[string]interface{} {
      "PAGE_ID": page.CODE,
      "CODE": page.CODE,
      "TITLE": page.Title,
      "DESCRIPTION": page.Description,
      "KEYWORDS": page.Keywords,
      "LANG": lang,
    }
  ref.UnionMaps(&propPage, data)
  ref.UnionMapsStr(&propPage, &trMap)
  
  glog.Infof("********: templateFile PAGE(%v) ", propPage)
  err := tmplPage.Execute(&tpl, propPage)
  if err != nil {
    glog.Errorf("ERR: templateFile(%s).Execute: '%v'", page_code, err)
    return "", false
  }
  
  if private {
    tmplPage := p.templates.getPrivate(page.Template, tpl.String(), style, lang)
    if tmplPage == nil {
      glog.Errorf("ERR: LOAD: templatePrivate(%s)", page_code)
      return "", false
    }
    tpl.Reset()
    glog.Infof("*******22:propPage templateFile PAGE(%v) ", propPage)
    err := tmplPage.Execute(&tpl, propPage)
    if err != nil {
      glog.Errorf("ERR: templateFile(%s).Execute: '%v'", page_code, err)
      return "", false
    }
  }
  return tpl.String(), true
}
