package ui

import (
  "sync"
  "time"
  "regexp"
  "fmt"
  "io"
  "errors"
  "strings"
  
  "path/filepath"
  "encoding/json"
  
  "github.com/golang/glog"
  "github.com/radovskyb/watcher"

  "net/http"
  "github.com/golang/gddo/httputil/header"
  
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
  "github.com/tdewolff/minify/v2/svg"
  
  "github.com/Lunkov/lib-tr"
  "github.com/Lunkov/lib-cache"
)

type UIInfo struct {
  Form   map[string]FormInfo   `json:"forms"   yaml:"forms"`
  View   map[string]ViewInfo   `json:"views"   yaml:"views"`
}

type UIStat struct {
  CntPage    int64     `json:"cnt_page"`
  CntForm    int64     `json:"cnt_form"`
  CntView    int64     `json:"cnt_view"`
  CntTr      int       `json:"cnt_tr"`
  CntLang    int       `json:"cnt_lang"`
}

type UI struct {
  watcherFiles    *watcher.Watcher
  uimu             sync.RWMutex
  mapJSONs         map[string][]byte // is string
  mapRenders       cache.ICache
  configTrPath     string
  minifyRender    *minify.M
  t               *tr.Tr
  forms           *Forms
  views           *Views
  pages           *Pages
  functions       *Functions
  templates       *Templates
}

func NewUI(templPath string, cfgForms *cache.CacheConfig, cfgViews *cache.CacheConfig, cfgPages *cache.CacheConfig, cfgRenders *cache.CacheConfig) (*UI) {
  t := tr.New()
  tmplts := NewTemplates(t, templPath)
  forms := NewForms(cfgForms, t, tmplts)
  views := NewViews(cfgViews, t, tmplts)
  funcs := NewFunctions(t, forms, views)
  tmplts.SetFunc(funcs)
  return &UI{
         forms      : forms,
         views      : views,
         pages      : NewPages(cfgPages, t, tmplts),
         mapJSONs   : make(map[string][]byte), // is string
         mapRenders : cache.NewConfig(cfgRenders),
         t          : t,
         functions  : funcs,
         templates  : tmplts,
  }
}

func (u *UI) ToJSON(role_name string, lang string, menus []string, forms []string, views []string) []byte {
  index := role_name + "#" + lang
  ijson, ok := u.mapJSONs[index]
  if ok {
    return ijson
  }
  var tui UIInfo
  tui.Form = u.forms.Filter(lang, forms)
  tui.View = u.views.Filter(lang, views)
  resJSON, _ := json.Marshal(tui)
  u.uimu.Lock()
  u.mapJSONs[index] = resJSON
  u.uimu.Unlock()
  return resJSON
}

func (u *UI) Init(configPath string, enableWatcher bool, enableMinify bool) {
  u.minifyRender = nil
  if enableMinify {
    u.minifyRender = minify.New()
    u.minifyRender.AddFunc("text/css", css.Minify)
    u.minifyRender.AddFunc("text/html", html.Minify)
    u.minifyRender.AddFunc("image/svg+xml", svg.Minify)
    u.minifyRender.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
    glog.Infof("LOG: Enable Minify HTML")
  }
  
  u.loadAll(configPath)
  
  if enableWatcher {
    u.watcherFiles = watcher.New()
    u.watcherFiles.SetMaxEvents(1)
    u.watcherFiles.FilterOps(watcher.Rename, watcher.Move, watcher.Remove, watcher.Create, watcher.Write)
    go func() {
      for {
        select {
        case event := <-u.watcherFiles.Event:	
          if glog.V(9) {
            glog.Infof("DBG: Watcher Event: %v", event)
          }
          // Ignore New Translate Files
          if filepath.Ext(event.Name()) != ".!yaml" {
            u.loadAll(configPath)
          }
        case err := <-u.watcherFiles.Error:
          glog.Fatalf("ERR: Watcher Event: %v", err)
        case <-u.watcherFiles.Closed:
          glog.Infof("LOG: Watcher Close")
          return
        }
      }
    }()
    // Start the watching process - it'll check for changes every 100ms.
    glog.Infof("LOG: Watcher Start (%s)", configPath)
    if err := u.watcherFiles.AddRecursive(configPath); err != nil {
      glog.Fatalf("ERR: Watcher AddRecursive: %v", err)
    }
    if err := u.watcherFiles.AddRecursive("./templates/"); err != nil {
      glog.Fatalf("ERR: Watcher AddRecursive: %v", err)
    }
    
    if glog.V(9) {
      // Print a list of all of the files and folders currently
      // being watched and their paths.
      for path, f := range u.watcherFiles.WatchedFiles() {
        glog.Infof("DBG: WATCH FILE: %s: %s\n", path, f.Name())
      }
    }
	  go func() {
      if err := u.watcherFiles.Start(time.Millisecond * 100); err != nil {
        glog.Fatalf("ERR: Watcher Start: %v", err)
      }
    }()
  }
}

func (u *UI) GetLangList() *map[string]map[string]string {
  return u.t.GetList()
}

func (u *UI) makeMimiHTML(s string) string {
  if u.minifyRender != nil {
    res, err := u.minifyRender.String("text/html", s)
    if err != nil {
      glog.Errorf("ERR: HTML Minify: %v", err)
    } else {
      return res
    }
  }
  return s
}

func (u *UI) loadAll(configPath string) {
  // Clear cache
  u.mapJSONs = make(map[string][]byte) // is string
  u.mapRenders.Clear()

  u.configTrPath = configPath + "/tr"

  u.t.LoadLangs(configPath + "/langs.yaml")
  u.t.LoadTrs(u.configTrPath)
  u.forms.Load(configPath)
  u.views.Load(configPath)
  u.pages.Load(configPath)

  u.t.SaveNew(u.configTrPath)
  
  glog.Infof("LOG: Load Langs: %d", u.t.LangCount())
  glog.Infof("LOG: Load Tr: %d",    u.t.Count())
  glog.Infof("LOG: Load Forms: %d", u.forms.Count())
  glog.Infof("LOG: Load Views: %d", u.views.Count())
  glog.Infof("LOG: Load Pages: %d", u.pages.Count())
}

func (u *UI) GetStat() *UIStat {
  return &UIStat{
                CntPage: u.pages.Count(),
                CntForm: u.forms.Count(),
                CntView: u.views.Count(),
                CntTr:   u.t.Count(),
                CntLang: u.t.LangCount(),
  }
}

func (u *UI) RenderForm(form_code string, lang string, style string, isModal bool, data *map[string]interface{}) string {
  index := u.forms.index(lang, form_code, style, isModal, data)
  var render string
  rForm, ok := u.mapRenders.Get(index, &render)
  if ok {
    render, _ = rForm.(string)
    return render
  }
  render = u.makeMimiHTML(u.forms.Render(form_code, lang, style, isModal, data))
  u.mapRenders.Set(index,  render)
  u.t.SaveNew(u.configTrPath)
  return render
}

func (u *UI) RenderPage(page_code string, lang string, style string, private bool, data *map[string]interface{}) string {
  index := u.pages.index(lang, page_code, style, private)
  var render string
  rPage, ok := u.mapRenders.Get(index, &render)
  if ok {
    render, _ = rPage.(string)
    return render
  }
  rend, ok := u.pages.Render(page_code, lang, style, private, data)
  if ok {
    render = u.makeMimiHTML(rend)
    u.mapRenders.Set(index, render)
    defer u.t.SaveNew(u.configTrPath)
    return render
  }
  return "!! ERROR: RENDER PAGE(" + page_code + ") !!"
}

type malformedRequest struct {
  status int
  msg    string
}

func (mr *malformedRequest) Error() string {
  return mr.msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
  if r.Header.Get("Content-Type") != "" {
    value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
    if value != "application/json" {
      msg := "Content-Type header is not application/json"
      return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
    }
  }

  r.Body = http.MaxBytesReader(w, r.Body, 1048576)

  dec := json.NewDecoder(r.Body)
  // TODO
  // dec.DisallowUnknownFields()

  err := dec.Decode(&dst)
  if err != nil {
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError

    switch {
    case errors.As(err, &syntaxError):
        msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.ErrUnexpectedEOF):
        msg := fmt.Sprintf("Request body contains badly-formed JSON")
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.As(err, &unmarshalTypeError):
        msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case strings.HasPrefix(err.Error(), "json: unknown field "):
        fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
        msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.EOF):
        msg := "Request body must not be empty"
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case err.Error() == "http: request body too large":
        msg := "Request body must not be larger than 1MB"
        return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

    default:
        return err
    }
  }

  err = dec.Decode(&struct{}{})
  if err != io.EOF {
    msg := "Request body must only contain a single JSON object"
    return &malformedRequest{status: http.StatusBadRequest, msg: msg}
  }

  return nil
}
