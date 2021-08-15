package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Page represents a wiki page stored in memory
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filepath.Join("data", filename), p.Body, 0600)
}

var viewPath = filepath.Join("tmpl", "view.html")
var editPath = filepath.Join("tmpl", "edit.html")
var templates = template.Must(template.ParseFiles(viewPath, editPath))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filepath.Join("data", filename))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// ViewHandler serves view.html
func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// EditHandler serves edit.html
func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

// SaveHandler accepts POSTs to /save/
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
