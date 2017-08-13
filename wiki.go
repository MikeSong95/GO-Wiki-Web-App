package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"errors"
)

// Defining the data structures
type Page struct {
	Title string
	Body []byte		// Byte slice - like an arry but with unspecified length. This is the type expected by the io libraries.
}

// Global variable to strore regexp
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$)")

// Global variable to parse all templates into a single template. Then use ExecuteTemplate to render a specific template.
// ParseFiles takes any number of string arguments (identify our template files) and parses them into templates that are named after the base file name.
// Prevents renderTemplate from parsing every time a page is rendered.
// template.Must is a convenience wrapper that panics when passed a non-nil error value, and otherwise returns the *Template unaltered.
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// Used to handle requests to the web root "/"
// http.Request = data structure that represents the client HTTP request
// http.ResponseWriter assembles the HTTP servers response. By writing to it, we can send data to the HTTP client.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)		// path component of the requested URL
	if err != nil {
		return
	}
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
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
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

// Handle the submission of forms located on the edit pages. 
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")	// Get the page content. It is of type string - we must convert it to []byte before it will fit into the Page struct.ß
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()	// Write the data to a file
	// An error that occurs during p.save() will be reported to the user
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Helper function to utilise HTTP templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)	// Sends a specified HTTP response code (Internal Server Error) and error message.
		return
	}
}

// Validates the path and extracts the page title
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)	// 404 not found http error
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression
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
	http.HandleFunc("/save/", saveHandler)	// Handle any requests under the path /save/
	http.ListenAndServe(":8080", nil)
}