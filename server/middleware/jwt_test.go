package middleware

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type jwtSuite struct{}

var _ = check.Suite(&jwtSuite{})

func (as *jwtSuite) Test_auth(c *check.C) {
	SetJwtSecret("walm-test")
	enstr, err := GenerateToken("username", "password")
	c.Assert(err, check.IsNil)

	pClaim, err := ParseToken(enstr)
	c.Assert(err, check.IsNil)
	c.Assert(pClaim.Username, check.Equals, "username")
	c.Assert(pClaim.Password, check.Equals, "password")
}
