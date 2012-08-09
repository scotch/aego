// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
memory: Stores the entity in instance memory, i.e. RAM.

Entites are stored in a map using the datastore.Key.Encode() as the key
with the stuct encoded using encoding/gob.

TODO(kylefinley) Add a method for expiring entites. Old entities should
also be removed when spaces is limited.

*/
package memory

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/gob"
	"errors"
	dserrors "github.com/scotch/hal/ds/errors"
	"reflect"
	"sync"
)

type Store struct {
	cache map[string]interface{}
	mu    sync.Mutex
}

// New creates a new memory.Store
func New() *Store {
	return &Store{
		cache: make(map[string]interface{}),
	}
}

// Clear removes all entites from the store.
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]interface{})
	return nil
}

// Count returns the total number of items in the store.
func (s *Store) Count(c appengine.Context) int64 {
	return int64(len(s.cache))
}

// PutMulti given a []*datastore.Key and a struct pointer adds multiple entities
// to the store.
func (s *Store) PutMulti(c appengine.Context, key []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	v := reflect.ValueOf(src)
	if len(key) != v.Len() {
		return nil, errors.New(
			"ds/memory: key and src slices have different length")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		var d bytes.Buffer
		elem := v.Index(i)
		enc := gob.NewEncoder(&d)
		err := enc.Encode(elem.Interface())
		s.cache[key[i].Encode()] = d
		if err != nil {
			multiErr[i] = err
			any = true
		}
	}
	if any {
		return key, multiErr
	}
	return key, nil
}

// Put given a *datastore.Key and a struct pointer adds a single entity
// to the store.
func (s *Store) Put(c appengine.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	k, err := s.PutMulti(c, []*datastore.Key{key}, []interface{}{src})
	if err != nil {
		return nil, err
	}
	return k[0], nil
}

// GetMulti given a []*datastore.Key returns multiple entities from the store
func (s *Store) GetMulti(c appengine.Context, key []*datastore.Key, dst interface{}) (err error) {
	v := reflect.ValueOf(dst)
	if len(key) != v.Len() {
		return errors.New("ds/memory: key and dst slices have different length")
	}
	if len(key) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		data := s.cache[key[i].Encode()]
		if data == nil {
			multiErr[i] = dserrors.ErrNoSuchEntity
			any = true
			continue
		}
		b := data.(bytes.Buffer)
		dec := gob.NewDecoder(&b)
		elem := v.Index(i)
		err = dec.Decode(elem.Interface())
		if err != nil {
			multiErr[i] = err
			any = true
		}
	}
	if any {
		return multiErr
	}
	return nil
}

// Get given a *datastore.Key returns a single entity from the store
func (s *Store) Get(c appengine.Context, key *datastore.Key, dst interface{}) error {
	err := s.GetMulti(c, []*datastore.Key{key}, []interface{}{dst})
	if me, ok := err.(dserrors.MultiError); ok {
		return me[0]
	}
	return err
}

// DeleteMulti given a []*datastore.Key deletes multiple entities from the store
func (s *Store) DeleteMulti(c appengine.Context, key []*datastore.Key) (err error) {
	if len(key) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		present := s.cache[key[i].Encode()]
		if present == nil {
			multiErr[i] = dserrors.ErrNoSuchEntity
			any = true
			continue
		}
		delete(s.cache, key[i].Encode())
	}
	if any {
		return multiErr
	}
	return err
}

// Delete given a *datastore.Key deletes a single entity from the store
func (s *Store) Delete(c appengine.Context, key *datastore.Key) error {
	err := s.DeleteMulti(c, []*datastore.Key{key})
	if me, ok := err.(dserrors.MultiError); ok {
		return me[0]
	}
	return err
}
