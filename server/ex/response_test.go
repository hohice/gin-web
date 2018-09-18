package ex

import (
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type apiResSuite struct {
}

var _ = Suite(&apiResSuite{})

func (ars *apiResSuite) TestResponse_returnBadRequest(c *C) {
	_, br := ReturnBadRequest()
	c.Assert(br.Code, Equals, http.StatusBadRequest)
	c.Assert(br.Message, Equals, "Invalid Name supplied!")
}

func (ars *apiResSuite) TestResponse_returnInternalServerError(c *C) {
	err := myerr{msg: "server error!"}
	_, ise := ReturnInternalServerError(err)

	c.Assert(ise.Code, Equals, http.StatusInternalServerError)
	c.Assert(ise.Message, Equals, "Internal Server error: server error!")
}

func (ars *apiResSuite) TestResponse_returnConfigExistError(c *C) {
	_, ise := ReturnConfigExistError()
	c.Assert(ise.Code, Equals, ERROR_CONFIG_EXIST)
}

func (ars *apiResSuite) TestResponse_returnConfigNotExistError(c *C) {
	_, ise := ReturnConfigNotExistError()
	c.Assert(ise.Code, Equals, ERROR_CONFIG_NOT_EXIST)
}

func (ars *apiResSuite) TestResponse_returnLimitError(c *C) {
	_, ise := ReturnLimitError()
	c.Assert(ise.Code, Equals, ERROR_LIMIT)
}

func (ars *apiResSuite) TestResponse_ReturnOK(c *C) {
	ok, _ := ReturnOK()
	c.Assert(ok, Equals, SUCCESS)
}

type myerr struct {
	msg string
}

func (me myerr) Error() string {
	return me.msg
}
