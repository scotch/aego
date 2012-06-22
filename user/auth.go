// Copyright 2012 The HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This file contains User methods connected to UserProfiles.

*/

package user

import (
	"appengine"
	"appengine/datastore"
	dserror "github.com/scotch/hal/ds/errors"
	"time"
)

// GetByAuthID gets a User by an associated UserProfile StringID.
func GetByAuthID(c appengine.Context, authID string) (u *User, err error) {
	q := datastore.NewQuery("User").
		Filter("AuthIDs =", authID).
		Limit(1)
	for t := q.Run(c); ; {
		u := new(User)
		key, err := t.Next(u)
		if err != nil {
			err = dserror.ErrNoSuchEntity
			return nil, err
		} else {
			u.Key = key
			return u, nil
		}
	}
	return
}

// CreateByAuthID creates a new user from an AuthID
func CreateByAuthID(c appengine.Context, authID string) (u *User, err error) {
	u = New()
	u.Created = time.Now()
	u.AuthIDs = []string{authID}
	u.Key = datastore.NewKey(c, "User", "", 0, nil)
	err = u.Put(c)
	return u, err
}

// GetOrInsertByAuthID creates or updates a User from a UserProfile Key
func GetOrInsertByAuthID(c appengine.Context, authID string) (u *User, err error) {
	u, err = GetByAuthID(c, authID)
	if err == nil {
		u.AddAuthID(authID)
		err = u.Put(c)
		return
	}
	u, err = CreateByAuthID(c, authID)
	return
}

// AddAuthID Adds an AuthID to the User's AuthIDs list
func (u *User) AddAuthID(authID string) {
	for _, id := range u.AuthIDs {
		if id == authID {
			return
		}
	}
	u.AuthIDs = append(u.AuthIDs, authID)
	return
}

//// UserProfiles returns the []UserProfiles owned by the User.
//func (u *User) UserProfiles(c appengine.Context) ([]user_profile.UserProfile, error) {
//	keys := make([]datastore.Key, len(u.AuthIDs))
//	ups := make([]user_profile.UserProfile, len(u.AuthIDs))
//	for _, id := range u.AuthIDs {
//		k := datastore.NewKey(c, "UserProfile", id, 0, nil)
//		append(keys, k)
//	}
//	err := ds.GetMulti(c, keys, ups)
//	return ups, err
//}
//
