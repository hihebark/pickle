package core

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type ServeMux struct {
	mutex sync.RWMutex
	md    Markdown
}

//Markdown
type Markdown struct {
	Name string
	Data string
}

const GITHUBAPI string = "https://api.github.com/markdown/raw"

//MarkdowntoHTML convert given markdown data to html.
func MarkdowntoHTML(data string) string {
	req, err := http.NewRequest("POST", GITHUBAPI, bytes.NewBufferString(data))
	if err != nil {
		fmt.Printf("! Error on request\n\t\t%v\n", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("! Error on response\n\t\t%v\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

//StartServer open the port 7069.
func StartServer(md Markdown) {
	fmt.Println("+ Stating server on: localhost:7069 | [::1]:7069")
	mux := &ServeMux{md: md}
	err := http.ListenAndServe(":7069", mux)
	if err != nil {
		fmt.Printf("! Error on starting server\n\t\t%v\n", err)
	}
}

//ServeHTTP hundle results route
func (mutex *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		mutex.mutex.RLock()
		defer mutex.mutex.RUnlock()
		showResults(w, r, mutex.md)
		return
	case strings.Contains(r.URL.Path, "css/"):
		assets := []string{
			"css/pickle.css",
			"css/syntax.css",
			"css/github.css",
		}
		for _, v := range assets {
			mutex.mutex.RLock()
			defer mutex.mutex.RUnlock()
			http.ServeFile(w, r, "template/assets/"+v)
		}
		return
	default:
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

//ShowResultsFile see if the directory is in the data/results and show all json files in it
func showResults(w http.ResponseWriter, r *http.Request, data Markdown) {
	htmlTemplate := template.New("index.html")
	htmlTemplate, err := htmlTemplate.ParseFiles("template/index.html")
	if err != nil {
		fmt.Printf("! Error html parser\n\t\t%v\n", err)
	}
	data.Data = html.UnescapeString(data.Data)
	htmlTemplate.Execute(w, data)
}
