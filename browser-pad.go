package main
import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

//Page struct

type Page struct {
	Title string
	Body  []byte
}

//Save to file function

func (page *Page) save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}

//Load page from file

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//Handler for /view/

func viewHandler(writer http.ResponseWriter, reader *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(writer, reader, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

//Handler for /exit/ 

func editHandler(writer http.ResponseWriter, reader *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(writer, "edit", page)
}

//Handler for saving

func saveHandler(writer http.ResponseWriter, reader *http.Request, title string) {
	body := reader.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, reader, "/view/"+title, http.StatusFound)
}

//Template files in same folder

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//Creates page out of template page

func renderTemplate(writer http.ResponseWriter, template string, page *Page) {
	err := templates.ExecuteTemplate(writer, template + ".html", page)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

//Validation for url path

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Function that creates the handlers

func makeHandler(make func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		match := validPath.FindStringSubmatch(reader.URL.Path)
		if match == nil {
			http.NotFound(writer, reader)
			return
		}
		make(writer, reader, match[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}


