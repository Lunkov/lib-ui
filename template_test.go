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

  vars_need := []string{ "Title", "User_DisplayName", "LANG"}
  vars, _ := requiredTemplateVars(tmplProp)

  assert.Equal(t, vars_need, vars)
  
  trs_need := []string{ "Сервисы", "Мой профиль", "Выход"}
  trs := findTrTemplate(tmplProp)
  assert.Equal(t, trs_need, trs)
}
