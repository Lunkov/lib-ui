package ui

import (
  "bytes"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-tr"
)

type InfoAction struct {
  Title   string    `json:"title"    yaml:"title"`
  Url     string    `json:"url"      yaml:"url"`
  Target  string    `json:"target"   yaml:"target"`
}

type ColumnInfo struct {
  Title                 string  `json:"title"                yaml:"title"`
  Type                  string  `json:"type"                 yaml:"type"`
  Field                 string  `json:"field"                yaml:"field"`
  Align                 string  `json:"align"                yaml:"align"`
  Sortable              bool    `json:"sortable"             yaml:"sortable"`
  Filter                string  `json:"filter"               yaml:"filter"`
  valuePrepareFunction  string  `json:"valuePrepareFunction" yaml:"valuePrepareFunction"`
}

type TableInfo struct {
  Name       string                 `json:"name"     yaml:"name"`
  DataUrl    string                 `json:"url"      yaml:"url"`
  Columns    []ColumnInfo           `json:"columns"  yaml:"columns"`
}

type MapInfo struct {
  Name           string         `json:"name"           yaml:"name"`
  Latitude       float64        `json:"latitude"       yaml:"latitude"`
  Longitude      float64        `json:"longitude"      yaml:"longitude"`
  DefaultZoom    float64        `json:"default_zoom"   yaml:"default_zoom"`
  MinZoom        float64        `json:"min_zoom"       yaml:"min_zoom"`
  MaxZoom        float64        `json:"max_zoom"       yaml:"max_zoom"`
  LayersUrl      string         `json:"layers_url"     yaml:"layers_url"`
}

type ViewInfo struct {
  CODE           string                 `json:"code"          yaml:"code"`
  Title          string                 `json:"title"         yaml:"title"`
  Description    string                 `json:"description"   yaml:"description"`
  MainModel      string                 `json:"main_model"    yaml:"main_model"`
  EditForm       string                 `json:"edit_form"     yaml:"edit_form"`
  ReadItemUrl    string                 `json:"read_item_url"   yaml:"read_item_url"`
  WriteItemUrl   string                 `json:"write_item_url"  yaml:"write_item_url"`
  ViewItemUrl    string                 `json:"view_item_url"   yaml:"view_item_url"`
  Table          TableInfo              `json:"table"         yaml:"table,flow"`
  Map            MapInfo                `json:"map"           yaml:"map,flow"`
  Function       string                 `json:"function"      yaml:"function"`
  Template       string                 `json:"template"      yaml:"template"`
  Actions        map[string]InfoAction  `json:"actions"       yaml:"actions,flow"`
}

var mapViews = make(map[string]ViewInfo)

func viewCount() int {
  return len(mapViews)
}

func (info *ViewInfo) Copy(in *ViewInfo) {
  buf, _ := json.Marshal(in)
  json.Unmarshal(buf, info)
}

func viewNew(info ViewInfo) string {
  info.TrDef()
  if glog.V(9) {
    glog.Infof("LOG: Load View: '%s'", info.CODE)
  }
  mapViews[info.CODE] = info
  return info.CODE
}

func viewGet(code string) *ViewInfo {
  i, ok := mapViews[code]
  if ok {
    return &i
  }
  return nil
}

func (info *ViewInfo) TrDef() {
  tr.SetDef(info.Title)
  tr.SetDef(info.Description)
  for _, column := range info.Table.Columns {
    tr.SetDef(column.Title)
  }
  for _, action := range info.Actions {
    tr.SetDef(action.Title)
  }
}

func (info *ViewInfo) Tr(lang string) {
  info.Title, _ = tr.Tr(lang, info.Title)
  info.Description, _ = tr.Tr(lang, info.Description)
  for _, column := range info.Table.Columns {
    column.Title, _ = tr.Tr(lang, column.Title)
  }
  for _, action := range info.Actions {
    action.Title, _ = tr.Tr(lang, action.Title)
  }
}

func viewFilter(lang string, views []string) map[string]ViewInfo {
  res := make(map[string]ViewInfo)
  for _, view_code := range views {
    m := viewGet(view_code)
    if m != nil {
      m.Tr(lang)
      res[view_code] = (*m)
    }
  }
  return res
}

func viewInit(configPath string) {
  env.LoadFromFiles(configPath, "", loadViewYAML)
}

func loadViewYAML(filename string, yamlFile []byte) int {
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return 0
  }
  var mapTmp ViewInfo

  err = yaml.Unmarshal(yamlFile, &mapTmp)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  viewNew(mapTmp)

  return 1
}

func indexView(lang string, view_code string, style string) string {
  return "VIEW#" + lang + "#" + view_code + "#" + style
}

func renderView(lang string, view_code string, style string, data *map[string]interface{}) (string, bool) {
  p := viewGet(view_code)
  if p == nil {
    glog.Errorf("ERR: VIEW(%s) NOT FOUND ", view_code)
    return "", false
  }
  var view ViewInfo
  view.Copy(p)
  view.Tr(lang)
  
  var tpl bytes.Buffer
  
  tmplView := getTemplate(view.Template, "views", style, lang)
  if tmplView == nil {
    glog.Errorf("ERR: LOAD: templateFile(%s)", view_code)
    return "", false
  }
  
  trMap := makeTrMap(tmplView, lang)
  propView := map[string]interface{} {
      "VIEW_ID": view.CODE,
      "CODE": view.CODE,
      "TITLE": view.Title,
      "DESCRIPTION": view.Description,
      "LANG": lang,
      "VIEW": view,
    }
  unionMap(&propView, data)
  unionMapStr(&propView, &trMap)
  
  err := tmplView.Execute(&tpl, &propView)
  if err != nil {
    glog.Errorf("ERR: templateFile(%s).Execute: '%v'", view_code, err)
    return "", false
  }
  return tpl.String(), true
}
