package ui

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "bytes"
  
  "github.com/Lunkov/lib-tr"
)

func TestCheckTemplate(t *testing.T) {
  
  translate := tr.New()
  templates := NewTemplates(translate, "./templates")
  
  // Public templates
  tmplProp := templates.Get("index111", "pages", "bootstrap4", "ru")
  assert.Nil(t, tmplProp)

  tmplProp = templates.Get("index", "pages", "bootstrap4", "ru")
  assert.NotNil(t, tmplProp)

  vars_need := []string{ "Title", "User_DisplayName", "LANG"}
  vars, _ := requiredTemplateVars(tmplProp)

  assert.Equal(t, vars_need, vars)
  
  trs_need := []string{ "Сервисы", "Мой профиль", "Выход"}
  trs := findTrTemplate(tmplProp)
  assert.Equal(t, trs_need, trs)

  // Private templates
  var tpl bytes.Buffer
  var ptpl bytes.Buffer

  tmplPage := templates.Get("private", "pages", "bootstrap4", "ru")
  assert.NotNil(t, tmplProp)

  vars_need = []string{ "Title", "LANG"}
  vars, _ = requiredTemplateVars(tmplPage)
  assert.Equal(t, vars_need, vars)
  
  propPage := map[string]interface{} {
      "Title": "Hello",
      "LANG": "ru",
    }
  
  err := tmplPage.Execute(&tpl, propPage)
  assert.Nil(t, err)
    
  assert.Equal(t, "  <body>\nHello\n\nLanguage: ru\n[[ .User_DisplayName ]]\n\n<ul id=\"hexGrid\">\n\n</ul>\n", tpl.String())
  
  tmplProp = templates.getPrivate("private", tpl.String(), "bootstrap4", "ru")
  assert.NotNil(t, tmplProp)

  vars_need = []string{ "User_DisplayName"}
  vars, _ = requiredTemplateVars(tmplProp)

  assert.Equal(t, vars_need, vars)
  

  propPage = map[string]interface{} {
      "User_DisplayName": "Serg",
    }
  
  err = tmplProp.Execute(&ptpl, propPage)
  assert.Nil(t, err)

  assert.Equal(t, "  <body>\nHello\n\nLanguage: ru\nSerg\n\n<ul id=\"hexGrid\">\n\n</ul>\n", ptpl.String())
  
}
