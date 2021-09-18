package logger

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qgymje/golib/httputil"
)

type DingTalkLogger struct {
	token      string
	app        string
	httpClient *http.Client
}

func NewDingTalkLogger(app string, token string) *DingTalkLogger {
	return &DingTalkLogger{
		app:        app,
		token:      token,
		httpClient: httputil.CreateHttpClient(3*time.Second, 10, 10),
	}
}

func (l *DingTalkLogger) Error(format string, args ...interface{}) {
	msg := `{
            "msgtype": "markdown",
            "markdown": {
                "title": "%s",
				"text": "[%s]:%s"
            }
        }`
	text := fmt.Errorf(format, args...)
	body := fmt.Sprintf(msg, l.app, l.app, text)
	log.Printf("%s\n", body)

	resp, err := httputil.DoRequestWithReader(l.httpClient, l.token, http.MethodPost, httputil.ContentTypeJSON, bytes.NewBufferString(body))
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s", resp)
}

func (l *DingTalkLogger) Fatal(format string, args ...interface{}) {
	log.Printf("[Fatal]"+format, args...)
}

func (l *DingTalkLogger) Debug(format string, args ...interface{}) {
	log.Printf("[Debug]"+format, args...)
}

func (l *DingTalkLogger) Info(format string, args ...interface{}) {
	log.Printf("[Info]"+format, args...)
}
