// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mail

import (
	"appengine/mail"
	"testing"
)

var (
	sender  string   = "Company Name <noreply@example.com>"
	subject string   = "Example Subject"
	to      []string = []string{"1@example.com"}
)

func TestTemplateMessage(t *testing.T) {

	v := map[string]string{
		"one": "one",
	}
	bptxt := "templates/test/01.txt"
	bphtml := "templates/test/01.html"
	msg := &TemplateMessage{
		Message: &mail.Message{
			Sender:  sender,
			To:      to,
			Subject: subject,
		},
		BodyValues:       v,
		BodyTmplPath:     bptxt,
		HTMLBodyTmplPath: bphtml,
	}

	err := msg.Execute()
	if err != nil {
		t.Errorf(`err: %v, want "%v"`, err, nil)
	}
	if b, h := string(msg.Body), string("one == one\n"); b != h {
		t.Errorf(`msg.Body: %q, want %q`, b, h)
	}
	if b, h := string(msg.HTMLBody), string("<h1>one == one</h1>\n"); b != h {
		t.Errorf(`msg.Body: %q, want %q`, b, h)
	}
}
