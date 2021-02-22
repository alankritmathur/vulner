package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var urlName string
var structure []string

type CrawlData struct {
	Name string
	Site []string
}

var dat = new(CrawlData)

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		// http.ServeFile(w, r, "./static/index.html")
		template, _ := template.ParseFiles("./static/index.html")
		template.Execute(w, dat)

	case "POST":
		template, _ := template.ParseFiles("./static/index.html")
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		urlName = name
		crawler()
		template.Execute(w, dat)
		fmt.Println(dat)
		// fmt.Fprintf(w, "Name = %s\n", name)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func crawler() {
	resp, err := http.Get(urlName)
	if err != nil {
		log.Fatal(err)
	}
	getLinks(resp.Body)
}

//Collect all links from response body and return it as an array of strings
func getLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//todo: links list shoudn't contain duplicates
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}

				}
			}
			structure = links
			dat.Name = urlName
			dat.Site = links
		}
	}
}

func main() {
	http.HandleFunc("/", hello)
	fmt.Printf("Starting server for testing HTTP POST...\n")
	port, err := os.LookupEnv("PORT")
	if !err {
		port = "3000"
	}
	fmt.Println(port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	dat.Name = ""
	dat.Site = []string{}
}
