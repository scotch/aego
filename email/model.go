// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"appengine"
	"appengine/datastore"
)

// User definition
type Email struct {
	UserKey *datastore.Key
}

func CreateOrUpdate(c appengine.Context, userkey *datastore.Key,
	email string) (*datastore.Key, error) {

	key := datastore.NewKey(c, "Email", email, 0, nil)
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		e := new(Email)
		err := datastore.Get(c, key, e)
		// If an email dosen't exist, create it.
		if err == datastore.ErrNoSuchEntity {
			e.UserKey = userkey
			if _, err := datastore.Put(c, key, e); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			c.Errorf("Email CreateOrUpdate err: %v", err)
			return err
		}
		return nil
	}, nil)
	return key, err
}
