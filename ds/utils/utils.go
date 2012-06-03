// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"reflect"
)

type multiArgType int

const (
	MultiArgTypeInvalid multiArgType = iota
	//MultiArgTypePropertyLoadSaver
	MultiArgTypeStruct
	MultiArgTypeStructPtr
	MultiArgTypeInterface
)

// multiValid is a batch version of Key.valid. It returns an error, not a
// []bool.
// NOT USED YET
// func ValidateKeys(key []*datastore.Key) error {
// 	// From App Engine datastore.
// 	invalid := false
// 	for _, k := range key {
// 		if !k.valid() {
// 			invalid = true
// 			break
// 		}
// 	}
// 	if !invalid {
// 		return nil
// 	}
// 	err := make(dserrors.MultiError, len(key))
// 	for i, k := range key {
// 		if !k.valid() {
// 			err[i] = dserrors.ErrInvalidKey
// 		}
// 	}
// 	return err
// }

// CheckMultiArg checks that v has type []S, []*S, []I, or []P, for some struct
// type S, for some interface type I, or some non-interface non-pointer type P
// such that P or *P implements PropertyLoadSaver.
//
// It returns what category the slice's elements are, and the reflect.Type
// that represents S, I or P.
//
// As a special case, PropertyList is an invalid type for v.
func CheckMultiArg(v reflect.Value) (m multiArgType, elemType reflect.Type) {
	// From App Engine datastore.

	if v.Kind() != reflect.Slice {
		return MultiArgTypeInvalid, nil
	}
	//if v.Type() == typeOfPropertyList {
	//return MultiArgTypeInvalid, nil
	//}
	elemType = v.Type().Elem()
	//if reflect.PtrTo(elemType).Implements(typeOfPropertyLoadSaver) {
	//return MultiArgTypePropertyLoadSaver, elemType
	//}
	switch elemType.Kind() {
	case reflect.Struct:
		return MultiArgTypeStruct, elemType
	case reflect.Interface:
		return MultiArgTypeInterface, elemType
	case reflect.Ptr:
		elemType = elemType.Elem()
		if elemType.Kind() == reflect.Struct {
			return MultiArgTypeStructPtr, elemType
		}
	}
	return MultiArgTypeInvalid, nil
}

func uuid() string {
	// From Russ Cox
	// https://groups.google.com/forum/?fromgroups#!msg/golang-nuts/owCogizIuZs/ZzmwkQGrlnEJ
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
