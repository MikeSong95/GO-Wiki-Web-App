package main

import (
	"fmt"
	"io/ioutil"
)

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
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}