package httputil

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/gin-gonic/gin"
)

func RunGinServerHandler(ctrl interface{}, method string, ct ContentType, payload string, handlerName string) (*httptest.ResponseRecorder, error) {
	body := bytes.NewBufferString(payload)
	req, err := http.NewRequest(method, "/", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", string(ct))
	resp := httptest.NewRecorder()

	router := gin.Default()

	handler, ok := reflect.ValueOf(ctrl).MethodByName(handlerName).Interface().(func(*gin.Context))
	if !ok {
		return nil, errors.New("ctrl must a func(*gin.Context)")
	}
	reflect.ValueOf(router).MethodByName(method).Call([]reflect.Value{reflect.ValueOf("/"), reflect.ValueOf(handler)})

	router.ServeHTTP(resp, req)

	return resp, nil
}
