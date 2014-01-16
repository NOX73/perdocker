package perd


type Lang struct{
  Name      string
  Ext       string   
  Image     string
  Command   string
}

func (l *Lang) uniqFileName () string {
  return uniqFileName() + l.Ext
}

func (l *Lang) RunCommand (filePath string) string {
  return l.Command + " " + filePath
}

var Ruby *Lang = &Lang{"ruby", ".rb", "perdocker/ruby", "ruby"}
var Nodejs *Lang = &Lang{"nodejs", ".js", "perdocker/nodejs", "node"}
