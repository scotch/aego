// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright: 2011 Google Inc. All Rights Reserved.
// license: Apache Software License, see LICENSE for details.

package user_profile

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/scotch/hal/ds"
	"github.com/scotch/hal/types"
	"strings"
	"time"
)

type UserProfile struct {
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
	// Auth maybe used by the provodier to store any information that it
	// may need.
	Auth []byte
	// PersonJSON is the JSON encoded representation of a aego/types.Person
	PersonJSON []byte
	// PersonRawJSON is the JSON encoded representation of the raw
	// response returned from a provider representing the User's Profile.
	PersonRawJSON []byte
	// Created is a time.Time representing with the UserProfile was created.
	Created time.Time
	// Created is a time.Time representing with the UserProfile was updated.
	Updated time.Time
}

// New creates a new UserProfile and set the Created to now
func New() *UserProfile {
	return &UserProfile{
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// GenAuthID generates a unique id for the UserProfile based on a string
// representing the provider and a unique id provided by the provider.
// Using this generator is prefered over compiling the key manually to
// ensure consistency.
func GenAuthID(provider, id string) string {
	return fmt.Sprintf("%v|%v", strings.ToLower(provider), id)
}

// NewKey generates a *datastore.Key based on a string representing
// the provider and a unique id provided by the provider.
func NewKey(c appengine.Context, provider, id string) *datastore.Key {
	authID := GenAuthID(provider, id)
	return datastore.NewKey(c, "UserProfile", authID, 0, nil)
}

// Get is a convience method for retrieveing an entity foom the store.
func Get(c appengine.Context, id string, up *UserProfile) (err error) {
	key := datastore.NewKey(c, "UserProfile", id, 0, nil)
	err = ds.Get(c, key, up)
	up.Key = key
	return
}

// Put is a convience method to save the UserProfile to the datastore and
// updated the Updated property to time.Now().
func (u *UserProfile) Put(c appengine.Context) error {
	u.Updated = time.Now()
	// TODO add error handeling for empty Provider and ID
	u.Key = NewKey(c, u.Provider, u.ID)
	key, err := ds.Put(c, u.Key, u)
	u.Key = key
	return err
}

func (u *UserProfile) SetPerson(p *types.Person) error {
	b, err := json.Marshal(p)
	u.PersonJSON = b
	return err
}

func (u *UserProfile) Person() (*types.Person, error) {
	p := new(types.Person)
	err := json.Unmarshal(u.PersonJSON, p)
	return p, err
}

// check aborts the current execution if err is non-nil.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

//// GetOrInsert creates a UserProfile using the email address as the key.
//// If a UserProfile already exists it updates it instead.
//func GetOrInsert(c appengine.Context, provider, id string,
//	auth, person, personRaw []byte) (*datastore.Key, error) {
//
//	key := datastore.NewKey(c, "UserProfile", id, 0, nil)
//	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
//		e := new(UserProfile)
//		err := datastore.Get(c, key, e)
//		if err != nil && err != datastore.ErrNoSuchEntity {
//			c.Errorf("CreateOrUpdate err: %v", err)
//			check(err)
//			return nil
//		}
//		// If a profile dosen't exist; create it.
//		if err == datastore.ErrNoSuchEntity {
//			e.Created = time.Now()
//		}
//		e.Provider = provider
//		e.Updated = time.Now()
//		e.Auth = auth
//		e.Person = person
//		e.PersonRaw = personRaw
//		if _, err := datastore.Put(c, key, e); err != nil {
//			check(err)
//		}
//		return nil
//	}, nil)
//	return key, err
//}