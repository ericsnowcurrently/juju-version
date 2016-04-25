// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package version_test

import (
	"strings"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/version"
)

type BinarySuite struct{}

var _ = gc.Suite(&BinarySuite{})

func binaryVersion(major, minor, patch, build int, tag, series, arch string) version.Binary {
	return version.Binary{
		Number: version.Number{
			Major: major,
			Minor: minor,
			Patch: patch,
			Build: build,
			Tag:   tag,
		},
		Series: series,
		Arch:   arch,
	}
}

func (*BinarySuite) TestParseBinary(c *gc.C) {
	parseBinaryTests := []struct {
		v      string
		err    string
		expect version.Binary
	}{{
		v:      "1.2.3-trusty-amd64",
		expect: binaryVersion(1, 2, 3, 0, "", "trusty", "amd64"),
	}, {
		v:      "1.2.3.4-trusty-amd64",
		expect: binaryVersion(1, 2, 3, 4, "", "trusty", "amd64"),
	}, {
		v:      "1.2-alpha3-trusty-amd64",
		expect: binaryVersion(1, 2, 3, 0, "alpha", "trusty", "amd64"),
	}, {
		v:      "1.2-alpha3.4-trusty-amd64",
		expect: binaryVersion(1, 2, 3, 4, "alpha", "trusty", "amd64"),
	}, {
		v:   "1.2.3",
		err: "invalid binary version.*",
	}, {
		v:   "1.2-beta1",
		err: "invalid binary version.*",
	}, {
		v:   "1.2.3--amd64",
		err: "invalid binary version.*",
	}, {
		v:   "1.2.3-trusty-",
		err: "invalid binary version.*",
	}}

	for i, test := range parseBinaryTests {
		c.Logf("test 1: %d", i)
		got, err := version.ParseBinary(test.v)
		if test.err != "" {
			c.Assert(err, gc.ErrorMatches, test.err)
		} else {
			c.Assert(err, jc.ErrorIsNil)
			c.Assert(got, gc.Equals, test.expect)
		}
	}

	for i, test := range parseTests {
		c.Logf("test 2: %d", i)
		v := test.v + "-trusty-amd64"
		got, err := version.ParseBinary(v)
		expect := version.Binary{
			Number: test.expect,
			Series: "trusty",
			Arch:   "amd64",
		}
		if test.err != "" {
			c.Assert(err, gc.ErrorMatches, strings.Replace(test.err, "version", "binary version", 1))
		} else {
			c.Assert(err, jc.ErrorIsNil)
			c.Assert(got, gc.Equals, expect)
		}
	}
}

func (*BinarySuite) TestBinaryMarshalUnmarshal(c *gc.C) {
	for _, m := range marshallers {
		c.Logf("encoding %v", m.name)
		type doc struct {
			Version *version.Binary
		}
		// Work around goyaml bug #1096149
		// SetYAML is not called for non-pointer fields.
		bp := version.MustParseBinary("1.2.3-trusty-amd64")
		v := doc{&bp}
		data, err := m.marshal(&v)
		c.Assert(err, jc.ErrorIsNil)
		var bv doc
		err = m.unmarshal(data, &bv)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(bv, gc.DeepEquals, v)
	}
}
