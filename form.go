package ui

import (
  "github.com/jinzhu/copier"
  
  "github.com/Lunkov/lib-tr"
)

// Form methods: POST or GET.
const (
  POST = "POST"
  GET  = "GET"
)

type InfoItem struct {
  Description  string      `json:"description"   yaml:"description"`
  Enum         []string    `json:"enum"          yaml:"enum"`
}

type InfoProperty struct {
  CODE           string       `json:"code"          yaml:"code"`
  Order          int          `json:"order"         yaml:"order"`
  Title          string       `json:"title"         yaml:"title"`
  Description    string       `json:"description"   yaml:"description"`
  DataType       string       `json:"data_type"     yaml:"data_type"`
  WidgetType     string       `json:"widget"        yaml:"widget"`
  Required       string       `json:"required"      yaml:"required"`
  Format         string       `json:"format"        yaml:"format"`
  Table          string       `json:"table"         yaml:"table"`
  Default        string       `json:"default"       yaml:"default"`
  OneOf          []InfoItem   `json:"oneOf"         yaml:"oneOf"`
  OneOfConst     string       `json:"oneOfConst"    yaml:"oneOfConst"`
}

type InfoButton struct {
  CODE           string       `json:"code"          yaml:"code"`
  Order          int          `json:"order"         yaml:"order"`
  Title          string       `json:"title"         yaml:"title"`
  Type           string       `json:"type"          yaml:"type"`
}

type InfoTab struct {
  CODE         string    `json:"code"         yaml:"code"`
  Order        int       `json:"order"        yaml:"order"`
  Title        string    `json:"title"        yaml:"title"`
  Description  string    `json:"description"  yaml:"description"`
  Fields       []string  `json:"fields"       yaml:"fields"`
}

type ArProperties []InfoProperty
type ArButtons    []InfoButton
type ArTabs       []InfoTab

type FormInfo struct {
  Code          string                     `json:"code"         yaml:"code"`
  EncType       string                     `json:"enctype"      yaml:"enctype"`
  Method        string                     `json:"method"       yaml:"method"`
  Action        string                     `json:"action"       yaml:"action"`
  Version       int32                      `json:"version"      yaml:"version"`
  Date          string                     `json:"date"         yaml:"date"`
  Date_start    string                     `json:"date_start"   yaml:"date_start"`
  Date_finish   string                     `json:"date_finish"  yaml:"date_finish"`
  Title         string                     `json:"title"        yaml:"title"`
  Type          string                     `json:"type"         yaml:"type"`
  Description   string                     `json:"description"  yaml:"description"`
  Tabs          ArTabs                     `json:"tabs"         yaml:"tabs,flow"`
  Properties    ArProperties               `json:"properties"   yaml:"properties,flow"`
  Buttons       ArButtons                  `json:"buttons"      yaml:"buttons,flow"`
}

func (info *FormInfo) TrDef(t *tr.Tr) {
  t.SetDef(info.Title)
  t.SetDef(info.Description)
  for _, tab := range info.Tabs {
    t.SetDef(tab.Title)
    t.SetDef(tab.Description)
  }
  for _, prop := range info.Properties {
    t.SetDef(prop.Title)
    t.SetDef(prop.Description)
  }
  for _, button := range info.Buttons {
    t.SetDef(button.Title)
  }
}

func (info *FormInfo) Tr(t *tr.Tr, lang string) {
  info.Title, _ = t.Tr(lang, info.Title)
  info.Description, _ = t.Tr(lang, info.Description)
  for i, tab := range info.Tabs {
    tab.Title, _ = t.Tr(lang, tab.Title)
    tab.Description, _ = t.Tr(lang, tab.Description)
    info.Tabs[i] = tab
  }
  for i, prop := range info.Properties {
    prop.Title, _ = t.Tr(lang, prop.Title)
    prop.Description, _ = t.Tr(lang, prop.Description)
    info.Properties[i] = prop
  }
  for i, button := range info.Buttons {
    button.Title, _ = t.Tr(lang, button.Title)
    info.Buttons[i] = button
  }
}

func (info *FormInfo) Copy(in *FormInfo) {
  copier.CopyWithOption(info, in, copier.Option{IgnoreEmpty: true, DeepCopy: true})
  // Use JSON
  // buf, _ := json.Marshal(in)
  // json.Unmarshal(buf, info)
}


