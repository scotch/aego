// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Modified from github.com/nathankerr/rest.go

Copyright (c) 2010 Nathan Kerr

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/
package rest

import (
	"net/http"
	"strings"
)

var resources = make(map[string]interface{})
var apiPrefix = "/-/api/v1"

// Lists all the items in the resource
// GET /resource
type index interface {
	Index(http.ResponseWriter, *http.Request)
}

// Creates a new resource item
// POST /resource
type create interface {
	Create(http.ResponseWriter, *http.Request)
}

// Views a resource item
// GET /resource/id
type find interface {
	Find(http.ResponseWriter, *http.Request, string)
}

// PUT /resource/id
type update interface {
	Update(http.ResponseWriter, *http.Request, string)
}

// DELETE /resource/id
type delete interface {
	Delete(http.ResponseWriter, *http.Request, string)
}

// Return options to use the service. If string is nil, then it is the base URL
// OPTIONS /resource/id
// OPTIONS /resource
type options interface {
	Options(http.ResponseWriter, *http.Request, string)
}

type Api struct {
	Prefix string
}

func NewApi(name string) *Api {
	return &Api{name}
}

// Add a resource route to http
func (a Api) Resource(prefix string, res interface{}) {
	apiPrefix = a.Prefix
	resources[prefix] = res
	url := apiPrefix + "/" + prefix
	http.Handle(url, http.HandlerFunc(rootHandler))
	http.Handle(url+"/", http.HandlerFunc(itemHandler))
}

// getResource takes a url and returns:
//  resource interface{}
//  id string
//  ok bool - returns false if the resource is missing.
func getResource(url string) (r interface{}, id string, false bool) {
	// Parse request URI to resource URI and (potential) ID
	path := strings.Split(url, apiPrefix)[1]
	s := strings.Split(path, "/")
	name := s[1]
	if len(s) > 2 {
		id = s[2]
	}
	r, ok := resources[name]
	return r, id, ok
}

// rootHandler handles urls without an id
func rootHandler(w http.ResponseWriter, r *http.Request) {
	resource, id, ok := getResource(r.URL.Path)
	if !ok {
		NotFound(w)
	}
	switch r.Method {
	case "GET":
		// Index
		if resIndex, ok := resource.(index); ok {
			resIndex.Index(w, r)
		} else {
			NotImplemented(w)
		}
	case "POST":
		// Create
		if resCreate, ok := resource.(create); ok {
			resCreate.Create(w, r)
		} else {
			NotImplemented(w)
		}
	case "OPTIONS":
		// automatic options listing
		if resOptions, ok := resource.(options); ok {
			resOptions.Options(w, r, id)
		} else {
			NotImplemented(w)
		}
	default:
		NotImplemented(w)
	}
}

// itemHandler handles urls containing an id
func itemHandler(w http.ResponseWriter, r *http.Request) {
	resource, id, ok := getResource(r.URL.Path)
	if !ok {
		NotFound(w)
	}
	switch r.Method {
	case "GET":
		// Find
		if resFind, ok := resource.(find); ok {
			resFind.Find(w, r, id)
		} else {
			NotImplemented(w)
		}
	case "PUT":
		// Update
		if resUpdate, ok := resource.(update); ok {
			resUpdate.Update(w, r, id)
		} else {
			NotImplemented(w)
		}
	case "DELETE":
		// Delete
		if resDelete, ok := resource.(delete); ok {
			resDelete.Delete(w, r, id)
		} else {
			NotImplemented(w)
		}
	case "OPTIONS":
		// automatic options
		if resOptions, ok := resource.(options); ok {
			resOptions.Options(w, r, id)
		} else {
			NotImplemented(w)
		}
	default:
		NotImplemented(w)
	}

}

// Emits a 404 Not Found
func NotFound(w http.ResponseWriter) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

// Emits a 501 Not Implemented
func NotImplemented(w http.ResponseWriter) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
}

// Emits a 201 Created with the URI for the new location
func Created(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	http.Error(w, "201 Created", http.StatusCreated)
}

// Emits a 200 OK with a location. Used when after a PUT
func Updated(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	http.Error(w, "200 OK", http.StatusOK)
}

// Emits a bad request with the specified instructions
func BadRequest(w http.ResponseWriter, instructions string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(instructions))
}

// Emits a 204 No Content
func NoContent(w http.ResponseWriter) {
	http.Error(w, "204 No Content", http.StatusNoContent)
}
