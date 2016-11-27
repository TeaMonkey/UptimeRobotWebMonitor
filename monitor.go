package main

import (
  "fmt"
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
  key, _ := loadApiKey("apikey.txt")
  fmt.Fprintf(w, "<h1>Uptime Robot Monitors</h1><div>API key: %s</div>", key)
  fmt.Fprintf(w, "<a href=\"/list\">List all monitors</a>")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
  key, _ := loadApiKey("apikey.txt")
  t, _ := template.ParseFiles("Templates/ListMonitors.html")

  //Call UTR and get a list of monitors back
  response, err := http.Get("https://api.uptimerobot.com/getMonitors?apiKey=" + key)
  if err != nil {
    fmt.Println(err)
    return
      //log.Fatal(err)
  }
  defer response.Body.Close()
  responseData, _ := ioutil.ReadAll(response.Body)

  responseString := string(responseData)
  fmt.Println(responseString)

  //Parse this list into a data structure
  m := &Monitor{FriendlyName: "Bob", Status: 1}


  //Pass data structure to the template for display

  t.Execute(w, m)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("Templates/SingleMonitor.html")

  m := &Monitor{FriendlyName: "Bob", Status: 1}

  t.Execute(w, m)
}

func main() {
  http.HandleFunc("/", indexHandler)        //Index
  http.HandleFunc("/list/", listHandler)    //List of avail monitors
  http.HandleFunc("/view/", viewHandler)    //View a single monitor

  http.ListenAndServe(":8080", nil)
}
