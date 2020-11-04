package ui

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCheckTemplate(t *testing.T) {
  
  tmplProp := getTemplate("index111", "pages", "bootstrap4", "ru")
  assert.Nil(t, tmplProp)

  tmplProp = getTemplate("index", "pages", "bootstrap4", "ru")
  assert.NotNil(t, tmplProp)


  // err_need  := []error{ errors.New("len(c.Args)=2: Node {{TR \"Сервисы\"}} not supported"), errors.New("len(a.Ident): Node {{.User.Avatar}} not supported"), errors.New("len(c.Args)=2: Node {{TR \"Мой профиль\"}} not supported"), errors.New("len(c.Args)=2: Node {{TR \"Выход\"}} not supported"), errors.New("parse.TextNode: Node {{range .Services}}\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"{{.Url}}/?k={{.User_JWT}}\" target=\"_blank\">\n            <img src=\"{{.Img}}\" alt=\"\">\n            <h1>{{.Name}}</h1>\n            <p>{{.Desc}}</p>\n          </a>\n        </div>\n      </li>\n{{end}} not supported")}
  vars_need := []string{ "Title", "User_DisplayName", "LANG"}
  vars, _ := requiredTemplateVars(tmplProp)

  assert.Equal(t, vars_need, vars)
  // assert.Equal(t, err_need,  err)
  
  trs_need := []string{ "Сервисы", "Мой профиль", "Выход"}
  trs := findTrTemplate(tmplProp)
  assert.Equal(t, trs_need, trs)
}
