// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/session"
	"net/http"
	"time"
)

var (
	ErrNoLoggedInUser = errors.New("hal/user: no logged in user")
)

// User definition
type User struct {
	// AuthIDs is a list of string represting additional authentications
	// stategies. E.g.
	//
	//    ["google|12345", "facebook|12345"]
	//
	AuthIDs []string
	// Email is the primary email address. Used for notifications.
	Email string
	// Emails is a list of additional email addresses. Used in quering.
	Emails []string
	// Password is a password hast used to verify the user.
	Password []byte
	// Roles is a list of role names that the user belongs to.
	Roles []string
	// Created is a time.Time of when the User was first created.
	Created time.Time
	// Updated is a time.Time of when the User was last updated.
	Updated time.Time
}

// New creates a new user and set the Created to now
func New() *User {
	return &User{
		//Email:   email,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// Get is a convience method for retrieveing an entity foom the store.
func Get(c appengine.Context, id int64) (*User, *datastore.Key, error) {
	u := &User{}
	key := datastore.NewKey(c, "User", "", id, nil)
	err := ds.Get(c, key, u)
	return u, key, err
}

// Put is a convience method to save the User to the datastore and
// updated the Updated property to time.Now().
func (u *User) Put(c appengine.Context, key *datastore.Key) (*datastore.Key, error) {
	u.Updated = time.Now()
	key, err := ds.Put(c, key, u)
	return key, err
}

// CurrentUserID returns the userID of the requesting user.
func CurrentUserID(r *http.Request) (int64, error) {
	s, err := session.Store.Get(r, "auth")
	if err != nil {
		return 0, err
	}
	userID, _ := s.Values["userid"].(int64)
	return userID, err
}

// SetCurrentUserID adds the provided userID to the current users session/cookie
func SetCurrentUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
	s, err := session.Store.Get(r, "auth")
	if err != nil {
		return err
	}
	s.Values["userid"] = userID
	s.Save(r, w)
	return nil
}

// Logout sets the session userid to 0, effectivly logging the user out.
func Logout(w http.ResponseWriter, r *http.Request) error {
	return SetCurrentUserID(w, r, 0)
}

// Current checks the requesting User's session to see if they have an
// account. If they do, the provided User struct is populated with the
// information that is saved in the datastore. If they don't an error is
// returned.
// NOTE: this method will likely change in un upcomming release. The
// appengine.Context argument will be removed.
func Current(c appengine.Context, r *http.Request, u *User) (*datastore.Key, error) {
	// TODO(kylefinley) We need either the appengine.Context or the
	// *http.Request not both. There's a testing issue that prevents
	// converting the request to a context mid function, once this is
	// resulved remove the appengine.Context argument.
	userID, _ := CurrentUserID(r)
	if userID != 0 {
		key := datastore.NewKey(c, "User", "", userID, nil)
		err := ds.Get(c, key, u)
		return key, err
	}
	return nil, ErrNoLoggedInUser
}

// IsAdmin returns true if the requesting user is an admin; otherwise
// returns false.
func IsAdmin(r *http.Request) (bool, error) {
	s, err := session.Store.Get(r, "auth")
	if err != nil {
		return false, err
	}
	isAdmin, _ := s.Values["isadmin"].(int64)
	if isAdmin == 1 {
		return true, nil
	}
	return false, nil
}
