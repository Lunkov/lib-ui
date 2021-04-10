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

type Templates struct {
  templPath    string
  templates    map[string]*template.Template
  translate   *tr.Tr
  functions   *Functions
}

func NewTemplates(t *tr.Tr, p string) *Templates {
  return &Templates{
        templPath: p,
        templates: make(map[string]*template.Template),
        translate: t,
      }
}

func (t *Templates) SetFunc(f *Functions) {
  t.functions = f
}

func (t *Templates) FuncMap(name string, path string, style string, lang string) template.FuncMap {
  return t.functions.FuncMap(name, path, style, lang)
}

/*
func stdFuncMap() template.FuncMap {
  return template.FuncMap{
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
}*/
/*
func (t *Templates) AppendFuncMap(funcmap *template.FuncMap) {
  for k, v := range (*funcmap) {
    t.functions[k] = v
  }
}*/

// Get Name of Template from file name
func fileNameWithoutExtension(fileName string) string {
  return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func (ts *Templates) appendBaseTemplate(t *template.Template, name string, path string, style string, lang string) *template.Template {
  scanPath := fmt.Sprintf("%s/%s/%s/base/", ts.templPath, style, path)
  scanPath, err := filepath.Abs(scanPath)
  if err != nil {
    glog.Errorf("ERR: Get AbsPath(%s): %v", scanPath, err)
    return nil
  }
  count := 0
  errScan := filepath.Walk(scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      filebase := fileNameWithoutExtension(filename)
      if glog.V(2) {
        glog.Infof("LOG: Loading template(%s) file: %s\n", filebase, filename)
      }

      index := fmt.Sprintf("TEMPLATE#BASE#%s#%s#%s", style, filebase, lang)

      t_base, ok := ts.templates[index]
      if !ok {
        contents, err := ioutil.ReadFile(filename)
        if err != nil {
          glog.Errorf("ERR: Get Template(%s:%s): %v", filebase, filename, err)
          return err
        }
        t_base = template.New(filebase).Funcs(ts.FuncMap(name, path, style, lang))
        t_base, err = t_base.Parse(string(contents))
        if err != nil {
          glog.Errorf("ERR: Parse Template(%s:%s): %v", filebase, filename, err)
          if glog.V(9) {
            glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", filename, string(contents))
          }
          return err
        }
        ts.makeTrMap(t_base, lang)
      }
      count ++
      t.AddParseTree(t_base.Name(), t_base.Tree)
    }
    return nil
  })
  if glog.V(9) {
    glog.Infof("DBG: Scan Path: %s, Templates: %d", scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: %s\n", errScan)
  }
  return t
}


func (ts *Templates) Get(name string, path string, style string, lang string) *template.Template {
  index := fmt.Sprintf("TEMPLATE#%s#%s#%s", style, name, lang)
  
  i, ok := ts.templates[index]
  if ok {
    return i
  }
  var err error

  filen := fmt.Sprintf("%s/%s/%s/%s.html", ts.templPath, style, path, name)
  filename, err := filepath.Abs(filen)
  if err != nil {
    glog.Errorf("ERR: Get AbsPath(%s): %v", filen, err)
    return nil
  }
 
  contents, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: Get Template(%s): %v", filename, err)
    return nil
  }
  t := template.New(filename).Funcs(ts.FuncMap(name, path, style, lang))
  t, err = t.Parse(string(contents))
  if err != nil {
    glog.Errorf("ERR: Parse Template(%s): %v", filename, err)
    if glog.V(9) {
      glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", filename, string(contents))
    }
    return nil
  }
  t = ts.appendBaseTemplate(t, name, path, style, lang)
  if glog.V(9) {
    glog.Infof("DBG: Load Template(%s) file=%s", index, filename)
  }
  ts.templates[index] = t
  return t
}

func (t *Templates) getPrivate(name string, contents string, style string, lang string) *template.Template {
  index := fmt.Sprintf("PRIVTEMPLATE#%s#%s#%s#p", style, name, lang)
  
  i, ok := t.templates[index]
  if ok {
    return i
  }
  var err error

  tmp := template.New(index).Funcs(t.FuncMap(name, "private", style, lang))
  tmp, err = tmp.Delims("[[", "]]").Parse(contents)
  if err != nil {
    glog.Errorf("ERR: Parse Private Template(%s): %v", name, err)
    if glog.V(9) {
      glog.Infof("DBG: ERROR: Parse Template(%s) html=%s", name, contents)
    }
    return nil
  }
  t.templates[index] = tmp
  return tmp
}

func (ts *Templates) makeTrMap(t *template.Template, lang string) map[string]string {
  resTr := make(map[string]string)
  trs := findTrTemplate(t)
  for _, v := range trs {
    ts.translate.SetDef(v)
    resTr[v], _ = ts.translate.Tr(lang, v)
  }
  return resTr
}

// Extract the template vars required from *simple* templates.
// Only works for top level, plain variables. Returns all problematic parse.Node as errors.
func requiredTemplateVars(t *template.Template) ([]string, []error) {
  var res []string
  var errors []error
  var ln *parse.ListNode
  if t == nil {
    return res, errors
  }
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
