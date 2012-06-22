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
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/user"
	"github.com/scotch/hal/user_profile"
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
	Authenticate(http.ResponseWriter, *http.Request,
		*user_profile.UserProfile) (string, error)
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

// createAndLogin saves the UserProfile to the datastore. And appends
// the Key.StringID() to the current User's AuthIDs. If a User has not
// yet been created, it creates one.
func createAndLogin(w http.ResponseWriter, r *http.Request,
	up *user_profile.UserProfile) (u *user.User, err error) {

	c := context.NewContext(r)
	// Save it.
	err = up.Put(c)
	if err != nil {
		return
	}
	u, err = user.Current(r)
	if err != nil {
		// If the User isn't logged in. Create an User and log them in.
		u, err = user.GetOrInsertByAuthID(c, up.Key.StringID())
		_ = user.CurrentUserSetID(w, r, u.Key.IntID())
		return
	}
	u.AddAuthID(up.Key.StringID())
	err = u.Put(c)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {

	up := user_profile.New()
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

	// Check the UserProfile for an ID and Provider.
	if up.ID == "" || up.Provider == "" {
		panic(`hal/auth: The UserProfile's "ID" or "Provider" is empty.` +
			`A Key can not be created.`)
	}

	if _, err = createAndLogin(w, r, up); err != nil {
		// TODO: set error message in session.
		http.Redirect(w, r, LoginURL, http.StatusFound)
		return
	}
	// If we've made it this far redirect to the SuccessURL
	http.Redirect(w, r, SuccessURL, http.StatusFound)
	return

}
