package ui

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-tr"
)

func TestCheckPage(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  res := tr.LoadLangs("./etc.test/langs.yaml")
  assert.Equal(t, true, res)
  tr.LoadTrs("./etc.test/tr")
  
  assert.Equal(t, 0, pageCount())
  
  f0 := FormInfo{Code: "welcome.1",
                 Title: "Welcome",
                 Properties: ArProperties{InfoProperty{Title: "Логин",  WidgetType: "text"}, InfoProperty{Title: "Пароль", WidgetType: "password"}},
                 Buttons:    ArButtons{InfoButton{Title: "Вход", Type: "submit", CODE: "submit"}, InfoButton{Title: "Закрыть", Type: "close", CODE: "close"}}}
  formNew(f0)
  // assert.Equal(t, 1, formCount())

  
  v0 := PageInfo{CODE: "page.home", Title: "Home page", Template: "index", Cache: true}
  pageNew(v0)
  
  assert.Equal(t, 1, pageCount())
  
  v1 := pageGet("page--1")
  assert.Nil(t, v1)

  v1 = pageGet("page.home")
  assert.Equal(t, &v0, v1)

  uid4, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  services := []map[string]interface{}{{"User_JWT": "111", "Name": "Service 1"}, {"User_JWT": "112", "Name": "Service 2"}, {"User_JWT": "131", "Name": "Service 3"}}
  data := map[string]interface{}{"id": uid4, "User_DisplayName": "Sergey Lunkov", "Services": services}
  
  render1_need := "\n<!doctype html>\n<html lang=\"ru\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\">\n    <link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/style.css\">\n    <title>#Домашняя страница#</title>\n  </head>\n  Hello, Sergey Lunkov!\n\n  <body>\n\t<div class=\"modal fade\" id=\"welcome.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Добро пожаловать</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"\">Логин</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Пароль</label>\n    <input type=\"password\" class=\"form-control\" id=\"\" placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"welcome.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Вход</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Закрыть</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n\n    <nav class=\"navbar navbar-dark bg-dark\">\n      <span class=\"navbar-brand mb-0 h1\"></span>\n      <ul class=\"navbar-nav mr-auto\">\n        <li class=\"nav-item active\">\n          <a class=\"nav-link\" href=\"/\">Сервисы</a>\n        </li>\n      </ul>\n      <li class=\"nav-item dropdown\">\n\t\t\t  <a class=\"nav-link dropdown-toggle row\" href=\"#\" role=\"button\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">\n\t\t\t\t<div class=\"column\">Sergey Lunkov</div><div class=\"column\"><img src=\"\" class=\"avatar\"></div>\n\t\t\t  </a>\n        <div class=\"dropdown-menu\" aria-labelledby=\"navbarDropdown\">\n          <a class=\"dropdown-item\" href=\"/profile\">Мой профиль</a>\n          <div class=\"dropdown-divider\"></div>\n          <a class=\"dropdown-item\" href=\"/logout\">Выход</a>\n        </div>\n       </li>\n    </nav>\n<br>\n<div class=\"container\">\n\nLanguage: ru_RU\n[[ .User_DisplayName ]]\n\n<ul id=\"hexGrid\">\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=111\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 1</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=112\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 2</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=131\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 3</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n</ul>\n\n--- VIEW 'table_edit_news' NOT FOUND ---\n\n\nBye, Sergey Lunkov!\n\n\n</div>\n    <script src=\"https://code.jquery.com/jquery-3.3.1.slim.min.js\" integrity=\"sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js\" integrity=\"sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js\" integrity=\"sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy\" crossorigin=\"anonymous\"></script>\n  </body>\n</html>\n"
  render1, okr := renderPage("ru_RU", "page.home", "bootstrap4", false, &data)
  assert.Equal(t, true, okr)
  assert.Equal(t, render1_need, render1)

  render1_need = "\n<!doctype html>\n<html lang=\"ru\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\">\n    <link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/style.css\">\n    <title>#Home page#</title>\n  </head>\n  Hello, Sergey Lunkov!\n\n  <body>\n\t<div class=\"modal fade\" id=\"welcome.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Welcome</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"\">Login</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Password</label>\n    <input type=\"password\" class=\"form-control\" id=\"\" placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"welcome.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Login</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Close</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n\n    <nav class=\"navbar navbar-dark bg-dark\">\n      <span class=\"navbar-brand mb-0 h1\"></span>\n      <ul class=\"navbar-nav mr-auto\">\n        <li class=\"nav-item active\">\n          <a class=\"nav-link\" href=\"/\">Services</a>\n        </li>\n      </ul>\n      <li class=\"nav-item dropdown\">\n\t\t\t  <a class=\"nav-link dropdown-toggle row\" href=\"#\" role=\"button\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">\n\t\t\t\t<div class=\"column\">Sergey Lunkov</div><div class=\"column\"><img src=\"\" class=\"avatar\"></div>\n\t\t\t  </a>\n        <div class=\"dropdown-menu\" aria-labelledby=\"navbarDropdown\">\n          <a class=\"dropdown-item\" href=\"/profile\">My profile</a>\n          <div class=\"dropdown-divider\"></div>\n          <a class=\"dropdown-item\" href=\"/logout\">Exit</a>\n        </div>\n       </li>\n    </nav>\n<br>\n<div class=\"container\">\n\nLanguage: en_US\n[[ .User_DisplayName ]]\n\n<ul id=\"hexGrid\">\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=111\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 1</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=112\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 2</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=131\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 3</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n</ul>\n\n--- VIEW 'table_edit_news' NOT FOUND ---\n\n\nBye, Sergey Lunkov!\n\n\n</div>\n    <script src=\"https://code.jquery.com/jquery-3.3.1.slim.min.js\" integrity=\"sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js\" integrity=\"sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js\" integrity=\"sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy\" crossorigin=\"anonymous\"></script>\n  </body>\n</html>\n"
  render1, okr = renderPage("en_US", "page.home", "bootstrap4", false, &data)
  assert.Equal(t, true, okr)
  assert.Equal(t, render1_need, render1)
  

  render1_need = "\n<!doctype html>\n<html lang=\"ru\">\n  <head>\n    <meta charset=\"utf-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\">\n    <link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/style.css\">\n    <title>#Home page#</title>\n  </head>\n  Hello, Sergey Lunkov!\n\n  <body>\n\t<div class=\"modal fade\" id=\"welcome.1\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Welcome</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <label for=\"\">Login</label>\n    <input type=\"text\" class=\"form-control\" id=\"\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"\">Password</label>\n    <input type=\"password\" class=\"form-control\" id=\"\" placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"welcome.1-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Login</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Close</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n\n    <nav class=\"navbar navbar-dark bg-dark\">\n      <span class=\"navbar-brand mb-0 h1\"></span>\n      <ul class=\"navbar-nav mr-auto\">\n        <li class=\"nav-item active\">\n          <a class=\"nav-link\" href=\"/\">Services</a>\n        </li>\n      </ul>\n      <li class=\"nav-item dropdown\">\n\t\t\t  <a class=\"nav-link dropdown-toggle row\" href=\"#\" role=\"button\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">\n\t\t\t\t<div class=\"column\">Sergey Lunkov</div><div class=\"column\"><img src=\"\" class=\"avatar\"></div>\n\t\t\t  </a>\n        <div class=\"dropdown-menu\" aria-labelledby=\"navbarDropdown\">\n          <a class=\"dropdown-item\" href=\"/profile\">My profile</a>\n          <div class=\"dropdown-divider\"></div>\n          <a class=\"dropdown-item\" href=\"/logout\">Exit</a>\n        </div>\n       </li>\n    </nav>\n<br>\n<div class=\"container\">\n\nLanguage: en_US\nSergey Lunkov\n\n<ul id=\"hexGrid\">\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=111\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 1</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=112\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 2</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n      <li class=\"hex\">\n        <div class=\"hexIn\">\n          <a class=\"hexLink\" href=\"/?k=131\" target=\"_blank\">\n            <img src=\"\" alt=\"\">\n            <h1>Service 3</h1>\n            <p></p>\n          </a>\n        </div>\n      </li>\n\n</ul>\n\n--- VIEW 'table_edit_news' NOT FOUND ---\n\n\nBye, Sergey Lunkov!\n\n\n</div>\n    <script src=\"https://code.jquery.com/jquery-3.3.1.slim.min.js\" integrity=\"sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js\" integrity=\"sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49\" crossorigin=\"anonymous\"></script>\n    <script src=\"https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js\" integrity=\"sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy\" crossorigin=\"anonymous\"></script>\n  </body>\n</html>\n"
  render1, okr = renderPage("en_US", "page.home", "bootstrap4", true, &data)
  assert.Equal(t, true, okr)
  assert.Equal(t, render1_need, render1)
  

  tr.SaveNew("./etc.test/tr")
}
