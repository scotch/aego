// Copyright 2012 HAL Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
hal/person provides structs representing various things. The objective of
this package is to standardize the discription of things that are commonly
used in web applications,
*/
package person

import ()

type Person struct {
	// AboutMe: A short biography for this person.
	AboutMe string `json:"aboutMe,omitempty"`

	// Birthday: The person's date of birth, represented as YYYY-MM-DD.
	Birthday string `json:"birthday,omitempty"`

	// CurrentLocation: The current location for this person.
	CurrentLocation string `json:"currentLocation,omitempty"`

	// Created: The time in Unix time when the User was created.
	Created int64 `json:"created,omitempty"`

	// Updated: The time in Unix time when the User was updated.
	Updated int64 `json:"updated,omitempty"`

	// DisplayName: The name of this person, suitable for display.
	DisplayName string `json:"displayName,omitempty"`

	// Email is the default email address
	Email string `json:"email,omitempty"`

	// Emails: A list of email addresses for this person.
	Emails []*PersonEmails `json:"emails,omitempty"`

	// Etag: ETag of this response for caching purposes.
	Etag string `json:"etag,omitempty"`

	// Gender: The person's gender. Possible values are: 
	// - "male" - Male gender. 
	// - "female" - Female gender. 
	// - "other" - Other.
	Gender string `json:"gender,omitempty"`

	// HasApp: If "true", indicates that the person has installed the app
	// that is making the request and has chosen to expose this install
	// state to the caller. A value of "false" indicates that the install
	// state cannot be determined (it is either not installed or the person
	// has chosen to keep this information private).
	HasApp bool `json:"hasApp,omitempty"`

	// Id: The ID of this person.
	ID string `json:"id,omitempty"`

	// Image: The representation of the person's profile photo.
	Image *PersonImage `json:"image,omitempty"`

	// Kind: Identifies this resource as a person. Value: "plus#person".
	Kind string `json:"kind,omitempty"`

	// LanguagesSpoken: The languages spoken by this person.
	LanguagesSpoken []string `json:"languagesSpoken,omitempty"`

	// Name: An object representation of the individual components of a
	// person's name.
	Name *PersonName `json:"name,omitempty"`

	// Nickname: The nickname of this person.
	Nickname string `json:"nickname,omitempty"`

	// ObjectType: Type of person within Google+. Possible values are:  
	// - "person" - represents an actual person. 
	// - "page" - represents a page.
	ObjectType string `json:"objectType,omitempty"`

	// Organizations: A list of current or past organizations with which
	// this person is associated.
	Organizations []*PersonOrganizations `json:"organizations,omitempty"`

	// PlacesLived: A list of places where this person has lived.
	PlacesLived []*PersonPlacesLived `json:"placesLived,omitempty"`

	// Password: A Password object used for password changes.
	Password *PersonPassword `json:"password,omitempty"`

	// Provider: The source of the Person Profile
	Provider *PersonProvider `json:"provider,omitempty"`

	// RelationshipStatus: The person's relationship status. Possible values are:  
	// - "single" - Person is single. 
	// - "in_a_relationship" - Person is in a relationship. 
	// - "engaged" - Person is engaged. 
	// - "married" - Person is married. 
	// - "its_complicated" - The relationship is complicated. 
	// - "open_relationship" - Person is in an open relationship. 
	// - "widowed" - Person is widowed. 
	// - "in_domestic_partnership" - Person is in a domestic partnership. 
	// - "in_civil_union" - Person is in a civil union.
	RelationshipStatus string `json:"relationshipStatus,omitempty"`

	// Roles: a list of strings indicating the roles the user has
	Roles []string `json:"roles,omitempty"`

	// Tagline: The brief description (tagline) of this person.
	Tagline string `json:"tagline,omitempty"`

	// Url: The URL of this person's profile.
	URL string `json:"url,omitempty"`

	// Urls: A list of URLs for this person.
	Urls []*PersonUrls `json:"urls,omitempty"`
}

func New() *Person {
	return &Person{}
}

func (p *Person) Validate() {}

type PersonEmails struct {
	// Primary: If "true", indicates this email address is the person's
	// primary one.
	Primary bool `json:"primary,omitempty"`

	// Type: The type of address. Possible values are: 
	// - "home" - Home email address. 
	// - "work" - Work email address. 
	// - "other" - Other.
	Type string `json:"type,omitempty"`

	// Value: The email address.
	Value string `json:"value,omitempty"`
}

type PersonImage struct {
	// URL: The URL of the person's profile photo. To re-size the image and
	// crop it to a square, append the query string ?sz=x, where x is the
	// dimension in pixels of each side.
	URL string `json:"url,omitempty"`
}

type PersonName struct {
	// FamilyName: The family name (last name) of this person.
	FamilyName string `json:"familyName,omitempty"`

	// Formatted: The full name of this person, including middle names,
	// suffixes, etc.
	Formatted string `json:"formatted,omitempty"`

	// GivenName: The given name (first name) of this person.
	GivenName string `json:"givenName,omitempty"`

	// HonorificPrefix: The honorific prefixes (such as "Dr." or "Mrs.") for
	// this person.
	HonorificPrefix string `json:"honorificPrefix,omitempty"`

	// HonorificSuffix: The honorific suffixes (such as "Jr.") for this
	// person.
	HonorificSuffix string `json:"honorificSuffix,omitempty"`

	// MiddleName: The middle name of this person.
	MiddleName string `json:"middleName,omitempty"`
}

type PersonOrganizations struct {
	// Department: The department within the organization.
	Department string `json:"department,omitempty"`

	// Description: A short description of the person's role in this
	// organization.
	Description string `json:"description,omitempty"`

	// EndDate: The date the person left this organization.
	EndDate string `json:"endDate,omitempty"`

	// Location: The location of this organization.
	Location string `json:"location,omitempty"`

	// Name: The name of the organization.
	Name string `json:"name,omitempty"`

	// Primary: If "true", indicates this organization is the person's
	// primary one (typically interpreted as current one).
	Primary bool `json:"primary,omitempty"`

	// StartDate: The date the person joined this organization.
	StartDate string `json:"startDate,omitempty"`

	// Title: The person's job title or role within the organization.
	Title string `json:"title,omitempty"`

	// Type: The type of organization. Possible values are: 
	// - "work" - Work. 
	// - "school" - School.
	Type string `json:"type,omitempty"`
}

type PersonPlacesLived struct {
	// Primary: If "true", this place of residence is this person's primary
	// residence.
	Primary bool `json:"primary,omitempty"`

	// Value: A place where this person has lived. For example: "Seattle,
	// WA", "Near Toronto".
	Value string `json:"value,omitempty"`
}

type PersonProvider struct {
	// Name: the name of the provider
	Name string `json:"name,omitempty"`

	// URL: The url of the provider
	URL string `json:"url,omitempty"`
}

type PersonUrls struct {
	// Primary: If "true", this URL is the person's primary URL.
	Primary bool `json:"primary,omitempty"`

	// Type: The type of URL. Possible values are: 
	// - "home" - URL for home. 
	// - "work" - URL for work. 
	// - "blog" - URL for blog. 
	// - "profile" - URL for profile. 
	// - "other" - Other.
	Type string `json:"type,omitempty"`

	// Value: The URL value.
	Value string `json:"value,omitempty"`
}

type PersonPassword struct {
	// New: the new password.
	New string `json:"new,omitempty"`

	// Current: The current password.
	Current string `json:"current,omitempty"`

	// IsSet: Indictor that the User has created a password.
	IsSet bool `json:"isSet"`
}
