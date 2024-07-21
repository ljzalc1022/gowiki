package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var titlePattern = "[a-zA-Z0-9]+"

var pageLink = regexp.MustCompile("[" + titlePattern + "]")
var validPath = regexp.MustCompile("^/(edit|save|view)/(" + titlePattern + ")$")

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

func replaceHyperLink(pattern []byte) []byte {
	pageName := string(pattern[1 : len(pattern)-1])
	re := fmt.Sprintf("<a href=%q>%v</a>", "/view/"+pageName, pageName)
	return []byte(re)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page, replace bool) {
	var renderedPage strings.Builder
	err := templates.ExecuteTemplate(&renderedPage, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderedPageBytes := []byte(renderedPage.String())

	// replacing the [PageName] with <a></a>
	if replace {
		renderedPageBytes = pageLink.ReplaceAllFunc(renderedPageBytes, replaceHyperLink)
	}

	fmt.Fprintf(w, "%s", renderedPageBytes)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p, true)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p, false)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

const homePage = "/view/FrontPage"

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, homePage, http.StatusFound)
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	http.HandleFunc("/{$}", homePageHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
