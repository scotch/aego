// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hal/datastore provides cached presistence for the Google App Engine datastore.

TODO(kylefinley) Document this.


*/

// TODO(kylefinley) Ideally this package would allow for a more modular approach.
// The API would look more like this:
//
// Register("Entity", &memcache.Store{}, &datastore.Store{})
//
// Stores would be listed in the order in which they should be check. This approach
// would allow for custom stores to be easily created as well. Creating A MongoDB
// Store would be as simple as creating a Store that implements Storer, and passing
// a pointer to the Register() method.
//
// The end goal is to create a consistant interface for all types of data persistence.

// TODO(kylefinley) add support for gob encoding of invalid datastore types.
// E.g.
// type T struct {
//   A map[string]interface{} `ds:"gob"`
// }
// Would encode T.A to gob
// Or maybe this should be handled by ds/datastore.

package ds

import (
	"appengine"
	"appengine/datastore"
	dsdatastore "github.com/scotch/hal/ds/appengine/datastore"
	dsmemcache "github.com/scotch/hal/ds/appengine/memcache"
	dsmemory "github.com/scotch/hal/ds/memory"
	"time"
)

func init() {
	// Set the default store to memcache > datastore.
	Register("default", true, true, false)
}

var (
	kinds    = make(map[string]*StoreConfig)
	Kinds    = kinds
	dsds     = dsdatastore.New()
	memcache = dsmemcache.New()
	memory   = dsmemory.New()
)

// Storer is not implemented yet.
type Storer interface {
	Get(appengine.Context, *datastore.Key, interface{}) error
	GetMulti(appengine.Context, []*datastore.Key, interface{}) error
	Put(appengine.Context, *datastore.Key, interface{}) (*datastore.Key, error)
	PutMulti(appengine.Context, []*datastore.Key, interface{}) ([]*datastore.Key, error)
	Delete(appengine.Context, *datastore.Key) error
	DeleteMulti(appengine.Context, []*datastore.Key) error
}

// StoreConfig is used to establish the stores that should be used for an Entity.
// A true value indicated that that store should be used. A false that that store
// should not be used.
type StoreConfig struct {
	Datastore, Memcache, Memory bool
}

// Register takes a string representing the entity kind, followed by bools for
// datastore, memcache and memory with a value of true indicating that the entity
// should be stored in that store.
//
// Register("Product", true, true, false)
//
// Would set the StoreConfig for products. Based on the config "Product" entities
// will be stored to the datastore, and memcache, but not memory.
func Register(kind string, useDatastore, useMemcache, useMemory bool) {
	kinds[kind] = &StoreConfig{useDatastore, useMemcache, useMemory}
	return
}

// getConfig returns a StoreConfig based on entity kind. If the entity kind
// has not been saved the default StoreConfig is returned.
func getConfig(knd string) *StoreConfig {
	s, ok := kinds[knd]
	if !ok {
		s = kinds["default"]
	}
	return s
}

// PutMulti given a []*datastore.Key and a list of struct pointers adds
// multiple entities to the store.
func PutMulti(c appengine.Context, key []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	// TODO(kylefinley) After the datastore put this method should also put
	// to the various caches that the entity belongs.
	key, err := dsds.PutMulti(c, key, src)
	// Clear caches
	_ = memcache.DeleteMulti(c, key)
	_ = memory.DeleteMulti(c, key)
	return key, err
}

// Put given a *datastore.Key and a struct pointer adds a single entity
// to the store.
func Put(c appengine.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	sc := getConfig(key.Kind())

	key, err := dsds.Put(c, key, src)
	if err != nil {
		return key, err
	}
	if sc.Memcache {
		key, err = memcache.Put(c, key, src)
	}
	if sc.Memory {
		key, err = memory.Put(c, key, src)
	}
	return key, err
}

// GetMulti given a []*datastore.Key returns multiple entities from the store
func GetMulti(c appengine.Context, key []*datastore.Key, dst interface{}) (err error) {
	// TODO(kylefinley) It would save time and quota if this method only pulled
	// from the datastore values that where not cached.
	// The challenge arises when loading keys from different entity groups.
	// Each entity group may have different storage configurations, which
	// complicates things.
	// For the time being pull entity directly from the datastore.
	return dsds.GetMulti(c, key, dst)
}

// Get given a *datastore.Key returns a single entity from the store
func Get(c appengine.Context, key *datastore.Key, dst interface{}) (err error) {
	sc := getConfig(key.Kind())
	needMemory := true
	needMemcache := true
	if sc.Memory {
		err = memory.Get(c, key, dst)
		if err == nil {
			needMemory = false
			goto Complete
		}
	}
	if sc.Memcache {
		err = memcache.Get(c, key, dst)
		if err == nil {
			needMemcache = true
			goto Complete
		}
	}
	if sc.Datastore {
		err = dsds.Get(c, key, dst)
		if err == nil {
			goto Complete
		}
	}
	return
Complete:
	// Save the entity to any caches that it should be saved to.
	if sc.Memcache && needMemcache {
		_, _ = memcache.Put(c, key, dst)
	}
	if sc.Memory && needMemory {
		_, _ = memory.Put(c, key, dst)
	}
	return
}

// DeleteMulti given a []*datastore.Key deletes multiple entities from the store
func DeleteMulti(c appengine.Context, key []*datastore.Key) (err error) {
	if len(key) == 0 {
		return nil
	}
	err = dsds.DeleteMulti(c, key)
	_ = memcache.DeleteMulti(c, key)
	_ = memory.DeleteMulti(c, key)
	return
}

// Delete given a *datastore.Key deletes a single entity from the store
func Delete(c appengine.Context, key *datastore.Key) (err error) {
	err = dsds.Delete(c, key)
	_ = memcache.Delete(c, key)
	_ = memory.Delete(c, key)
	return
}

func AllocateIDs(c appengine.Context, kind string, parent *datastore.Key, n int) (low, high int64, err error) {
	// TODO: added for testing. Allocating IDs for mememache and memory
	// should not be used in production.
	sc := getConfig(kind)
	if sc.Datastore {
		return datastore.AllocateIDs(c, kind, parent, n)
	}
	t := time.Now()
	//var l int64
	l := t.UnixNano()
	return l, l + int64(n), nil
}
