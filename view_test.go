package ui

import (
  "flag"
  
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestCheckView(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  Init("./etc.test/", false, true)
  assert.Equal(t, 1, viewCount())
  
  render1_need := "<div class=\"modal fade\" id=\"news\" tabindex=\"-1\" role=\"dialog\" aria-hidden=\"true\">\n<form method=\"POST\" action=\"/api/news_public\" enctype=\"multipart/form-data\">\n  <div class=\"modal-dialog\">\n    <div class=\"modal-content\">\n      <div class=\"modal-header\">\n        <h5 class=\"modal-title\">Редактирование новости</h5>\n        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-label=\"Close\">\n          <span aria-hidden=\"true\">&times;</span>\n        </button>\n      </div>\n      <div class=\"modal-body\">\n<div class=\"form-group\">\n    <input type=\"hidden\" class=\"form-control\" id=\"id\" name=\"id\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"title\">Заголовок</label>\n    <input type=\"text\" class=\"form-control\" id=\"title\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"public_date\">Дата публикации</label>\n    <input type=\"text\" class=\"form-control\" id=\"public_date\" name=\"public_date\" placeholder=\"\" value=\"\">\n</div>\n<div class=\"form-group\">\n    <label for=\"description\">Описание</label>\n    <textarea id=\"description\" name=\"description\"></textarea>\n<script>\n$(\"#description\").sceditor({\n\t  plugins: \"xhtml\",\n\t  toolbar: \"bold,italic,underline,strike,subscript,superscript|left,center,right,justify|size,color,removeformat|cut,copy,pastetext|bulletlist,orderedlist,indent,outdent|table|code,quote|horizontalrule,image,email,link,unlink|youtube,time|ltr,rtl|print,maximize,source\",\n\t  charset: \"utf-8\",\n\t  style: '/static/css/sceditor/themes/default.min.css',\n\t  emoticonsEnabled: false,\n\t  resizeEnabled: false,\n\t  height:600,\n\t  spellcheck: true,\n\t});\n</script>\n</div>\n<div class=\"form-group\">\n    <label for=\"images\">Файлы</label>\n    <input type=\"file\" class=\"form-control\" id=\"images\" name=\"images[]\"  data-show-upload=\"false\" data-show-caption=\"true\" multiple placeholder=\"\">\n</div>\n\n      </div>\n      <div class=\"modal-footer\">\n<button id=\"news-SUBMIT-BUTTON\" type=\"button\" class=\"btn btn-primary\">\n  <span>Save</span>\n</button>\n<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Close</button>\n      </div>\n    </div>\n  </div>\n</form>\n</div>\n\n<div class=\"row\">\n\t<div class=\"col-lg-10\">\n    \n\t\t<table\n\t\t  id=\"table-table_edit_news\"\n\t\t  data-toolbar=\"#toolbar-table_edit_news\"\n\t\t  data-toggle=\"table\"\n\n\t\t  data-id-field=\"id\"\n\t\t  data-click-to-select=\"true\"\n\t\t  data-url=\"\"\n\t\t  data-response-handler=\"responseTableHandler\">\n\t\t</table>\n\n\t</div>\n</div>\n\n<script>\n\t\nclass tableEdit {\n  constructor() {\n    this.edit_form = \"news\";\n    this.view_url = \"\\/page\\/1news\\/ru_RU?item=\";\n    this.read_item_url = \"\";\n    this.write_item_url = \"\";\n    this.table = $('#table-table_edit_news');\n  }\n  \n  $this.table.bootstrapTable('destroy').bootstrapTable({\n      locale: \"en_US\", \n      columns: [\n      \n        {\n          field: 'id',\n          title: 'ID',\n          sortable:  false ,\n          align: 'center'\n        },\n      \n        {\n          field: 'title',\n          title: 'Title',\n          sortable:  true ,\n          align: 'left'\n        },\n      \n        {\n          field: 'public_at',\n          title: 'Public At',\n          sortable:  true ,\n          align: 'center'\n        },\n      \n        {\n          field: 'operate',\n          title: 'Действия',\n          align: 'left',\n          clickToSelect: false,\n          events: window.operateEvents,\n          formatter: operateFormatter\n        }]\n    })\n}\n</script>\n"
  render1, okr := renderView("en_US", "table_edit_news", "bootstrap4", nil)
  assert.Equal(t, true, okr)
  assert.Equal(t, render1_need, render1)
  
}
