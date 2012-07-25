// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/types"
	"time"
)

// User definition
type User struct {
	// The datastore Key
	Key *datastore.Key `datastore:"-"`
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
	// Password is a password hash used to verify the user.
	Password []byte
	// Roles is a list of role names that the user belongs to.
	Roles []string
	// Created is a time.Time of when the User was first created.
	Created time.Time
	// Updated is a time.Time of when the User was last updated.
	Updated time.Time
	// Person is an Object representing personal information about the user.
	Person *types.Person `datastore:"-"`
	// PersonJSON is the Person object converted to JSON, for storage purposes.
	PersonJSON []byte `datastore:"Person"`
}

// New creates a new user and set the Created to now
func New() *User {
	return &User{
		Person:  new(types.Person),
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func (u *User) Decode() error {
	if u.PersonJSON != nil {
		var p *types.Person
		err := json.Unmarshal(u.PersonJSON, &p)
		u.Person = p
		return err
	}
	return nil
}

func (u *User) Encode() error {
	// Update Person

	// Sanity check, maybe we should raise an error instead.
	if u.Person == nil {
		u.Person = new(types.Person)
	}
	u.Person.ID = u.Key.StringID()
	u.Person.Roles = u.Roles
	// TODO(kylefinley) consider alternatives to returning miliseconds.
	// Convert time to unix miliseconds for javascript
	u.Person.Created = u.Created.UnixNano() / 1000000
	u.Person.Updated = u.Updated.UnixNano() / 1000000
	// We don't want to return the password hash. So, we simply return a bool indicating that
	// the user has set there password.
	if len(u.Password) != 0 {
		u.Person.Password = &types.PersonPassword{IsSet: true}
	} else {
		u.Person.Password = &types.PersonPassword{IsSet: false}
	}
	// Convert to JSON
	j, err := json.Marshal(u.Person)
	u.PersonJSON = j
	return err
}

// Put is a convience method to save the User to the datastore and
// updated the Updated property to time.Now(). This method should
// always be usdd when saving a user, fore it does some necessary
// preprocessing.
func (u *User) Put(c appengine.Context) error {

	// If we are saving for the first time lets get an id so that we
	// can save the id to the json data before saving the entity. This
	// prevents us from having to save twice.
	if u.Key == nil || u.Key.StringID() == "" {
		id, err := ds.AllocateID(c, "User")
		if err != nil {
			return err
		}
		u.Key = datastore.NewKey(c, "User", id, 0, nil)
	}
	u.Updated = time.Now()
	u.Encode()
	key, err := ds.Put(c, u.Key, u)
	u.Key = key
	return err
}

// Get is a convience method for retrieveing an entity foom the store.
func Get(c appengine.Context, id string) (u *User, err error) {
	u = &User{}
	key := datastore.NewKey(c, "User", id, 0, nil)
	err = ds.Get(c, key, u)
	u.Key = key
	u.Decode()
	return u, err
}
