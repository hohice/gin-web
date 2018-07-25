package ex

import (
	check "gopkg.in/check.v1"
)

type msgSuite struct {
}

var _ = check.Suite(&msgSuite{})

func (ms *msgSuite) Test_GetMsg(c *check.C) {
	str := GetMsg(SUCCESS)
	c.Assert(str, check.Equals, "ok")
	str = GetMsg(0)
	c.Assert(str, check.Equals, "Internal Server error")
}
