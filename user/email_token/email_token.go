// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email_token

import (
	"appengine"
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	"time"
)

type EmailToken struct {
	Key          *datastore.Key `datastore:"-"`
	Token        string         `datastore:"-"`
	UserID       string         `datastore:",noindex"`
	EmailAddress string         `datastore:",noindex"`
	Created      time.Time
}

func New(c context.Context) (e *EmailToken) {
	t := genToken()
	k := datastore.NewKey(c, "EmailToken", t, 0, nil)
	e = &EmailToken{
		Key:     k,
		Token:   t,
		Created: time.Now(),
	}
	return
}

func Get(c context.Context, token string) (e *EmailToken, err error) {
	e = &EmailToken{}
	k := datastore.NewKey(c, "EmailToken", token, 0, nil)
	if err = ds.Get(c, k, e); err != nil {
		return nil, err
	}
	e.Key = k
	e.Token = k.StringID()
	return
}

func (e *EmailToken) Put(c appengine.Context) (err error) {
	if e.Key == nil {
		panic("email_token: Key not set.")
	}
	key, err := ds.Put(c, e.Key, e)
	e.Key = key
	return
}

func (e *EmailToken) Delete(c context.Context) error {
	return ds.Delete(c, e.Key)
}
