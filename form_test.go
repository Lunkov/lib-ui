package ui

import (
  "flag"

  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/google/uuid"
  "github.com/Lunkov/lib-tr"
)

func TestCheckValue2String(t *testing.T) {
  
  res, ok := value2String("112")
  assert.Equal(t, true, ok)
  assert.Equal(t, "112", res)
  
  res, ok = value2String(112)
  assert.Equal(t, true, ok)
  assert.Equal(t, "112", res)

  res, ok = value2String(-112)
  assert.Equal(t, true, ok)
  assert.Equal(t, "-112", res)

  res, ok = value2String(-112.4568)
  assert.Equal(t, true, ok)
  assert.Equal(t, "-112.4568", res)


  // WARNING
  res, ok = value2String(-3453455645645672.45345646474568)
  assert.Equal(t, true, ok)
  assert.Equal(t, "-3.4534556456456725e+15", res)

}

func TestCheckForm(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

  res := tr.LoadLangs("./etc.test/langs.yaml")
  assert.Equal(t, true, res)
  tr.LoadTrs("./etc.test/tr")
  assert.Equal(t, 4, tr.Count())
  
  tr_str_need := "Password"
  tr_str, _ := tr.Tr("en_US", "Пароль")
  assert.Equal(t, tr_str_need, tr_str)
  tr_str, _ = tr.Tr("en", "Password")
  assert.Equal(t, tr_str_need, tr_str)

  assert.Equal(t, 0, formCount())
  
  v0 := FormInfo{Code: "welcome.1",
                 Title: "Welcome",
                 Properties: ArProperties{InfoProperty{Order: 20, Title: "Пароль", WidgetType: "password"}, InfoProperty{Title: "Логин",  Order: 10, WidgetType: "text"}, InfoProperty{Title: "Имя",  Order: 15, WidgetType: "text"}, InfoProperty{Title: "Фамилия",  Order: 17, WidgetType: "text"}},
                 Buttons:    ArButtons{InfoButton{Order: 20, Title: "Закрыть", Type: "close", CODE: "close"},InfoButton{Order: 10, Title: "Вход", Type: "submit", CODE: "submit"}}}
  formNew(v0)
  
  assert.Equal(t, 1, formCount())
  
  v1 := formGet("view--1")
  assert.Nil(t, v1)

  v1_need := FormInfo{Code:"welcome.1", EncType:"multipart/form-data", Method:"POST", Action:"", Version:0, Date:"", Date_start:"", Date_finish:"", Title:"Welcome", Type:"", Description:"", Tabs:ArTabs(nil), Properties:ArProperties{InfoProperty{CODE:"", Order:10, Title:"Логин", Description:"", DataType:"", WidgetType:"text", Required:"", Format:"", Table:"", Default:"", OneOf:[]InfoItem(nil), OneOfConst:""}, InfoProperty{CODE:"", Order:15, Title:"Имя", Description:"", DataType:"", WidgetType:"text", Required:"", Format:"", Table:"", Default:"", OneOf:[]InfoItem(nil), OneOfConst:""}, InfoProperty{CODE:"", Order:17, Title:"Фамилия", Description:"", DataType:"", WidgetType:"text", Required:"", Format:"", Table:"", Default:"", OneOf:[]InfoItem(nil), OneOfConst:""}, InfoProperty{CODE:"", Order:20, Title:"Пароль", Description:"", DataType:"", WidgetType:"password", Required:"", Format:"", Table:"", Default:"", OneOf:[]InfoItem(nil), OneOfConst:""}}, Buttons:ArButtons{InfoButton{CODE:"submit", Order:10, Title:"Вход", Type:"submit"}, InfoButton{CODE:"close", Order:20, Title:"Закрыть", Type:"close"}}}
  v1 = formGet("welcome.1")
  assert.Equal(t, &v1_need, v1)
  
  render1_need := "<div class=\"modal fade\" id=\"welcome.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Добро пожаловать</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"\">Логин</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Имя</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Фамилия</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Пароль</label>\n    <input type=\"password\" class=\"form-control\" id=\"\" placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"welcome.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Вход</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Закрыть</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n"
  render1 := renderForm("ru_RU", "welcome.1", "bootstrap4", true, nil)
  assert.Equal(t, render1_need, render1)

  render1_need = "<div class=\"modal fade\" id=\"welcome.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Welcome</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"\">Login</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Name</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Last name</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Password</label>\n    <input type=\"password\" class=\"form-control\" id=\"\" placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"welcome.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Login</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Close</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n"
  render1 = renderForm("en_US", "welcome.1", "bootstrap4", true, nil)
  assert.Equal(t, render1_need, render1)

  vOrg := FormInfo{Code: "organization.1", Action: "http://localhost:3000/organization/", Title: "Organization", EncType:"multipart/form-data", Method:"POST",
                   Properties: ArProperties{InfoProperty{CODE: "name", Order: 10, Title: "Наименование",  WidgetType: "text"}, InfoProperty{CODE: "url_logo", Order: 20, Title: "Лого", WidgetType: "text"}},
                   Buttons:    ArButtons{InfoButton{Order: 10, Title: "Вход", Type: "submit", CODE: "submit"}, InfoButton{Order: 20, Title: "Закрыть", Type: "close", CODE: "close"}}}
  
  formNew(vOrg)
  
  uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  data := map[string]interface{}{"id": uid4, "code": "test.org.4", "name": "ORG `Test NEW 21`", "address_legal.city": "Moscow", "address_legal.index": "127282", "address_legal.country": "Russia", "bank.0.bik": "0065674747", "bank.0.account": "2342345446560065674747"}
  
  render1_need = "<div class=\"modal fade\" id=\"organization.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"http://localhost:3000/organization/\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Организация</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"name\">Наименование</label>\n    <input type=\"text\" class=\"form-control\" id=\"name\" placeholder=\"\" value=\"ORG `Test NEW 21`\">\n</div>\n<div class=\"form-group\">\n    <label for=\"url_logo\">Лого</label>\n    <input type=\"text\" class=\"form-control\" id=\"url_logo\" placeholder=\"\" value=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"organization.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Вход</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Закрыть</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n"
  render1 = renderForm("ru_RU", "organization.1", "bootstrap4", true, &data)
  assert.Equal(t, render1_need, render1)

  render1_need = "<div class=\"modal fade\" id=\"organization.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"http://localhost:3000/organization/\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Organization</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"name\">Name</label>\n    <input type=\"text\" class=\"form-control\" id=\"name\" placeholder=\"\" value=\"ORG `Test NEW 21`\">\n</div>\n<div class=\"form-group\">\n    <label for=\"url_logo\">Logo</label>\n    <input type=\"text\" class=\"form-control\" id=\"url_logo\" placeholder=\"\" value=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"organization.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Login</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Close</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n"
  render1 = renderForm("en_US", "organization.1", "bootstrap4", true, &data)
  assert.Equal(t, render1_need, render1)

  tr.SaveNew("./etc.test/tr")
}

func TestCheckFormTabs(t *testing.T) {
  vOrg := FormInfo{Code: "organization.1", Action: "http://localhost:3000/organization/", Title: "Organization", EncType:"multipart/form-data", Method:"POST",
                   Properties: ArProperties{InfoProperty{ Order: 10, CODE: "name", Title: "Наименование",  WidgetType: "text"}, InfoProperty{Order: 20, CODE: "url_logo", Title: "Лого", WidgetType: "text"}, InfoProperty{CODE: "address_legal.country", Title: "Страна", WidgetType: "text", Order: 30}, InfoProperty{CODE: "address_legal.city", Title: "Город", WidgetType: "text", Order: 40}},
                   Tabs:       ArTabs{ InfoTab{ Order: 10, CODE: "main", Title: "Главная", Fields: []string{"name", "url_logo"}}, InfoTab{Order: 20, CODE: "address_legal", Title: "Юридический адрес", Fields: []string{"address_legal.country", "address_legal.city"}}, InfoTab{CODE: "bank", Title: "Банковские реквизиты", Fields: []string{"bank.0.bik", "bank.0.account"}}},
                   Buttons:    ArButtons{InfoButton{Title: "Вход", Type: "submit", CODE: "submit", Order: 10}, InfoButton{Title: "Закрыть", Type: "close", CODE: "close", Order: 20}}}
  
  formNew(vOrg)
  
  uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  data := map[string]interface{}{"id": uid4, "code": "test.org.4", "name": "ORG `Test NEW 21`", "address_legal.city": "Moscow", "address_legal.index": "127282", "address_legal.country": "Russia", "bank.0.bik": "0065674747", "bank.0.account": "2342345446560065674747"}
  
  render1_need := "<div class=\"modal fade\" id=\"organization.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"http://localhost:3000/organization/\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Организация</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"address_legal.country\">Страна</label>\n    <input type=\"text\" class=\"form-control\" id=\"address_legal.country\" placeholder=\"\" value=\"Russia\">\n</div>\n<div class=\"form-group\">\n    <label for=\"address_legal.city\">Город</label>\n    <input type=\"text\" class=\"form-control\" id=\"address_legal.city\" placeholder=\"\" value=\"Moscow\">\n</div>\n<ul class=\"nav nav-tabs\" id=\"TAB-organization.1\" role=\"tablist\">\n<li class=\"nav-item\">\n  <a class=\"nav-link active\" id=\"bank-tab\" data-toggle=\"tab\" href=\"#bank\" role=\"tab\" aria-controls=\"bank\" aria-selected=\"true\">Банковские реквизиты</a>\n</li>\n<li class=\"nav-item\">\n  <a class=\"nav-link \" id=\"main-tab\" data-toggle=\"tab\" href=\"#main\" role=\"tab\" aria-controls=\"main\" aria-selected=\"true\">Главная</a>\n</li>\n<li class=\"nav-item\">\n  <a class=\"nav-link \" id=\"address_legal-tab\" data-toggle=\"tab\" href=\"#address_legal\" role=\"tab\" aria-controls=\"address_legal\" aria-selected=\"true\">Юридический адрес</a>\n</li>\n\n</ui>\n<div class=\"tab-content\" id=\"TAB-organization.1-CONTENT\">>\n<div class=\"tab-pane fade active show\" id=\"bank\" role=\"tabpanel\" aria-labelledby=\"bank-tab\">\n  \n</div>\n<div class=\"tab-pane fade \" id=\"main\" role=\"tabpanel\" aria-labelledby=\"main-tab\">\n  <div class=\"form-group\">\n    <label for=\"name\">Наименование</label>\n    <input type=\"text\" class=\"form-control\" id=\"name\" placeholder=\"\" value=\"ORG `Test NEW 21`\">\n</div>\n<div class=\"form-group\">\n    <label for=\"url_logo\">Лого</label>\n    <input type=\"text\" class=\"form-control\" id=\"url_logo\" placeholder=\"\" value=\"\">\n</div>\n\n</div>\n<div class=\"tab-pane fade \" id=\"address_legal\" role=\"tabpanel\" aria-labelledby=\"address_legal-tab\">\n  <div class=\"form-group\">\n    <label for=\"address_legal.country\">Страна</label>\n    <input type=\"text\" class=\"form-control\" id=\"address_legal.country\" placeholder=\"\" value=\"Russia\">\n</div>\n<div class=\"form-group\">\n    <label for=\"address_legal.city\">Город</label>\n    <input type=\"text\" class=\"form-control\" id=\"address_legal.city\" placeholder=\"\" value=\"Moscow\">\n</div>\n\n</div>\n\n</div>\n<script>\n  $(function () {\n    $('#TAB-organization.1 li:last-child a').tab('show')\n  })\n</script>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"organization.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Вход</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Закрыть</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n"
  render1 := renderForm("ru_RU", "organization.1", "bootstrap4", true, &data)
  assert.Equal(t, render1_need, render1)
}
