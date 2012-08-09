// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/ds/appengine/datastore provides App Engine datastore persistence.

This package is a wrapper for the appengine/datastore package.
*/
package datastore

import (
	"appengine"
	"appengine/datastore"
	dserrors "github.com/scotch/hal/ds/errors"
)

type Store struct{}

// New creates a new datastore.Store
func New() *Store {
	return &Store{}
}

// Count returns the total number of items in the store.
// this method is for testing ONLY.
func (s *Store) Count(c appengine.Context) int64 {
	q := datastore.NewQuery("X")
	cnt, _ := q.Count(c)
	return int64(cnt)
}

// PutMulti given a []*datastore.Key and a struct pointer adds multiple entities
// to the store".
func (s *Store) PutMulti(c appengine.Context, key []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	// TODO(kylefinley) Error codes should be converted to hal/ds errors.
	return datastore.PutMulti(c, key, src)
}

// Put given a *datastore.Key and a struct pointer adds a single entity
// to the store.
func (s *Store) Put(c appengine.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	// TODO(kylefinley) Error codes should be converted to hal/ds errors.
	return datastore.Put(c, key, src)
}

// GetMulti given a []*datastore.Key returns multiple entities from the store
func (s *Store) GetMulti(c appengine.Context, key []*datastore.Key, dst interface{}) error {
	// TODO(kylefinley) Error codes should be converted to hal/ds errors.
	// This needs to be optimized.
	dserr := datastore.GetMulti(c, key, dst)
	if dserr != nil {
		err := make(appengine.MultiError, len(key))
		for i, e := range dserr.(appengine.MultiError) {
			if e != nil {
				if e == datastore.ErrNoSuchEntity {
					err[i] = dserrors.ErrNoSuchEntity
				} else {
					err[i] = e
				}
			}
		}
		return err
	}
	return nil
}

// Get given a *datastore.Key returns a single entity from the store
func (s *Store) Get(c appengine.Context, key *datastore.Key, dst interface{}) (err error) {
	err = datastore.Get(c, key, dst)
	if err == datastore.ErrNoSuchEntity {
		err = dserrors.ErrNoSuchEntity
	}
	return
}

// DeleteMulti given a []*datastore.Key deletes multiple entities from the store
func (s *Store) DeleteMulti(c appengine.Context, key []*datastore.Key) (err error) {
	// TODO(google) if the supplied key does not exist, Datastore should
	// return datastore.ErrNoSuchEntity instead of panicing.
	// TODO(kylefinley) figure out a way to catch the panic here.
	defer func() {
		if r := recover(); r != nil {
			err = dserrors.ErrNoSuchEntity
			return
		}
	}()
	err = datastore.DeleteMulti(c, key)
	return
}

// Delete given a *datastore.Key deletes a single entity from the store
func (s *Store) Delete(c appengine.Context, key *datastore.Key) (err error) {
	// TODO(google) if the supplied key does not exist, Datastore should
	// return datastore.ErrNoSuchEntity instead of panicing.
	// TODO(kylefinley) figure out a way to catch the panic here.
	defer func() {
		if r := recover(); r != nil {
			err = dserrors.ErrNoSuchEntity
			return
		}
	}()
	err = datastore.Delete(c, key)
	return
}
