// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/scotch/hal/ds"
	"strings"
	"time"
)

var ErrInvalidAddress = errors.New("email: invalid address")

// Validate returns an ErrInvalidEmail if the supplied
// string does not contains an "@" and a ".".
func Validate(address string) error {

	// TODO maybe use a regex here instead?
	if len(address) < 5 {
		return ErrInvalidAddress
	}
	if !strings.Contains(address, "@") {
		return ErrInvalidAddress
	}
	if !strings.Contains(address, ".") {
		return ErrInvalidAddress
	}
	return nil
}

const (
	unconfirmed = iota
	pending
	confirmed
	suspended
)

type Email struct {
	Key     *datastore.Key `datastore:"-"`
	Address string
	UserId  int64
	Status  int64
	// Created is a time.Time of when the Email was first created.
	Created time.Time
	// Updated is a time.Time of when the Email was last updated.
	Updated time.Time
}

// New creates a new email entity and set Created to now
func New() *Email {
	return &Email{
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func (e *Email) Put(c appengine.Context) (err error) {
	e.Updated = time.Now()
	e.Key, err = ds.Put(c, e.Key, e)
	return
}

func Get(c appengine.Context, address string) (*Email, error) {
	e := New()
	key := datastore.NewKey(c, "Email", strings.ToLower(address), 0, nil)
	err := ds.Get(c, key, e)
	e.Key = key
	return e, err
}

// func CreateOrUpdate(c appengine.Context, userkey *datastore.Key,
// 	email string) (*datastore.Key, error) {
// 
// 	key := datastore.NewKey(c, "Email", email, 0, nil)
// 	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
// 		e := new(Email)
// 		err := datastore.Get(c, key, e)
// 		// If an email dosen't exist, create it.
// 		if err == datastore.ErrNoSuchEntity {
// 			e.UserKey = userkey
// 			if _, err := datastore.Put(c, key, e); err != nil {
// 				return err
// 			}
// 			return nil
// 		}
// 		if err != nil {
// 			c.Errorf("Email CreateOrUpdate err: %v", err)
// 			return err
// 		}
// 		return nil
// 	}, nil)
// 	return key, err
// }
