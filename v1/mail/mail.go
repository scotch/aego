// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mail provides methods for sending e-mails
*/
package mail

import (
	"appengine/delay"
	"appengine/mail"
	"bytes"
	"github.com/scotch/aego/v1/context"
	"html/template"
)

var SendLater = delay.Func("send", Send)

type TemplateMessage struct {
	*mail.Message
	BodyValues       map[string]string
	BodyTmplPath     string
	HTMLBodyTmplPath string
}

func parseTmpl(tmplPath string, values map[string]string) (s string, err error) {
	var b bytes.Buffer
	tmpl := template.Must(template.ParseFiles(tmplPath))
	if err := tmpl.Execute(&b, values); err != nil {
		return s, err
	}
	return string(b.Bytes()), nil
}

func (m *TemplateMessage) Execute() (err error) {
	if m.BodyTmplPath != "" {
		if m.Body, err = parseTmpl(m.BodyTmplPath, m.BodyValues); err != nil {
			return
		}
	}
	if m.HTMLBodyTmplPath != "" {
		if m.HTMLBody, _ = parseTmpl(m.HTMLBodyTmplPath, m.BodyValues); err != nil {
			return
		}
	}
	return nil
}

func Send(c context.Context, m *TemplateMessage) (err error) {
	if err = m.Execute(); err != nil {
		return
	}
	if err = mail.Send(c, m.Message); err != nil {
		c.Errorf("email: sending confirmation email to %v failed: %v", m.To, err)
		return
	}
	return nil
}
