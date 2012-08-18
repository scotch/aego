// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"appengine/delay"
	"fmt"
	"github.com/scotch/hal/config"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/mail"
	"github.com/scotch/hal/user/token"
	"net/http"
)

const (
	URL_PREFIX = "/-/email/verify"
)

var SendConfirmAddressLater = delay.Func("confirmaddress", SendConfirmAddress)

var defaultConfig = map[string]string{
	"SiteTitle":        "Company Name",
	"SiteURL":          "http://localhost:8080",
	"SenderAddress":    "Company Name <noreply@example.com>",
	"Subject":          "Please verify your email '%s'",
	"BodyTmplPath":     "github.com/scotch/hal/user/email/templates/confirm.txt",
	"BodyHTMLTmplPath": "github.com/scotch/hal/user/email/templates/confirm.html",
}

func init() {
	http.HandleFunc(URL_PREFIX, verifyHandler)
}

func genVerifyURL(baseURL string, token string) string {
	return baseURL + URL_PREFIX + "?token=" + token
}

func SendConfirmAddress(c context.Context, e *Email) (err error) {
	v, err := config.GetOrInsert(c, "app", defaultConfig)
	if err != nil {
		return err
	}
	et := token.New(c)
	et.EmailAddress = e.Address
	if err = et.Put(c); err != nil {
		return err
	}
	v.Values["VerfyURL"] = genVerifyURL(v.Values["SiteURL"], et.Token)
	m := &mail.TemplateMessage{}
	m.Sender = v.Values["SenderAddress"]
	m.To = []string{e.Address}
	m.Subject = fmt.Sprintf(v.Values["Subject"], e.Address)
	m.BodyValues = v.Values
	m.BodyTmplPath = v.Values["BodyTmplPath"]
	mail.SendLater.Call(c, m)
	return nil
}

func (e *Email) SendConfirmAddress(c context.Context) (err error) {
	SendConfirmAddressLater.Call(c, e)
	return nil
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	//code := r.FormValue("code")
	var e *Email
	var et *token.Token
	var err error
	c := context.NewContext(r)
	code := r.URL.Query().Get("code")
	errURL := "/"
	successURL := "/"
	et, err = token.Get(c, code)
	if err != nil {
		goto Error
	}
	e, err = Get(c, et.EmailAddress)
	if err != nil {
		goto Error
	}
	e.Status = verified
	if err = e.Put(c); err != nil {
		goto Error
	}
	http.Redirect(w, r, successURL, http.StatusFound)

Error:
	// TODO added error to session
	http.Redirect(w, r, errURL, http.StatusNotFound)
}
