// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acl

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"github.com/scotch/hal/ds"
	"strings"
)

type Perm struct{}

func get(c appengine.Context, key *datastore.Key) (p *Perm, err error) {
	p = &Perm{}
	err = ds.Get(c, key, p)
	return
}

func put(c appengine.Context, key *datastore.Key) (p *Perm, err error) {
	p = &Perm{}
	_, err = ds.Put(c, key, p)
	return
}
func genKey(c appengine.Context, group, perm string, objKey *datastore.Key) (
	key *datastore.Key) {

	id := fmt.Sprintf("%s|%s", strings.ToLower(perm), group)
	key = datastore.NewKey(c, "Perm", id, 0, objKey)
	return
}

func Auth(c appengine.Context, group, perm string, objKey *datastore.Key) error {
	key := genKey(c, group, perm, objKey)
	if _, err := put(c, key); err != nil {
		return err
	}
	return nil
}

func Can(c appengine.Context, group, perm string, objKey *datastore.Key) (bool, error) {
	key := genKey(c, group, perm, objKey)
	if _, err := get(c, key); err != nil {
		return false, err
	}
	return true, nil
}
