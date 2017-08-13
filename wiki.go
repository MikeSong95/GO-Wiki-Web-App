package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// Defining the data structures
type Page struct {
	Title string
	Body []byte		// Byte slice - like an arry but with unspecified length. This is the type expected by the io libraries.
}

// Used to handle requests to the web root "/"
// http.Request = data structure that represents the client HTTP request
// http.ResponseWriter assembles the HTTP servers response. By writing to it, we can send data to the HTTP client.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]		// path component of the requested URL
	p, err := loadPage(title)
	// If the page is not found, redirects the client to the EDIT page so that the content may be created.
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)	// Adds an HTTP status code of http.StatusFound (302) and a Location header to the HTTP response.
		return
	}
	renderTemplate(w, "view", p)
}

// Loads page and displays and HTML form for editing the page
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	/* The renderTemplate function replaces the following code for better code re-use:
	 * t, _ := template.ParseFiles("edit.html")	
	 * t.Execute(w,p)
	 */ 
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

// Persistant storage
// Create a save method on Page
// This is a method named save that takes as its receiver p, a pointer to Page. 
// It takes no parameters, and returns a value of type error. 
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)	// 0600 = r-w permissions for current user only
}

// Load pages
// Caller to this function checks second parameter for error value.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)	// Use the blank identifier "_" to throw away unused return values.
	// Error handling. err == nil implies successfully loaded page.
	if err != nil {
		return nil, err 	// Error to be handled by the caller.
	}
	return &Page{Title: title, Body: body}, nil
}

// Main function
// Test what we've written
func main() {
	http.HandleFunc("/view/", viewHandler)	// Handle any requests under the path /view/
	http.HandleFunc("/edit/", editHandler)	// Handle any requests under the path /edit/
	http.ListenAndServe(":8080", nil)
}