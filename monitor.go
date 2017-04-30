package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Monitors struct {
	Monitors []Monitor `xml:"monitors>monitor"`
}

type Monitor struct {
	Id                  int    `xml:"id"`
	FriendlyName        string `xml:"friendlyname"`
	Url                 string `xml:"url"`
	Type                int    `xml:"type"`
	SubType             string `xml:"subtype"`
	KeywordType         int    `xml:"keywordtype"`
	KeywordValue        int    `xml:"keywordvalue"`
	HttpUsername        string `xml:"httpusername"`
	HttpPassword        string `xml:"httppassword"`
	Port                int    `xml:"port"`
	Interval            int    `xml:"interval"`
	Status              int    `xml:"status"`
	AllTimeUpTimerRatio int64  `xml:"alltimeuptimeratio"`
}

func loadApiKey(fileName string) (string, error) {
	//This method adds a byte 10 to the end of the byte array for any file read
	//This 10 causes the http.Get call to throw a 'malformed HTTP status code "XXX'
	//error. Error in go 1.7.4 and 1.8.1, don't thik go 1.6 did this
	//Using the bufio.NewScanner thing as this works OK
	//key, err := ioutil.ReadFile(fileName)
	key := ""

	f, err := os.Open(fileName)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		key += scanner.Text()
	}
	defer f.Close()

	if err != nil {
		return "", err
	}

	// fmt.Println(key) // print the content as 'bytes'

	// str := string(key) // convert content to a 'string'

	// fmt.Println(str) // print the content as a 'string'

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

	//Call UTL and get a list of monitors back
	url := "https://api.uptimerobot.com/getMonitors?apiKey=" + key
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("HTTP GET error: ", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal("HTTP status error: ", response.StatusCode)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("HTTP/IO read body error: ", err)
		return
	}

	log.Println(string(data))

	// defer response.Body.Close()
	// responseData, _ := ioutil.ReadAll(response.Body)

	// responseString := string(responseData)
	// fmt.Println(responseString)

	//Parse this list into a data structure
	//m := &Monitor{FriendlyName: "Bob", Status: 1}
	m := Monitors{}
	if err := xml.Unmarshal(data, &m); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	} else if len(m.Monitors) != 0 {
		//Pass data structure to the template for display
		log.Fatal("HERE")
		t.Execute(w, m)
	}
	log.Fatal("Array size: ", len(m.Monitors))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Templates/SingleMonitor.html")

	m := &Monitor{FriendlyName: "Bob", Status: 1}

	t.Execute(w, m)
}

func main() {
	http.HandleFunc("/", indexHandler)     //Index
	http.HandleFunc("/list/", listHandler) //List of avail monitors
	http.HandleFunc("/view/", viewHandler) //View a single monitor

	http.ListenAndServe(":8080", nil)
}
