// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 Erik Unger. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package hal/email provides methods surronding:

- Email storage
- User lookup by email address
- Email addess validation
- Email address verification
*/
package email

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"github.com/scotch/hal/ds"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidAddress      = errors.New("email: invalid address")
	ErrAddressInUse        = errors.New("email: address in use by another user")
	ErrAddressAlreadyAdded = errors.New("email: address has already been added")
)

var emailRegexp *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9\-+~_%]+[a-zA-Z0-9\-+~_%.]*@([a-z]+[a-z0-9\\-]*\.)+[a-z][a-z]+`)

// Validate returns an ErrInvalidEmail if the supplied
// string does not contains an "@" and a ".".
func Validate(address string) error {
	if !emailRegexp.Match([]byte(strings.TrimSpace(address))) {
		return ErrInvalidAddress
	}
	return nil
}

const (
	unverified = iota
	pending
	verified
	primay
)

type Email struct {
	Key     *datastore.Key `json:"-",datastore:"-"`
	Address string         `json:"address"`
	UserID  string         `json:"-"`
	Status  int64          `json:"status"`
	// Created is a time.Time of when the Email was first created.
	Created time.Time `json:"created"`
	// Updated is a time.Time of when the Email was last updated.
	Updated time.Time `json:"updated"`
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
	e.Address = strings.ToLower(e.Key.StringID())
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

func GetMulti(c appengine.Context, ids []string) (ee []*Email, err error) {
	key := make([]*datastore.Key, len(ids))
	for k, id := range ids {
		key[k] = datastore.NewKey(c, "Email", id, 0, nil)
	}
	ee = make([]*Email, len(ids))
	for i := range ee {
		ee[i] = new(Email)
	}
	err = ds.GetMulti(c, key, ee)
	for i := range ee {
		ee[i].Key = key[i]
	}
	return
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
