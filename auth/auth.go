// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/auth provides multi-provider Authentication.

Example Usage:

  import (
    "github.com/scotch/hal/auth"
    "github.com/scotch/hal/auth/google"
  )

  // Register the Google Provider.
  googleProvider := google.Provider.New("12345", "ABCD")
  Register("google", &googleProvider)
  // Register additional providers.
  // ...


*/

package auth

import (
	aeuser "appengine/user"
	"github.com/scotch/hal/auth/profile"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/user"
	"net/http"
	"strings"
)

var (
	// BaseURL represents the base url to be used for providers. For
	// example if the base url is /auth/ all provider urls would be at
	// /auth/<provider name>
	BaseURL = "/-/auth/"
	// LoginURL is a string representing the URL to be redirected to on
	// errors.
	LoginURL = "/-/auth/login"
	// LogoutURL is a string representing the URL to be used to remove
	// the auth cookie.
	LogoutURL = "/-/auth/logout"
	// SuccessURL is a string representing the URL to be direct to on a
	// successful login.
	SuccessURL = "/"
)

var providers = make(map[string]authenticater)

type authenticater interface {
	Authenticate(http.ResponseWriter, *http.Request, *profile.Profile) (string, error)
}

// Register adds an Authenticater for the auth service.
//
// It takes a string which is used for the url, and a pointer to an
// authentication provider that implements Authenticater.
// E.g.
//
//   googleProvider := google.Provider.New("12345", "ABCD")
//   Register("google", &googleProvider)
//
func Register(key string, auth authenticater) {
	providers[key] = auth
	// Set the start url e.g. /-/auth/google to be handled by the handler.
	http.HandleFunc(BaseURL+key, handler)
	// Set the callback url e.g. /-/auth/google/callback to be handled by the handler.
	http.HandleFunc(BaseURL+key+"/callback", handler)
}

// breakURL parse an url and returns the provider key. If the URL is
// invalid it returns and empty string "".
func breakURL(url string) (name string) {
	if p := strings.Split(url, BaseURL); len(p) > 1 {
		name = strings.Split(p[1], "/")[0]
	}
	return
}

// CreateAndLogin does the following:
//
//  - Search for an existing user - session -> Profile -> email address
//  - Saves the Profile to the datastore
//  - Creates a User or appends the AuthID to the Requesting user's account
//  - Logs in the User
//  - Adds the admin role to the User if they are an GAE Admin.
func CreateAndLogin(w http.ResponseWriter, r *http.Request,
	up *profile.Profile) (u *user.User, err error) {

	var id string
	var idUP string
	var idSess string
	var saveUser bool

	c := context.NewContext(r)
	up.SetKey(c)
	// Check the session for a UserID
	idSess, _ = user.CurrentUserID(r)
	// Check for an existing Profile
	up2 := profile.New()
	if err := profile.Get(c, up.Key.StringID(), up2); err == nil {
		idSess = up2.UserID
	}
	if idUP != "" || idSess != "" {
		if idUP == idSess {
			id = idSess
		} else {
			// TWO USER ACCOUNTS FOR 1 USEPROFILE
			// TODO implement some type of user merge here.
			// for the time being use the logged in User's ID
			id = idSess
		}
		// Get the user
		if u, err = user.Get(c, id); err != nil {
			// if user is not found we have some type of syncing problem.
			c.Criticalf(`auth: userID: %v was saved to Profile / Session, but was not found in the datastore`, id)
			return
		}
	} else {
		// New user
		u = user.New()
		if err = u.SetKey(c); err != nil {
			return
		}
		saveUser = true
	}
	// Add AuthID
	if u.AddAuthID(up.Key.StringID()) {
		saveUser = true
	}
	// If current user is an admin in GAE add role to User
	if aeuser.IsAdmin(c) {
		// Save the roll to the session
		_ = user.CurrentUserSetRole(w, r, "admin", true)
		// Add the role to the user's roles.
		if u.AddRole("admin") {
			saveUser = true
		}
	}
	// Log the user in.
	_ = user.CurrentUserSetID(w, r, u.Key.StringID())
	if saveUser {
		err = u.Put(c)
	}
	// TODO should the Profile always be saved?
	up.UserID = u.Key.StringID()
	err = up.Put(c)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	up := profile.New()
	k := breakURL(r.URL.Path)
	p := providers[k]
	url, err := p.Authenticate(w, r, up)
	if err != nil {
		// TODO: set error message in session.
		http.Redirect(w, r, LoginURL, http.StatusFound)
		return
	}
	// If we have a url the Provider wants to make a redirect before
	// proceeding.
	if url != "" {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	// If we don't have a URL or an error then the user has been authenticated.
	// Check the Profile for an ID and Provider.
	if up.ID == "" || up.Provider == "" {
		panic(`hal/auth: The Profile's "ID" or "Provider" is empty.` +
			`A Key can not be created.`)
	}
	if _, err = CreateAndLogin(w, r, up); err != nil {
		// TODO: set error message in session.
		http.Redirect(w, r, LoginURL, http.StatusFound)
		return
	}
	// If we've made it this far redirect to the SuccessURL
	http.Redirect(w, r, SuccessURL, http.StatusFound)
	return
}
