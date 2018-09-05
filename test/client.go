package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func Get(uri string, router *gin.Engine) (int, []byte) {
	req := httptest.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	return result.StatusCode, body
}

func PostJson(uri string, param map[string]interface{}, router *gin.Engine) (int, []byte) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	return result.StatusCode, body
}

func PostForm(uri string, param map[string]string, router *gin.Engine) (int, []byte) {
	req := httptest.NewRequest("POST", uri+ParseToStr(param), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	return result.StatusCode, body
}

func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}
