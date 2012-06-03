// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ds

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	dserrors "github.com/scotch/hal/ds/errors"
	"reflect"
	"testing"
	"time"
)

type A struct {
	S string
	I int
	B []byte
	T time.Time
}

type B struct {
	S string
	I int
	B []byte
	T time.Time
}

type C struct {
	S string
	I int
	B []byte
	T time.Time
}

var (
	now = time.Now()
)

func setup() {
	Register("A", true, true, true)
	Register("B", true, true, false)
	Register("C", true, false, false)
	Register("D", false, false, true)
	_ = memory.Clear()
	_ = memcache.Clear()
}

func tearDown() {
	context.Close()
}

func TestRegister(t *testing.T) {
	defer tearDown()
	Register("A", true, true, true)
	_, ok := kinds["A"]
	if ok == false {
		t.Errorf(`ok == %s; expected true`, ok)
	}
	Register("B", true, true, false)
	_, ok = kinds["B"]
	if ok == false {
		t.Errorf(`ok == %s; expected true`, ok)
	}
}

func TestPutGet(t *testing.T) {
	setup()
	c := context.NewContext(nil)
	defer tearDown()
	var x interface{}
	// Put A
	x = &A{S: "a", I: 1, B: []byte{}, T: now}
	key := datastore.NewKey(c, "A", "1", 0, nil)
	key, err := Put(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Test Store for values
	x = new(A)
	err = memory.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	err = memcache.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	err = dsds.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Put B
	x = &B{S: "a", I: 1, B: []byte{}, T: now}
	key = datastore.NewKey(c, "B", "1", 0, nil)
	key, err = Put(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Test Store for values
	x = new(B)
	err = memory.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
	err = memcache.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	err = dsds.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Put C
	x = &C{S: "a", I: 1, B: []byte{}, T: now}
	key = datastore.NewKey(c, "C", "1", 0, nil)
	key, err = Put(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	x = new(C)
	// Test Store for values
	err = memory.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
	err = memcache.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
	err = dsds.Get(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
}

func TestPutGetMulti(t *testing.T) {
	setup()
	c := context.NewContext(nil)
	defer tearDown()
	// Put As
	key := []*datastore.Key{
		datastore.NewKey(c, "A", "1", 0, nil),
		datastore.NewKey(c, "A", "2", 0, nil),
	}
	x := []*A{
		&A{S: "a", I: 1, B: []byte{}, T: now},
		&A{S: "a", I: 1, B: []byte{}, T: now},
	}
	key, err := PutMulti(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Test Store for values
	x = []*A{&A{}, &A{}}
	err = dsds.GetMulti(c, key, x)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
}

func TestDelete(t *testing.T) {
	setup()
	c := context.NewContext(nil)
	defer tearDown()
	var x interface{}
	// Put A
	x = &A{S: "a", I: 1, B: []byte{}, T: now}
	key := datastore.NewKey(c, "A", "1", 0, nil)
	key, _ = Put(c, key, x)
	// Delete A
	err := Delete(c, key)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Test Store for absence of values
	x = &A{}
	err = memory.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
	err = memcache.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
	err = dsds.Get(c, key, x)
	if err != dserrors.ErrNoSuchEntity {
		t.Errorf(`err = %s; expected %s`, err, dserrors.ErrNoSuchEntity)
	}
}

func TestDeleteMulti(t *testing.T) {
	setup()
	c := context.NewContext(nil)
	defer tearDown()
	var x interface{}
	// Put As
	x = []*A{
		&A{S: "a", I: 1, B: []byte{}, T: now},
		&A{S: "a", I: 1, B: []byte{}, T: now},
	}
	key := []*datastore.Key{
		datastore.NewKey(c, "A", "1", 0, nil),
		datastore.NewKey(c, "A", "2", 0, nil),
	}
	key, _ = PutMulti(c, key, x)
	// Delete As
	err := DeleteMulti(c, key)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	// Test Store for absence of values
	x = []*A{&A{}, &A{}}
	err = memory.GetMulti(c, key, x)
	if reflect.TypeOf(err) != reflect.TypeOf(dserrors.MultiError{}) {
		t.Errorf(`err = %s; expected %s`, reflect.TypeOf(err), reflect.TypeOf(dserrors.MultiError{}))
	}
	err = memcache.GetMulti(c, key, x)
	if reflect.TypeOf(err) != reflect.TypeOf(dserrors.MultiError{}) {
		t.Errorf(`err = %s; expected %s`, reflect.TypeOf(err), reflect.TypeOf(dserrors.MultiError{}))
	}
	// TODO(kylefinley) add this test back.
	// err = dsds.GetMulti(c, key, x)
	// if reflect.TypeOf(err) != reflect.TypeOf(dserrors.MultiError{}) {
	// 	t.Errorf(`err = %s; expected %s`, reflect.TypeOf(err), reflect.TypeOf(dserrors.MultiError{}))
	// }
}

func TestAllocateIDs(t *testing.T) {
	setup()
	c := context.NewContext(nil)
	defer tearDown()

	cnt := 5
	low, high, err := AllocateIDs(c, "D", nil, cnt)
	if err != nil {
		t.Errorf(`err = %s; expected nil`, err)
	}
	ncnt := int(high - low)
	if ncnt != cnt {
		t.Errorf(`ncnt = %v, %v`, ncnt, cnt)
		t.Errorf(`low = %v`, low)
		t.Errorf(`high = %v`, high)
	}
}
