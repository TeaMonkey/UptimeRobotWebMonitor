package main

import (
  //"fmt"
  "io/ioutil"
  "net/http"
  "html/template"
)

type Monitor struct {
  Id  int
  FriendlyName  string
  Url string
  Type int
  SubType string
  KeywordType int
  KeywordValue int
  HttpUsername string
  HttpPassword string
  Port int
  Interval int
  Status int
  AllTimeUpTimerRatio int64
}

func loadApiKey(fileName string) (string, error) {
  key, err := ioutil.ReadFile(fileName)
  if err != nil {
    return "", err
  }
  return string(key), nil
}

//func viewHandler(w http.ResponseWriter, r *http.Request) {
//  key, _ := loadApiKey("apikey.txt")
//  fmt.Fprintf(w, "<h1>Uptime Robot Monitors</h1><div>API key: %s</div>", key)
//}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("Templates/SingleMonitor.html")

  m := &Monitor{FriendlyName: "Bob", Status: 1}

  t.Execute(w, m)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.ListenAndServe(":8080", nil)
}
