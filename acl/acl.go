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
func genID(objKey *datastore.Key, groupId, perm string) string {
	return fmt.Sprintf("%s|%s|%s", objKey.String(), groupId, strings.ToLower(perm))
}

func genKey(c appengine.Context, groupId, perm string, objKey *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, "Perm", genID(objKey, groupId, perm), 0, nil)
}

func Auth(c appengine.Context, groupId, perm string, objKey *datastore.Key) error {
	key := genKey(c, groupId, perm, objKey)
	if _, err := put(c, key); err != nil {
		return err
	}
	return nil
}

func Can(c appengine.Context, groupId, perm string, objKey *datastore.Key) (bool, error) {
	key := genKey(c, groupId, perm, objKey)
	if _, err := get(c, key); err != nil {
		return false, err
	}
	return true, nil
}
