package ui

import (
  "fmt"
  "strings"
  "os"
  "path/filepath"
  "io/ioutil"
  "html/template"
  "text/template/parse"
  
  "github.com/golang/glog"
  "github.com/Lunkov/lib-tr"
)

var memTemplate = make(map[string]*template.Template)
var fmap_custom = make(template.FuncMap)

func AppendFuncMap(funcmap template.FuncMap) {
  fmap_custom = funcmap
}

func funcMap(name string, path string, style string, lang string) template.FuncMap {
  fmap := template.FuncMap{
                  "TR": func(str string) string {
                      t, _ := tr.Tr(lang, str)
                      return t
                  },
                  "TR_LANG": func() string {
                      return lang
                  },
                  "TR_LANG_NAME": func() string {
                      return tr.LangName(lang)
                  },
                  "FORM": func(form_code string, data map[string]interface{}) template.HTML {
                      return template.HTML(renderForm(lang, form_code, style, false, &data))
                  },
                  "MODAL": func(form_code string, data map[string]interface{}) template.HTML {
                      return template.HTML(renderForm(lang, form_code, style, true, &data))
                  },
                  "VIEW": func(view_code string, data map[string]interface{}) template.HTML {
                      res, ok := renderView(lang, view_code, style, &data)
                      if ok {
                        return template.HTML(res)
                      }
                      return template.HTML("--- VIEW '" + view_code + "' NOT FOUND ---")
                  },
                  "attr":func(s string) template.HTMLAttr{
                      return template.HTMLAttr(s)
                  },
                  "safe": func(s string) template.HTML {
                      return template.HTML(s)
                  },
                  "url": func(s string) template.URL {
                      return template.URL(s)
                  },
                }
  for k, v := range fmap_custom {
    fmap[k] = v
  }
  return fmap
}

// Get Name of Template from file name
func fileNameWithoutExtension(fileName string) string {
  return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func appendBaseTemplate(t *template.Template, name string, path string, style string, lang string) *template.Template {
  scanPath := fmt.Sprintf("./templates/%s/%s/base/", style, path)
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      filebase := fileNameWithoutExtension(filename)
      if glog.V(2) {
        glog.Infof("LOG: Loading template(%s) file: %s\n", filebase, filename)
      }

      index := fmt.Sprintf("TEMPLATE#BASE#%s#%s#%s", style, filebase, lang)

      t_base, ok := memTemplate[index]
      if !ok {
        contents, err := ioutil.ReadFile(filename)
        if err != nil {
          glog.Errorf("ERR: Get Template(%s:%s): %v", filebase, filename, err)
          return err
        }
        t_base = template.New(filebase).Funcs(funcMap(name, path, style, lang))
        t_base, err = t_base.Parse(string(contents))
        if err != nil {
          glog.Errorf("ERR: Parse Template(%s:%s): %v", filebase, filename, err)
          if glog.V(9) {
            glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", filename, string(contents))
          }
          return err
        }
        makeTrMap(t_base, lang)
      }
      count ++
      t.AddParseTree(t_base.Name(), t_base.Tree)
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Templates: %d\n", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: %s\n", errScan)
  }
  return t
}


func getTemplate(name string, path string, style string, lang string) *template.Template {
  index := fmt.Sprintf("TEMPLATE#%s#%s#%s", style, name, lang)
  
  i, ok := memTemplate[index]
  if ok {
    return i
  }
  var err error

  filename := fmt.Sprintf("./templates/%s/%s/%s.html", style, path, name)
 
  contents, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: Get Template(%s): %v", filename, err)
    return nil
  }
  t := template.New(filename).Funcs(funcMap(name, path, style, lang))
  t, err = t.Parse(string(contents))
  if err != nil {
    glog.Errorf("ERR: Parse Template(%s): %v", filename, err)
    if glog.V(9) {
      glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", filename, string(contents))
    }
    return nil
  }
  t = appendBaseTemplate(t, name, path, style, lang)
  
  memTemplate[index] = t
  return t
}

func getPrivateTemplate(name string, contents string, style string, lang string) *template.Template {
  index := fmt.Sprintf("PRIVTEMPLATE#%s#%s#%s#p", style, name, lang)
  
  i, ok := memTemplate[index]
  if ok {
    return i
  }
  var err error

  t := template.New(index).Funcs(funcMap(name, "mem", style, lang))
  t, err = t.Delims("[[", "]]").Parse(contents)
  if err != nil {
    glog.Errorf("ERR: Parse Private Template(%s): %v", name, err)
    if glog.V(9) {
      glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", name, contents)
    }
    return nil
  }
  memTemplate[index] = t
  return t
}

// UNION MAPS
func unionMap(srcMap *map[string]interface{}, newMap *map[string]interface{}) {
  if srcMap != nil && newMap != nil {
    for k, v := range (*newMap) {
      (*srcMap)[k] = v
    }
  }
}

func unionMapStr(srcMap *map[string]interface{}, newMap *map[string]string) {
  if srcMap != nil && newMap != nil {
    for k, v := range (*newMap) {
      (*srcMap)[k] = v
    }
  }
}

func makeTrMap(t *template.Template, lang string) map[string]string {
  resTr := make(map[string]string)
  trs := findTrTemplate(t)
  for _, v := range trs {
    tr.SetDef(v)
    resTr[v], _ = tr.Tr(lang, v)
  }
  return resTr
}

// Extract the template vars required from *simple* templates.
// Only works for top level, plain variables. Returns all problematic parse.Node as errors.
func requiredTemplateVars(t *template.Template) ([]string, []error) {
  var res []string
  var errors []error
  var ln *parse.ListNode
  ln = t.Tree.Root
Node:
  for _, n := range ln.Nodes {
    if nn, ok := n.(*parse.ActionNode); ok {
      p := nn.Pipe
      if len(p.Decl) > 0 {
        errors = append(errors, fmt.Errorf("len(p.Decl): Node %v not supported", n))
        continue Node
      }
      for _, c := range p.Cmds {
        if len(c.Args) != 1 {
          errors = append(errors, fmt.Errorf("len(c.Args)=%d: Node %v not supported", len(c.Args), n))
          continue Node
        }
        if a, ok := c.Args[0].(*parse.FieldNode); ok {
          if len(a.Ident) != 1 {
              errors = append(errors, fmt.Errorf("len(a.Ident): Node %v not supported", n))
              continue Node
          }
          res = append(res, a.Ident[0])
        } else {
          errors = append(errors, fmt.Errorf("parse.FieldNode: Node %v not supported", n))
          continue Node
        }

      }
    } else {
      if _, ok := n.(*parse.TextNode); !ok {
        errors = append(errors, fmt.Errorf("parse.TextNode: Node %v not supported", n))
        continue Node
      }
    }
  }
  return res, errors
}

// Extract the template vars required from *simple* templates.
// Only works for top level, plain variables. Returns all problematic parse.Node as errors.
func findTrTemplate(t *template.Template) []string {
  var res []string
  if t == nil || t.Tree == nil  || t.Tree.Root == nil {
    return res
  }
  var ln *parse.ListNode
  ln = t.Tree.Root
Node:
  for _, n := range ln.Nodes {
    if nn, ok := n.(*parse.ActionNode); ok {
      p := nn.Pipe
      for _, c := range p.Cmds {
        if len(c.Args) == 2 {
          if c.Args[0].String() == "TR" {
            str := strings.ReplaceAll(c.Args[1].String(), "\"", "")
            str = strings.ReplaceAll(str, "'", "")
            res = append(res, str)
          } else {
            continue Node
          }
        }
      }
    }
  }
  return res
}
