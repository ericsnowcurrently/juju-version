// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package version implements version parsing.
package version

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/juju/utils/series"
	"gopkg.in/mgo.v2/bson"
)

var binaryPat = regexp.MustCompile(`^(\d{1,9})\.(\d{1,9})(\.|-(\w+))(\d{1,9})(\.\d{1,9})?-([^-]+)-([^-]+)$`)

// Binary specifies a binary version of juju.v
type Binary struct {
	Number
	Series string
	Arch   string
}

// MustParseBinary parses a binary version and panics if it does
// not parse correctly.
func MustParseBinary(s string) Binary {
	b, err := ParseBinary(s)
	if err != nil {
		panic(err)
	}
	return b
}

// ParseBinary parses a binary version of the form "1.2.3-series-arch".
func ParseBinary(s string) (Binary, error) {
	m := binaryPat.FindStringSubmatch(s)
	if m == nil {
		return Binary{}, fmt.Errorf("invalid binary version %q", s)
	}
	var b Binary
	b.Major = atoi(m[1])
	b.Minor = atoi(m[2])
	b.Tag = m[4]
	b.Patch = atoi(m[5])
	if m[6] != "" {
		b.Build = atoi(m[6][1:])
	}
	b.Series = m[7]
	b.Arch = m[8]
	_, err := series.GetOSFromSeries(b.Series)
	return b, err
}

// String returns the string representation of the binary version.
func (b Binary) String() string {
	return fmt.Sprintf("%v-%s-%s", b.Number, b.Series, b.Arch)
}

// GetBSON implements bson.Getter.
func (b Binary) GetBSON() (interface{}, error) {
	return b.String(), nil
}

// SetBSON implements bson.Setter.
func (b *Binary) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return err
	}
	v, err := ParseBinary(s)
	if err != nil {
		return err
	}
	*b = v
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Binary) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Binary) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := ParseBinary(s)
	if err != nil {
		return err
	}
	*b = v
	return nil
}

// MarshalYAML implements yaml.v2.Marshaller interface.
func (b Binary) MarshalYAML() (interface{}, error) {
	return b.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface.
func (b *Binary) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var vstr string
	err := unmarshal(&vstr)
	if err != nil {
		return err
	}
	v, err := ParseBinary(vstr)
	if err != nil {
		return err
	}
	*b = v
	return nil
}
