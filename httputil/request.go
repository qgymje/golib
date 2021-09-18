package httputil

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ContentType string

const (
	ContentTypeForm      ContentType = "application/x-www-form-urlencoded"
	ContentTypeMultipart ContentType = "multipart/form-data"
	ContentTypeJSON      ContentType = "application/json"
)

func DoRequestWithReader(httpClient *http.Client, urlStr string, method string, contentType ContentType, body io.Reader) ([]byte, error) {
	var (
		err  error
		resp *http.Response
	)
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest(method, urlStr, bytes.NewBuffer(bodyData))
		req.Header.Set("Content-Type", string(contentType))
		if resp, err = httpClient.Do(req); err == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func DoRequest(httpClient *http.Client, urlStr string, method string, contentType ContentType, params url.Values) ([]byte, error) {
	var (
		err  error
		resp *http.Response
	)

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest(method, urlStr, bytes.NewBufferString(params.Encode()))
		req.Header.Set("Content-Type", string(contentType))
		if resp, err = httpClient.Do(req); err == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func NewUploadRequest(urlStr string, params map[string]string, fieldName, filename string) (*http.Request, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fi.Name())
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	for k, v := range params {
		_ = writer.WriteField(k, v)
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
