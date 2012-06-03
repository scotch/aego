// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
ds/model: EXPERIMENTAL

*/
package model

import (
	"appengine/datastore"
	"github.com/scotch/hal/context"
	"testing"
)

type Person struct {
	Model
	FirstName  string
	FamilyName string
}

func (p *Person) FullName() string {
	return p.FirstName + " " + p.FamilyName
}

func setup() {
	//Register("Person", Person, true, true, false)
	defer tearDown()
}

func tearDown() {
	context.Close()
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`err => %v, want nil`, err)
	}
}

// func TestNew(t *testing.T) {
// 	setup()
// 	defer tearDown()
// 	c := context.NewContext(nil)
// 	p := &Person{}
// 	//x := New(c, "Person", "1", 0, nil)
// 	x := New(c, p, "1", 0, nil)
// 	k := datastore.NewKey(c, "Person", "1", 0, nil)
// 	if p.key.String() != k.String() {
// 		t.Errorf(`x.Key => %v, want %v`, x.key, k)
// 	}
// 	if p.Kind() != "Person" {
// 		t.Errorf(`x.Kind() => %v, want %v`, p.Kind(), "Person")
// 	}
// }

func TestSetKey(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	x := &Person{FirstName: "Kyle", FamilyName: "Finley"}
	x.SetKey(c, "Person", "1", 0, nil)
	k := datastore.NewKey(c, "Person", "1", 0, nil)
	if x.Key.String() != k.String() {
		t.Errorf(`x.Key => %v, want %v`, x.Key, k)
	}
	//if x.Kind() != "Person" {
	//t.Errorf(`x.Kind() => %v, want %v`, x.Kind(), "Person")
	//}
}

func TestPut(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)

	x := &Person{FirstName: "Kyle", FamilyName: "Finley"}
	x.SetKey(c, "Person", "1", 0, nil)
	err := x.Put(c)
	if err != nil {
		t.Errorf(`err => %v, want nil`, err)
	}
	k := datastore.NewKey(c, "Person", "1", 0, nil)
	if x.Key.String() != k.String() {
		t.Errorf(`x.Key => %v, want %v`, x.Key, k)
	}
	if x.FirstName != "Kyle" {
		t.Errorf(`x.FirstName = %v, want "Kyle"`, x.FirstName)
	}
}

func TestGet(t *testing.T) {
	setup()
	defer tearDown()
	c := context.NewContext(nil)
	// Put it.
	//Register("Person", &Person{})
	x1 := &Person{FirstName: "Kyle"}
	x1.SetKey(c, "Person", "1", 0, nil)
	_ = x1.Put(c)
	// Get it.
	// TODO: I don't like this api. I don't want to have to perform
	// type assertion here. Until there's a better way use the hal/ds api.
	// var p Person
	// x2, _ := Get(c, &p, "1", 0)
	// if x1.FirstName == x2.(Person).FirstName {
	// 	t.Errorf(`x2.FirstName => %v, want %v`, x2.(Person).FirstName, x1.FirstName)
	// }
	// x3 := x2.(Person)
	// if x1.FirstName == x3.FirstName {
	// 	t.Errorf(`x2.FirstName => %v, want %v`, x3.FirstName, x1.FirstName)
	// }
	// var x3 Person
	// err := Get(c, "1", &x3)
	// checkErr(t, err)
	// if x1.FirstName == x3.FirstName {
	// 	t.Errorf(`x2.FirstName => %v, want %v`, x3.FirstName, x1.FirstName)
	// }

}

// func TestGetOrInsert(t *testing.T) {
// 	setup()
// 	defer tearDown()
// 	c := context.NewContext(nil)
// 	var p Person
// 	x := p.GetOrInsert(c, "1", nil)
// 	x := &Person{}.GetOrInsert(c, "1", nil)
// 	k := datastore.NewKey(c, "1", 0, nil)
// 	if x.Key != k {
// 		t.Errorf(`x.Key => %v, want %v`, x.Key, k)
// 	}
// }
// func TestDelete(t *testing.T) {
// 	setup()
// 	defer tearDown()
// 	c := context.NewContext(nil)
// 	var p Person
// 	x1 := p.New(c, "1", 0, nil)
// 	x1.FirstName = "Kyle"
// 	x1.FamilyName = "Finley"
// 	_ := x1.Put(c)
// 	x2, err := p.Get(c, "1", nil)
// 	checkErr(t, err)
// 	x2.Delete(c)
// 	x3, err := p.Get(c, "1", nil)
// 	if err != dserrors.ErrNoSuchEntity {
// 		t.Errorf(`err => %v, want %v`, err, dserrors.ErrNoSuchEntity)
// 	}
// }
// 
// func TestAllocateIDs(t *testing.T) {
// 	setup()
// 	defer tearDown()
// 	c := context.NewContext(nil)
// 	keys := Person{}.AllocateIDs(c)
// }
// 
// func TestQuery(t *testing.T) {
// 	setup()
// 	c := context.NewContext(nil)
// 	results := Person{}.Query()
// }
