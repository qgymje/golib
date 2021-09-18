package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qgymje/golib/provider"
	"github.com/qgymje/golib/utility"
)

type requestlog struct {
	RequestId       string `json:"request_id"`
	RequestURI      string `json:"request_uri"`
	RequestBody     string `json:"request_body"`
	IP              string `json:"ip"`
	Method          string `json:"method"`
	UserAgent       string `json:"user_agent"`
	RequestTime     string `json:"request_time"`
	RequestDuration string `json:"request_duration"`
	ResponseStatus  int    `json:"response_status"`
	ResponseBody    string `json:"response_body"`
	ResponseError   string `json:"response_error"`
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *BodyLogWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func BodyLogger(logger provider.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		blw := &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		req := ctx.Request
		var body []byte
		if req.Body != nil {
			body, _ = ioutil.ReadAll(req.Body)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		ctx.Next()

		/*
			requri, _ := url.Parse(req.URL.RequestURI())
			uv := requri.Query()
			uv.Del("request_id")
			uv.Del("response_error")
			requri.RawQuery = uv.Encode()
		*/

		logdata := &requestlog{
			RequestId: req.FormValue("request_id"),
			//RequestURI:      requri.String(),
			RequestURI:      ctx.Request.URL.RequestURI(),
			Method:          req.Method,
			RequestBody:     strings.TrimSpace(string(body)),
			UserAgent:       req.UserAgent(),
			IP:              utility.GetRemoteIP(req),
			RequestTime:     start.Format(time.RFC3339),
			RequestDuration: time.Since(start).String(),
			ResponseStatus:  ctx.Writer.Status(),
			ResponseBody:    blw.Body.String(),
			ResponseError:   req.FormValue("response_error"),
		}
		jsondata, _ := json.Marshal(logdata)
		logger.Info(string(jsondata))
	}
}
