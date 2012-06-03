// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appengine_openid

import (
	aeuser "appengine/user"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/types"
	"github.com/scotch/hal/user_profile"
	"net/http"
)

type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{
		"AppEngineOpenID",
		"http://appengine.google.com",
	}
}

// Authenticate process the request and returns a populated UserProfile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request,
	up *user_profile.UserProfile) (url string, err error) {

	c := context.NewContext(r)

	// Set provider info.

	up.Provider = p.Name
	up.ProviderURL = p.URL

	// Check for current User.

	u := aeuser.Current(c)

	if u == nil {

		url := r.FormValue("provider")
		redirectURL := r.URL.Path + "/callback"
		loginUrl, err := aeuser.LoginURLFederated(c, redirectURL, url)
		return loginUrl, err

	}

	if u.FederatedIdentity != "" {
		up.ID = u.FederatedIdentity
	} else {
		up.ID = u.ID
	}

	per := new(types.Person)
	per.Emails = []*types.PersonEmails{
		&types.PersonEmails{true, "home", u.Email},
	}
	per.Url = u.FederatedIdentity

	up.SetPerson(per)

	return "", nil
}