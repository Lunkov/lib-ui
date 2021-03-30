package ui

import (
  "github.com/jinzhu/copier"
  "github.com/Lunkov/lib-tr"
)

type PageInfo struct {
  CODE           string                 `json:"code"              yaml:"code"`
  Title          string                 `json:"title"             yaml:"title"`
  Keywords       string                 `json:"keywords"          yaml:"keywords"`             // MAX 250
  Description    string                 `json:"description"       yaml:"description"`          // MAX 140
  Template       string                 `json:"template"          yaml:"template"`
  Cache          bool                   `json:"cache"             yaml:"cache"`
  MultiLanguage  bool                   `json:"multi_language"    yaml:"multi_language"`
}

func (info *PageInfo) TrDef(t *tr.Tr) {
  t.SetDef(info.Title)
  t.SetDef(info.Description)
  t.SetDef(info.Keywords)
}

func (info *PageInfo) Tr(t *tr.Tr, lang string) {
  info.Title, _       = t.Tr(lang, info.Title)
  info.Description, _ = t.Tr(lang, info.Description)
  info.Keywords, _    = t.Tr(lang, info.Keywords)
}

func (info *PageInfo) Copy(in *PageInfo) {
  copier.CopyWithOption(info, in, copier.Option{IgnoreEmpty: true, DeepCopy: true})
  // Use JSON
  //buf, _ := json.Marshal(in)
  //json.Unmarshal(buf, info)
}

