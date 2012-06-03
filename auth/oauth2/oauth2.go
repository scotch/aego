// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"appengine/urlfetch"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/scotch/hal/context"
	"github.com/scotch/hal/user_profile"
	"net/http"
	"net/url"
	"strings"
)

type Provider struct {
	Name         string
	URL          string
	ClientID     string
	ClientSecret string
	Scope        string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
}

func New(name, url, clientID, clientSecret, scope, authURL, tokenURL string) *Provider {
	return &Provider{
		Name:         name,
		URL:          url,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
	}
}

// Config returns the configuration information for OAuth2.
func (p *Provider) Config(url *url.URL) *oauth.Config {
	return &oauth.Config{
		ClientId:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Scope:        p.Scope,
		AuthURL:      p.AuthURL,
		TokenURL:     p.TokenURL,
		RedirectURL: fmt.Sprintf("%s://%s/-/auth/%s/callback", url.Scheme, url.Host,
			strings.ToLower(p.Name)),
	}
}

func (p *Provider) start(r *http.Request) string {
	return p.Config(r.URL).AuthCodeURL(r.URL.RawQuery)
}

func (p *Provider) callback(r *http.Request) error {
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	t := &oauth.Transport{
		Config: p.Config(r.URL),
		Transport: &urlfetch.Transport{
			Context: context.NewContext(r),
		},
	}
	_, err := t.Exchange(code)
	return err
}

func (p *Provider) Authenticate(r *http.Request) (
	redirectURL string, up *user_profile.UserProfile, err error) {

	//c := context.NewContext(r)
	return "", up, nil
}
