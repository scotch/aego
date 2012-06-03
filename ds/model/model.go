// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
ds/model: EXPERIMENTAL

*/
package model

import (
	//"appengine"
	"appengine/datastore"
	//dserrors "github.com/scotch/hal/ds/errors"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/ds"
	//"reflect"
)

var models = make(map[string]interface{})

type Interface interface {
	//Put(c appengine.Context)
}

type Model struct {
	Key *datastore.Key `datastore:",-"`
}

//func New(c appengine.Context, kind string, stringID string, intID int64, parent *datastore.Key) *Model {
//	//t := reflect.TypeOf(m)
//	//v := reflect.ValueOf(m)
//	//kind := t.Name()
//	//elemType := v.Type().Elem()
//	//m = elemType.Kind()
//	//i := v.Interface()
//	key := datastore.NewKey(c, kind, stringID, intID, parent)
//	//n := new(i)
//	//n.key = key
//	//return n
//	return &Model{key: key}
//}
//
//func New(c appengine.Context, s interface{}, stringID string, intID int64, parent *datastore.Key) {
//	t := reflect.TypeOf(s)
//	kind := t.Name()
//	key := datastore.NewKey(c, kind, stringID, intID, parent)
//	newObjPtr := reflect.New(reflect.TypeOf(s).Elem())
//	newObj := reflect.Indirect(newObjPtr)
//	newObj.key = key
//	return newObjPtr.Interface()
//}

func Register(kind string, model interface{}) {
	models[kind] = model
	return
}

func (m *Model) SetKey(c context.Context, kind, stringID string, intID int64,
	parent *datastore.Key) {

	m.Key = datastore.NewKey(c, kind, stringID, intID, parent)
	return
}

func (m *Model) Put(c context.Context) error {
	k, err := ds.Put(c, m.Key, m)
	m.Key = k
	return err
}

// 
// func Get(c context.Context, k interface{}, stringID string, intID int64) (
// 	interface{}, error) {
// 
// 	typ := reflect.TypeOf(k)
// 	// if a pointer to a struct is passed, get the type of the dereferenced object
// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}
// 	knd := reflect.TypeOf(k).Name()
// 	key := datastore.NewKey(c, knd, stringID, intID, nil)
// 	err := ds.Get(c, key, k)
// 	//knd.Key = key
// 	//k.(Model).Key = key
// 	//t := reflect.TypeOf(kind).Kind()
// 	return k, err
// }
// 
func Get(c context.Context, stringID string, k interface{}) error {

	//typ := reflect.TypeOf(k)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	//if typ.Kind() == reflect.Ptr {
	//typ = typ.Elem()
	//}
	//knd := reflect.TypeOf(k).Name()
	knd := "Person"
	key := datastore.NewKey(c, knd, stringID, 0, nil)
	err := ds.Get(c, key, k)
	//knd.Key = key
	//k.(Model).Key = key
	//t := reflect.TypeOf(kind).Kind()
	return err
}

func (m *Model) GetByID(c context.Context, model interface{}, id string) {

}
