// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
//"appengine"
//"appengine/datastore"
//"encoding/json"
//"fmt"
//"github.com/scotch/hal/ds"
//"time"
)

type Owner struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type Thumb struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Image struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	MimeType string   `json:"mimeType"`
	URL      string   `json:"url"`
	Kind     string   `json:"kind"`
	Size     int64    `json:"size"`
	Owners   []*Owner `json:"owners"`
	Thumbs   []*Thumb `json:"thumbs"`
}

// // Collection is used for JSON requests and responses
// // NOT USED YET
// type Collection struct {
// 	Kind  string     `json:"kind"`
// 	Items []*Product `json:"items"`
// }
// 
// // Entity is the struct for storing the data in the datastore.
// type Entity struct {
// 	Key     *datastore.Key `datastore:"-"`
// 	Data    *Product       `datastore:"-"`
// 	JSON    []byte
// 	Owners  []string
// 	Created time.Time
// 	Updated time.Time
// }
// 
// func NewEntity(c appengine.Context, p *Product, id int64) (e *Entity) {
// 	e = &Entity{
// 		Data:    p,
// 		Created: time.Now(),
// 		Updated: time.Now(),
// 	}
// 	e.SetKey(c, id)
// 	e.Encode()
// 	return
// }
// 
// func DeleteEntity(c appengine.Context, id int64) error {
// 	k := datastore.NewKey(c, "Product", "", id, nil)
// 	err := ds.Delete(c, k)
// 	return err
// }
// 
// func Get(c appengine.Context, id int64) (*Entity, error) {
// 	var e *Entity
// 	k := datastore.NewKey(c, "Product", "", id, nil)
// 	err := ds.Get(c, k, &e)
// 	e.Key = k
// 	e.Decode()
// 	return e, err
// }
// 
// func GetAll(c appengine.Context) ([]*Entity, error) {
// 	q := datastore.NewQuery("Product").Order("-Updated").Limit(50)
// 	var ents []*Entity
// 	_, err := q.GetAll(c, &ents)
// 	return ents, err
// }
// 
// func (e *Entity) SetKey(c appengine.Context, id int64) (err error) {
// 	if id == 0 {
// 		id, _, err = ds.AllocateIDs(c, "Product", nil, 1)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	e.Key = datastore.NewKey(c, "Product", "", id, nil)
// 	return
// }
// 
// func (e *Entity) Decode() error {
// 	if e.JSON != nil {
// 		var p *Product
// 		err := json.Unmarshal(e.JSON, &p)
// 		e.Data = p
// 		return err
// 	}
// 	return nil
// }
// 
// func (e *Entity) Encode() error {
// 	if e.Data != nil {
// 		e.Data.ID = fmt.Sprintf("%v", e.Key.IntID())
// 		j, err := json.Marshal(e.Data)
// 		e.JSON = j
// 		return err
// 	}
// 	return nil
// }
// 
// func (e *Entity) Put(c appengine.Context) (err error) {
// 	err = e.Encode()
// 	if err != nil {
// 		return
// 	}
// 	e.Updated = time.Now()
// 	key, err := ds.Put(c, e.Key, e)
// 	e.Key = key
// 	return
// }
