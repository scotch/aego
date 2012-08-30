// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/scotch/aego/v1/context"
	"testing"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
}

func tearDown() {
	context.Close()
}

func TestGetOrInsert(t *testing.T) {
	defer context.Close()
	c := context.NewContext(nil)

	// Set it.

	m := map[string]string{
		"A": "1",
	}
	cnfg, err := GetOrInsert(c, "first", m)

	// Confirm.

	checkErr(t, err)
	if x := cnfg.Values["A"]; x != "1" {
		t.Errorf(`config["A"]: %v, want %v`, x, "1")
	}

	// The orginal map should be returned.

	m = map[string]string{
		"A": "2",
	}

	cnfg, err = GetOrInsert(c, "first", m)

	// Confirm.

	checkErr(t, err)
	if x := cnfg.Values["A"]; x != "1" {
		t.Errorf(`config["A"]: %v, want %v`, x, "1")
	}
}

func TestEdit(t *testing.T) {
	defer context.Close()
	c := context.NewContext(nil)

	// Set it.

	m := map[string]string{
		"A": "1",
	}
	cnfg, err := GetOrInsert(c, "first", m)

	// Change it.

	cnfg.Values["A"] = "2"

	// Save it.

	err = cnfg.Put(c)
	checkErr(t, err)

	// Confirm.

	cnfg, err = Get(c, "first")

	checkErr(t, err)
	if x := cnfg.Values["A"]; x != "2" {
		t.Errorf(`cnfg["A"]: %v, want %v`, x, "2")
	}
}
