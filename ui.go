package ui

import (
  "sync"
  "time"
  "regexp"
  "encoding/json"
  
  "github.com/golang/glog"
  "github.com/radovskyb/watcher"
  
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
  "github.com/tdewolff/minify/v2/svg"
  
  "github.com/Lunkov/lib-tr"
)

type UIInfo struct {
  Form   map[string]FormInfo   `json:"forms"   yaml:"forms"`
  View   map[string]ViewInfo   `json:"views"   yaml:"views"`
}

type UIStat struct {
  CntPage    int     `json:"cnt_page"`
  CntForm    int     `json:"cnt_form"`
  CntView    int     `json:"cnt_view"`
  CntTr      int     `json:"cnt_tr"`
  CntLang    int     `json:"cnt_lang"`
}

var watcherFiles watcher.Watcher
var uimu sync.RWMutex
var mapJSONs = make(map[string][]byte) // is string
var mapRenders = make(map[string]string)
var configTrPath string
var minifyRender *minify.M

func ToJSON(role_name string, lang string, menus []string, forms []string, views []string) []byte {
  index := role_name + "#" + lang
  ijson, ok := mapJSONs[index]
  if ok {
    return ijson
  }
  var tui UIInfo
  tui.Form = formFilter(lang, forms)
  tui.View = viewFilter(lang, views)
  resJSON, _ := json.Marshal(tui)
  uimu.Lock()
  mapJSONs[index] = resJSON
  uimu.Unlock()
  return resJSON
}

func Init(configPath string, enableWatcher bool, enableMinify bool) {
  minifyRender = nil
  if enableMinify {
    minifyRender = minify.New()
    minifyRender.AddFunc("text/css", css.Minify)
    minifyRender.AddFunc("text/html", html.Minify)
    minifyRender.AddFunc("image/svg+xml", svg.Minify)
    minifyRender.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
    glog.Infof("LOG: Enable Minify HTML")
  }
  
  loadAll(configPath)
  
  if enableWatcher {
    watcherFiles := watcher.New()
    watcherFiles.SetMaxEvents(1)
    watcherFiles.FilterOps(watcher.Rename, watcher.Move, watcher.Remove, watcher.Create, watcher.Write)
    go func() {
      for {
        select {
        case event := <-watcherFiles.Event:	
          glog.Infof("DBG: Watcher Event: %v", event)
          // Ignore New Translate Files
          if event.Name() != "!tr_new.yaml" {
            loadAll(configPath)
          }
        case err := <-watcherFiles.Error:
          glog.Fatalf("ERR: Watcher Event: %v", err)
        case <-watcherFiles.Closed:
          glog.Infof("LOG: Watcher Close")
          return
        }
      }
    }()
    // Start the watching process - it'll check for changes every 100ms.
    glog.Infof("LOG: Watcher Start (%s)", configPath)
    if err := watcherFiles.AddRecursive(configPath); err != nil {
      glog.Fatalf("ERR: Watcher AddRecursive: %v", err)
    }
    if err := watcherFiles.AddRecursive("./templates/"); err != nil {
      glog.Fatalf("ERR: Watcher AddRecursive: %v", err)
    }
    
    if glog.V(9) {
      // Print a list of all of the files and folders currently
      // being watched and their paths.
      for path, f := range watcherFiles.WatchedFiles() {
        glog.Infof("DBG: WATCH FILE: %s: %s\n", path, f.Name())
      }
    }
	  go func() {
      if err := watcherFiles.Start(time.Millisecond * 100); err != nil {
        glog.Fatalf("ERR: Watcher Start: %v", err)
      }
    }()
  }
}

func makeMimiHTML(s string) string {
  if minifyRender != nil {
    res, err := minifyRender.String("text/html", s)
    if err != nil {
      glog.Errorf("ERR: HTML Minify: %v", err)
    } else {
      return res
    }
  }
  return s
}

func loadAll(configPath string) {
  // Clear cache
  mapJSONs = make(map[string][]byte) // is string
  mapRenders = make(map[string]string)

  tr.LoadLangs(configPath + "/langs.yaml")
  tr.LoadTrs(configPath + "/tr")
  formInit(configPath + "/forms")
  viewInit(configPath + "/views")
  pageInit(configPath + "/pages")

  configTrPath = configPath + "/tr"
  tr.SaveNew(configTrPath)
  
  glog.Infof("LOG: Load Langs: %d", tr.LangCount())
  glog.Infof("LOG: Load Tr: %d",    tr.Count())
  glog.Infof("LOG: Load Forms: %d", formCount())
  glog.Infof("LOG: Load Views: %d", viewCount())
  glog.Infof("LOG: Load Pages: %d", pageCount())
}

func GetStat() UIStat {
  return UIStat{
                CntPage: pageCount(),
                CntForm: formCount(),
                CntView: viewCount(),
                CntTr:   tr.Count(),
                CntLang: tr.LangCount(),
  }
}

func RenderForm(lang string, form_code string, style string, isModal bool, data *map[string]interface{}) string {
  index := indexForm(lang, form_code, style, isModal, data)
  uimu.RLock()
  rForm, ok := mapRenders[index]
  uimu.RUnlock()
  if ok {
    return rForm
  }
  render := renderForm(lang, form_code, style, isModal, data)
  uimu.Lock()
  mapRenders[index] = makeMimiHTML(render)
  tr.SaveNew(configTrPath)
  uimu.Unlock()
  return mapRenders[index]
}

func RenderPage(lang string, page_code string, style string, private bool, data *map[string]interface{}) string {
  index := indexPage(lang, page_code, style, private)
  uimu.RLock()
  rPage, ok := mapRenders[index]
  uimu.RUnlock()
  if ok {
    return rPage
  }
  render, ok := renderPage(lang, page_code, style, private, data)
  if ok {
    uimu.Lock()
    mapRenders[index] = makeMimiHTML(render)
    tr.SaveNew(configTrPath)
    uimu.Unlock()
    return mapRenders[index]
  }
  return "!! ERROR: RENDER PAGE(" + page_code + ") !!"
}

