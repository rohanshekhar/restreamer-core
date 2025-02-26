package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datarhei/core/ffmpeg"
	"github.com/datarhei/core/http/api"
	"github.com/datarhei/core/http/errorhandler"
	"github.com/datarhei/core/http/validator"
	"github.com/datarhei/core/restream"
	"github.com/datarhei/core/restream/store"

	"github.com/alecthomas/jsonschema"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func DummyRestreamer() restream.Restreamer {
	store := store.NewDummyStore(store.DummyConfig{})

	ffmpeg, _ := ffmpeg.New(ffmpeg.Config{
		Binary: "ffmpeg",
	})

	rs, _ := restream.New(restream.Config{
		Store:  store,
		FFmpeg: ffmpeg,
	})

	return rs
}

func DummyEcho() *echo.Echo {
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true
	router.HTTPErrorHandler = errorhandler.HTTPErrorHandler
	router.Logger.SetOutput(ioutil.Discard)
	router.Validator = validator.New()

	return router
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

func Request(t *testing.T, httpstatus int, router *echo.Echo, method, path string, data io.Reader) *Response {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, data)
	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	router.ServeHTTP(w, req)

	response := CheckResponse(t, w.Result())

	require.Equal(t, httpstatus, w.Code, response.Data)

	return response
}

func CheckResponse(t *testing.T, res *http.Response) *Response {
	response := &Response{
		Code: res.StatusCode,
	}

	body, err := ioutil.ReadAll(res.Body)
	require.Equal(t, nil, err)

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		err := json.Unmarshal(body, &response.Data)
		require.Equal(t, nil, err)
	} else {
		response.Data = body
	}

	if response.Code != http.StatusOK {
		if err, ok := response.Data.(api.Error); ok {
			response.Message = err.Message
		}
	}

	return response
}

func Validate(t *testing.T, datatype, data interface{}) bool {
	schema, _ := jsonschema.Reflect(datatype).MarshalJSON()

	schemaLoader := gojsonschema.NewStringLoader(string(schema))
	documentLoader := gojsonschema.NewGoLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	require.Equal(t, nil, err)
	require.Equal(t, true, result.Valid(), result.Errors())

	return true
}

func Read(t *testing.T, path string) io.Reader {
	data, err := ioutil.ReadFile(path)
	require.Equal(t, nil, err)

	return bytes.NewReader(data)
}
