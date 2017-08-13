package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Used to handle requests to the web root "/"
// http.Request = data structure that represents the client HTTP request
// http.ResponseWriter assembles the HTTP servers response. By writing to it, we can send data to the HTTP client.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]		// path component of the requested URL
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)	// Write to w, the http.ResponseWriter to send data to HTTP client
}

// Defining the data structures
type Page struct {
	Title string
	Body []byte		// Byte slice - like an arry but not specified length. This is the type expected by the io libraries.
}

// Persistant storage
// Create a save method on Page
/* This is a method named save that takes as its receiver p, a pointer to Page. 
 * It takes no parameters, and returns a value of type error. */
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
	http.ListenAndServe(":8080", nil)
}