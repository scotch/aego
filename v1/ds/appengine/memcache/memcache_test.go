// Copyright 2012 The AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memcache

import (
	"appengine"
	"appengine/datastore"
	"github.com/scotch/aego/v1/context"
	dse "github.com/scotch/aego/v1/ds/errors"
	"testing"
)

var s = new(Store)

func setup() {
	_ = s.Clear()
}

func tearDown() {
	context.Close()
}

type X struct {
	Name string
	S    string
	Ss   []string
	I    int
	Ii   []int
	B    []byte
	K    *datastore.Key
	Kk   []*datastore.Key
}

func NewX(c appengine.Context, name string) (key *datastore.Key, x *X) {
	x = &X{
		Name: name,
		S:    "1",
		Ss:   []string{"2", "3"},
		I:    1,
		Ii:   []int{2, 3},
		B:    []byte{'1'},
		K:    datastore.NewKey(c, "X", "1", 0, nil),
		Kk: []*datastore.Key{
			datastore.NewKey(c, "X", "2", 0, nil),
			datastore.NewKey(c, "X", "3", 0, nil),
		},
	}
	key = datastore.NewKey(c, "X", name, 0, nil)
	return
}

func CheckX(c appengine.Context, t *testing.T, x *X, name string) {

	if x.Name != name {
		t.Errorf("x.Name: %v; want: %v.", x.Name, name)
	}
	if x.S != "1" {
		t.Errorf("x.S: %v; want: %v.", x.S, "1")
	}
	ss := []string{"2", "3"}
	if len(x.Ss) != len(ss) {
		t.Errorf("x.Ss: %v; want: %v.", x.Ss, ss)
	}
	if x.I != 1 {
		t.Errorf("x.I: %v; want: %v.", x.I, 1)
	}
	ii := []int{2, 3}
	if len(x.Ii) != len(ii) {
		t.Errorf("x.Ii: %v; want: %v.", x.Ii, ii)
	}
	b := []byte{'1'}
	if len(x.B) != len(b) {
		t.Errorf("x.B: %v; want: %v.", x.B, b)
	}
	k := datastore.NewKey(c, "X", "1", 0, nil)
	if x.K.String() != k.String() {
		t.Errorf("x.K: %v; want: %v.", x.K, k)
	}
	kk := []*datastore.Key{
		datastore.NewKey(c, "X", "2", 0, nil),
		datastore.NewKey(c, "X", "3", 0, nil),
	}
	if len(x.Kk) != len(kk) {
		t.Errorf("x.Kk: %v; want: %v.", x.Kk, kk)
	}
	return
}

func TestPut(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	k, x1 := NewX(c, "1")
	if cnt := s.Count(c); cnt != 0 {
		t.Errorf(`Before put; s.Count(c) = %v; want %v`, cnt, 0)
	}
	_, err := s.Put(c, k, x1)
	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}
	if cnt := s.Count(c); cnt != 1 {
		t.Errorf(`After put; s.Count(c): %v; want %v`, cnt, 1)
	}
}

func TestPutMulti(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	if cnt := s.Count(c); cnt != 0 {
		t.Errorf(`Before Put; s.Count(c): %v; want %v`, cnt, 0)
	}
	k1, x1 := NewX(c, "X1")
	k2, x2 := NewX(c, "X2")
	k3, x3 := NewX(c, "X3")
	keys := []*datastore.Key{k1, k2, k3}
	xs := []*X{x1, x2, x3}

	keys, err := s.PutMulti(c, keys, xs)
	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}
	if cnt := s.Count(c); cnt != 3 {
		t.Errorf(`After Put; s.Count(c): %v; want %v`, cnt, 3)
	}
}

func TestGet(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	// Put.

	k, x := NewX(c, "1")
	_, err := s.Put(c, k, x)

	// Get.

	r := new(X)
	err = s.Get(c, k, r)
	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}

	// Confirm.

	CheckX(c, t, r, "1")

	// Get non-existence

	k = datastore.NewKey(c, "X", "fakekey", 0, nil)
	r = new(X)
	err = s.Get(c, k, r)

	if err != dse.ErrNoSuchEntity {
		t.Errorf("err: %v; want: %v.", err, dse.ErrNoSuchEntity)
	}
}

func TestGetMulti(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	if cnt := s.Count(c); cnt != 0 {
		t.Errorf(`Before Put; s.Count(c): %v; want %v`, cnt, 0)
	}

	// Put.

	k1, x1 := NewX(c, "1")
	k2, x2 := NewX(c, "2")
	k3, x3 := NewX(c, "3")
	keys := []*datastore.Key{k1, k2, k3}
	xs := []*X{x1, x2, x3}
	keys, err := s.PutMulti(c, keys, xs)

	// Get.

	xs = []*X{&X{}, &X{}, &X{}}
	err = s.GetMulti(c, keys, xs)

	// Confirm.

	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}
	CheckX(c, t, xs[0], "1")
	CheckX(c, t, xs[1], "2")
	CheckX(c, t, xs[2], "3")

	// Get non-existence

	k4, _ := NewX(c, "4")
	k5, _ := NewX(c, "5")
	xs = []*X{&X{}, &X{}, &X{}, &X{}, &X{}}
	keys = []*datastore.Key{k1, k2, k3, k4, k5}

	err = s.GetMulti(c, keys, xs)

	if err.Error() != "ds: no such entity (and 1 other error)" {
		t.Errorf("err: %v; want: %v.", err.Error(), "ds: no such entity (and 1 other error)")
	}

	// Check.

	CheckX(c, t, xs[0], "1")
	CheckX(c, t, xs[1], "2")
	CheckX(c, t, xs[2], "3")
}

func TestDelete(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	// Put.

	k, x := NewX(c, "1")
	_, err := s.Put(c, k, x)

	// Delete.
	err = s.Delete(c, k)
	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}

	// Confirm.
	if cnt := s.Count(c); cnt != 0 {
		t.Errorf(`After delete; Count: %v; want %v`, cnt, 0)
	}

	// Delete non-existence

	k = datastore.NewKey(c, "X", "fakekey", 0, nil)
	r := new(X)
	err = s.Get(c, k, r)

	if err != dse.ErrNoSuchEntity {
		t.Errorf("err: %v; want: %v.", err, dse.ErrNoSuchEntity)
	}
}

func TestDeleteMulti(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	// Put.

	k1, x1 := NewX(c, "1")
	k2, x2 := NewX(c, "2")
	k3, x3 := NewX(c, "3")
	keys := []*datastore.Key{k1, k2, k3}
	xs := []*X{x1, x2, x3}
	keys, err := s.PutMulti(c, keys, xs)

	// Delete.

	err = s.DeleteMulti(c, keys)

	// Confirm.

	if err != nil {
		t.Errorf("err: %v; want: %v.", err, nil)
	}
	if cnt := len(xs); cnt != 3 {
		t.Errorf(`After DeleteMulti; len(xs): %v; want %v`, cnt, 3)
	}

	// Delete non-existence

	// Put.

	k1, x1 = NewX(c, "1")
	k2, x2 = NewX(c, "2")
	k3, x3 = NewX(c, "3")
	keys = []*datastore.Key{k1, k2, k3}
	xs = []*X{x1, x2, x3}
	keys, err = s.PutMulti(c, keys, xs)

	k4, _ := NewX(c, "4")
	k5, _ := NewX(c, "5")
	xs = []*X{&X{}, &X{}, &X{}, &X{}, &X{}}
	keys = []*datastore.Key{k1, k2, k3, k4, k5}

	err = s.DeleteMulti(c, keys)

	if err.Error() != "ds: no such entity (and 1 other error)" {
		t.Errorf("err: %v; want: %v.", err.Error(), "ds: no such entity (and 1 other error)")
	}
}
