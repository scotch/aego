// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
memcache: Stores the entity in App Engine memcache.

Entites are stored using the datastore.Key.Encode() as the key
with the stuct encoded using encoding/gob.

*/

package memcache

import (
	"appengine"
	"appengine/datastore"
	aemc "appengine/memcache"
	"bytes"
	"encoding/gob"
	"errors"
	dserrors "github.com/scotch/hal/ds/errors"
	"github.com/scotch/hal/ds/utils"
	"reflect"
)

type Store struct{}

// New creates a new memcache.Store
func New() *Store {
	return &Store{}
}

// Clear removes all entites from the store.
func (s *Store) Clear() error {
	// TODO(kylefinley) This method should clear memcache.
	return nil
}

// Count returns the total number of items in the store.
func (s *Store) Count(c appengine.Context) int64 {
	stats, _ := aemc.Stats(c)
	return int64(stats.Items)
}

// PutMulti given a []*datastore.Key and a struct pointer adds multiple entities
// to the store.
func (s *Store) PutMulti(c appengine.Context, key []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	v := reflect.ValueOf(src)
	if len(key) != v.Len() {
		return nil, errors.New(
			"ds/memory: key and src slices have different length")
	}
	// TODO(kylefinley) we should use PutMulti here.
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		// TODO(kylefinley) memcache has a 1mb size limit.
		// We should make sure the entity doesn't exceed that amount.
		var d bytes.Buffer
		elem := v.Index(i)
		enc := gob.NewEncoder(&d)
		err := enc.Encode(elem.Interface())
		if err != nil {
			multiErr[i] = err
			any = true
			continue
		}
		id := key[i].Encode()
		item := &aemc.Item{
			Key:   id,
			Value: d.Bytes(),
		}
		if err := aemc.Set(c, item); err != nil {
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

	multiArgType, _ := utils.CheckMultiArg(v)
	if multiArgType == utils.MultiArgTypeInvalid {
		return errors.New("ds/memory: dst has invalid type")
	}
	if len(key) != v.Len() {
		return errors.New("ds/memory: key and dst slices have different length")
	}
	if len(key) == 0 {
		return nil
	}
	// if err := utils.ValidateKeys(key); err != nil {
	// 	return err
	// }

	// TODO(kylefinley) we should use GetMulti here.
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		data, err := aemc.Get(c, key[i].Encode())
		if err != nil && err != aemc.ErrCacheMiss {
			multiErr[i] = dserrors.ErrNoSuchEntity
			any = true
			continue
		}
		// TODO(kylefinley) maybe this can be removed. Can memcache be trusted?
		if data == nil || data.Value == nil {
			multiErr[i] = dserrors.ErrNoSuchEntity
			any = true
			continue
		}
		var b bytes.Buffer
		b.Write(data.Value)
		dec := gob.NewDecoder(&b)
		elem := v.Index(i)
		if multiArgType == utils.MultiArgTypeStruct {
			elem = elem.Addr()
		}
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
	// TODO(kylefinley) we should use DeleteMulti here.
	multiErr, any := make(dserrors.MultiError, len(key)), false
	for i := range key {
		err := aemc.Delete(c, key[i].Encode())
		if err != nil {
			multiErr[i] = dserrors.ErrNoSuchEntity
			any = true
		}
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
