// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright: 2011 Google Inc. All Rights Reserved.
// license: Apache Software License, see LICENSE for details.

package profile

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/person"
	"strings"
	"time"
)

type Profile struct {
	// Key is the datastore key. It is not saved back to the datastore
	// it is just embeded here for convience.
	Key *datastore.Key `datastore:"-"`
	// ID represents a unique ID from the Provider.
	// This ID does not have to be unique to this application just to the
	// provider.
	ID string
	// A String representing the Provider that performed the
	// authentication. The Provider should be in the proper case for
	// example a User who was authenticated through Google should have
	// "Google" here and not "google"
	Provider string
	// ProviderURL is the URL that is commonly accepted as the
	// originator of the authentication. For example Google plus would
	// be http://plus.google.com and not http://google.com.
	ProviderURL string `datastore:",noindex"`
	// UserID is the string ID of the User that the Profile belongs to.
	UserID string
	// Auth maybe used by the provodier to store any information that it
	// may need.
	Auth []byte
	// Person is an Object representing personal information about the user.
	Person *person.Person `datastore:"-"`
	// PersonJSON is the Person object converted to JSON, for storage purposes.
	PersonJSON []byte `datastore:"Person"`
	// PersonRawJSON is the JSON encoded representation of the raw
	// response returned from a provider representing the User's Profile.
	PersonRawJSON []byte
	// Created is a time.Time representing with the Profile was created.
	Created time.Time
	// Created is a time.Time representing with the Profile was updated.
	Updated time.Time
}

// New creates a new Profile and set the Created to now
func New() *Profile {
	return &Profile{
		Person:  new(person.Person),
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// GenAuthID generates a unique id for the Profile based on a string
// representing the provider and a unique id provided by the provider.
// Using this generator is prefered over compiling the key manually to
// ensure consistency.
func GenAuthID(provider, id string) string {
	return fmt.Sprintf("%s|%s", strings.ToLower(provider), id)
}

// NewKey generates a *datastore.Key based on a string representing
// the provider and a unique id provided by the provider.
func NewKey(c appengine.Context, provider, id string) *datastore.Key {
	authID := GenAuthID(provider, id)
	return datastore.NewKey(c, "AuthProfile", authID, 0, nil)
}

// SetKey creates and embeds a ds.Key to the entity.
func (u *Profile) SetKey(c appengine.Context) (err error) {
	u.Key = NewKey(c, u.Provider, u.ID)
	return
}

// Encode is called prior to save. Any fields that need to be updated
// prior to save are updated here.
func (u *Profile) Encode() error {
	// Update Person

	// Sanity check, TODO maybe we should raise an error instead.
	if u.Person == nil {
		u.Person = new(person.Person)
	}
	u.Person.Provider = &person.PersonProvider{
		Name: u.Provider,
		URL:  u.ProviderURL,
	}
	u.Person.Kind = fmt.Sprintf("%s#person", strings.ToLower(u.Provider))
	u.Person.ID = u.ID
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
func (u *Profile) Decode() error {
	if u.PersonJSON != nil {
		var p *person.Person
		err := json.Unmarshal(u.PersonJSON, &p)
		u.Person = p
		return err
	}
	return nil
}

// Get is a convience method for retrieveing an entity from the ds.
func Get(c appengine.Context, id string, up *Profile) (err error) {
	key := datastore.NewKey(c, "AuthProfile", id, 0, nil)
	err = ds.Get(c, key, up)
	up.Key = key
	up.Decode()
	return
}

// Put is a convience method to save the Profile to the datastore and
// updated the Updated property to time.Now().
func (u *Profile) Put(c appengine.Context) error {
	// TODO add error handeling for empty Provider and ID
	u.SetKey(c)
	u.Updated = time.Now()
	u.Encode()
	key, err := ds.Put(c, u.Key, u)
	u.Key = key
	return err
}
