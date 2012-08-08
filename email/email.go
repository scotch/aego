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

var (
	ErrInvalidAddress      = errors.New("email: invalid address")
	ErrAddressInUse        = errors.New("email: address in use by another user")
	ErrAddressAlreadyAdded = errors.New("email: address has already been added")
)

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
	UserID  string
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

func (e *Email) SetKey(c appengine.Context, address string) {
	e.Key = datastore.NewKey(c, "Email", address, 0, nil)
	return
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

func AddForUser(c appengine.Context, address, userID string, status int64) (e *Email, err error) {

	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		e, err = Get(c, address)
		if e.UserID != "" {
			if e.UserID != userID {
				return ErrAddressInUse
			}
			return ErrAddressAlreadyAdded
		}
		e = New()
		e.SetKey(c, address)
		e.UserID = userID
		e.Status = status
		return e.Put(c)
	}, nil)

	return e, err
}
