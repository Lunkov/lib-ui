package ui

import (
  "html/template"

  "github.com/Lunkov/lib-tr"
  "github.com/golang/glog"
)

type Functions struct {
  forms     *Forms
  views     *Views  
  f          template.FuncMap
  t         *tr.Tr
}

func NewFunctions(t *tr.Tr, forms *Forms, views *Views) *Functions {
  return &Functions{
        forms: forms,
        views: views,
        t:     t,
      }
}

func (f *Functions) AppendFuncMap(funcmap *template.FuncMap) {
  // u.templates.AppendFuncMap(funcmap)
}

func (f *Functions) FuncMap(name string, path string, style string, lang string) template.FuncMap {
  fmap := template.FuncMap{
                  "TR": func(str string) string {
                      t, _ := f.t.Tr(lang, str)
                      return t
                  },
                  "TR_LANG": func() string {
                      return lang
                  },
                  "TR_LANG_NAME": func() string {
                      return f.t.LangName(lang)
                  },
                  "FORM": func(form_code string, data map[string]interface{}) template.HTML {
                      glog.Infof("@@@@@@@@@@@ DBG: templateFile FORM(%s) ", form_code)
                      if f.forms == nil {
                        return template.HTML("--- FORM '" + form_code + "' NOT INITED ---")
                      }
                      return template.HTML(f.forms.Render(form_code, lang, style, false, &data))
                  },
                  "MODAL": func(form_code string, data map[string]interface{}) template.HTML {
                      glog.Infof("@@@@@@@@@@22 DBG: templateFile FORM(%s) ", form_code)
                      if f.forms == nil {
                        return template.HTML("--- FORM '" + form_code + "' NOT INITED ---")
                      }
                      return template.HTML(f.forms.Render(form_code, lang, style, true, &data))
                  },
                  "VIEW": func(view_code string, data map[string]interface{}) template.HTML {
                      if f.views == nil {
                        return template.HTML("--- VIEW '" + view_code + "' NOT INITED ---")
                      }
                      res, ok := f.views.Render(view_code, lang, style, &data)
                      if ok {
                        return template.HTML(res)
                      }
                      return template.HTML("--- VIEW '" + view_code + "' NOT FOUND ---")
                  },
                  "attr":func(s string) template.HTMLAttr{
                      return template.HTMLAttr(s)
                  },
                  "safe": func(s string) template.HTML {
                      return template.HTML(s)
                  },
                  "url": func(s string) template.URL {
                      return template.URL(s)
                  },
                }
  //for k, v := range fmap {
  //  f.f[k] = v
  //}
  return fmap
}

