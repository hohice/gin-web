package util

import (
	"testing"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type utilResSuite struct {
}

var _ = check.Suite(&utilResSuite{})

func (ars *utilResSuite) Test_GetPathParams(c *check.C) {
	param := gin.Param{
		Key:   "key_1",
		Value: "val_1",
	}
	gc := new(gin.Context)
	gc.Params = append(gc.Params, param)
	c.Assert(gc.Param("key_1"), check.Equals, "val_1")

	val, err := GetPathParams(gc, []string{"key_1"})
	c.Assert(err, check.IsNil)
	//c.Assert(val, check.IsNil)
	c.Assert(val["key_1"], check.Equals, "val_1")

}
