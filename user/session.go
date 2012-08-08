// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/email"
	"github.com/scotch/hal/session"
	"net/http"
)

var (
	ErrNoLoggedInUser = errors.New("user: no logged in user")
)

// CurrentUserID returns the userId of the requesting user.
func CurrentUserID(r *http.Request) (string, error) {
	s, err := session.Store.Get(r, "user")
	if err != nil {
		return "", err
	}
	id, _ := s.Values["userid"].(string)
	return id, err
}

// CurrentUserIDByEmail returns the userId of the requesting user. Or the userID
// associated with the provided email.
func CurrentUserIDByEmail(r *http.Request, emailAddress string) (string, error) {
	// TODO: User merge if the session UserID is different then the email UserID
	// search session
	sessID, _ := CurrentUserID(r)
	if sessID != "" {
		// TODO: maybe confirm that the UserID exists?
		// There are case where the session may have an incorrect UserID.
		return sessID, nil
	}
	// search by email
	c := context.NewContext(r)
	e, err := email.Get(c, emailAddress)
	if err != nil {
		return "", err
	}
	return e.UserID, nil
}

// CurrentUserSetID adds the provided userId to the current users session/cookie
func CurrentUserSetID(w http.ResponseWriter, r *http.Request, userId string) error {
	s, err := session.Store.Get(r, "user")
	if err != nil {
		c := appengine.NewContext(r)
		c.Criticalf("user: There was an error retrieving the session Error: %v", err)
		return err
	}
	s.Values["userid"] = userId

	return s.Save(r, w)
}

// Current checks the requesting User's session to see if they have an
// account. If they do, the provided User struct is populated with the
// information that is saved in the datastore. If they don't an error is
// returned.
func Current(r *http.Request) (*User, error) {
	id, _ := CurrentUserID(r)

	if id != "" {
		c := context.NewContext(r)
		u := new(User)
		key := datastore.NewKey(c, "User", id, 0, nil)
		err := ds.Get(c, key, u)
		u.Key = key
		return u, err
	}
	return nil, ErrNoLoggedInUser
}

// Logout sets the session userid to "", effectivly logging the user out.
// TODO maybe delete cookie, instead.
func Logout(w http.ResponseWriter, r *http.Request) error {
	return CurrentUserSetID(w, r, "")
}
