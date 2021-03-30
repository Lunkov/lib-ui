package ui

import (
  "github.com/jinzhu/copier"
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

func (info *ViewInfo) TrDef(t *tr.Tr) {
  t.SetDef(info.Title)
  t.SetDef(info.Description)
  for _, column := range info.Table.Columns {
    t.SetDef(column.Title)
  }
  for _, action := range info.Actions {
    t.SetDef(action.Title)
  }
}

func (info *ViewInfo) Tr(t *tr.Tr, lang string) {
  info.Title, _ = t.Tr(lang, info.Title)
  info.Description, _ = t.Tr(lang, info.Description)
  for _, column := range info.Table.Columns {
    column.Title, _ = t.Tr(lang, column.Title)
  }
  for _, action := range info.Actions {
    action.Title, _ = t.Tr(lang, action.Title)
  }
}

func (info *ViewInfo) Copy(in *ViewInfo) {
  copier.CopyWithOption(info, in, copier.Option{IgnoreEmpty: true, DeepCopy: true})
  // Use JSON
  //buf, _ := json.Marshal(in)
  //json.Unmarshal(buf, info)
}
