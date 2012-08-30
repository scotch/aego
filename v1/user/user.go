// Copyright 2012 The AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/scotch/aego/v1/ds"
	"github.com/scotch/aego/v1/person"
	"time"
)

const (
	active = iota
	suspended
	deleted
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
	// Roles is a list of role names that the user belongs to.
	Roles []string
	// Status is the state of the account values include:
	//  0 - Active
	//  1 - Suspended
	//  2 - Deleted
	Status int8
	// Person is an Object representing personal information about the user.
	Person *person.Person `datastore:"-"`
	// PersonJSON is the Person object converted to JSON, for storage purposes.
	PersonJSON []byte `datastore:"Person"`
	// Created is a time.Time of when the User was first created.
	Created time.Time
	// Updated is a time.Time of when the User was last updated.
	Updated time.Time
}

// New creates a new user and set the Created to now
func New() *User {
	return &User{
		Person:  new(person.Person),
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// AllocateID is a convenience method for allocating a UserID
func AllocateID(c appengine.Context) (id string, err error) {
	id, err = ds.AllocateID(c, "User")
	return
}

// SetKey creates and embeds a ds.Key into the entity.
func (u *User) SetKey(c appengine.Context) (err error) {
	// If we are saving for the first time lets get an id so that we
	// can save the id to the json data before saving the entity. This
	// prevents us from having to save twice.
	if u.Key == nil || u.Key.StringID() == "" {
		id, err := AllocateID(c)
		if err != nil {
			return err
		}
		u.Key = datastore.NewKey(c, "User", id, 0, nil)
	}
	return
}

// Encode is called prior to save. Any fields that need to be updated
// prior to save are updated here.
func (u *User) Encode() error {
	// Update Person

	// Sanity check, TODO maybe we should raise an error instead.
	if u.Person == nil {
		u.Person = new(person.Person)
	}
	u.Person.ID = u.Key.StringID()
	u.Person.Roles = u.Roles
	// TODO(kylefinley) consider alternatives to returning miliseconds.
	// Convert time to unix miliseconds for javascript
	u.Person.Created = u.Created.UnixNano() / 1000000
	u.Person.Updated = u.Updated.UnixNano() / 1000000
	// Convert to JSON
	j, err := json.Marshal(u.Person)
	u.PersonJSON = j
	return err
}

// Decode is called after the entity has been retrieved from the the ds.
func (u *User) Decode() error {
	if u.PersonJSON != nil {
		var p *person.Person
		err := json.Unmarshal(u.PersonJSON, &p)
		u.Person = p
		return err
	}
	return nil
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

// Put is a convience method to save the User to the datastore and
// updated the Updated property to time.Now(). This method should
// always be usdd when saving a user, fore it does some necessary
// preprocessing.
func (u *User) Put(c appengine.Context) (err error) {
	if err = u.SetKey(c); err != nil {
		return
	}
	u.Updated = time.Now()
	u.Encode()
	key, err := ds.Put(c, u.Key, u)
	u.Key = key
	return err
}
