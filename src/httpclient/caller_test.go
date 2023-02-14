package httpclient_test

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sghaida/go-stuff/src/cauth"
	"github.com/sghaida/go-stuff/src/httpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var url = "https://crudcrud.com"
var endpointPath string

func TestMain(m *testing.M) {
	resp, _ := http.Get(url)
	defer func() {
		_ = resp.Body.Close()
	}()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	// define selector
	matcher := goquery.Single(
		"body > " +
			"section.hero.is-dark.is-medium.is-bold > " +
			"div > " +
			"div > " +
			"div.endpoint-url.notification.is-light.is-family-code.is-size-7-mobile",
	)
	// find the item and extract the text
	endpointPath = doc.FindMatcher(matcher).Text()
	// clean up the text
	endpointPath = strings.TrimSpace(endpointPath)
	endpointPath = strings.Replace(endpointPath, "https://crudcrud.com", "", -1)
	endpointPath = fmt.Sprintf("%s/crudOps", endpointPath)
	endpointPath = strings.TrimPrefix(endpointPath, "/")

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCaller_CallWithTimeout(t *testing.T) {

	url := "https://crudcrud.com"

	token := "some_jwt_token"
	headers := map[string]string{
		"Accept":       `application/json`,
		"Content-Type": `application/json`,
	}
	client := createClient(t, url, endpointPath, token, headers)
	type reqPayload struct {
		ID   string `json:"_id,omitempty"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	req := reqPayload{
		Name: "test",
		Age:  11,
	}
	var response []reqPayload

	payload, _ := json.Marshal(req)

	t.Run("post data", func(t *testing.T) {

		resp, err := client.Call(httpclient.POST, nil, nil, payload)
		if err != nil {
			assert.Failf(t, "expected to call successfully, recieved an error %s", err.Error())
		}
		body, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			assert.Failf(t, "expected to get the body successfully, got %s", err.Error())
		}
		b := string(body)
		fmt.Println(b)
	})

	t.Run("get data", func(t *testing.T) {

		resp, err := client.Call(httpclient.GET, nil, nil, nil)
		if err != nil {
			assert.Failf(t, "expected to call successfully, recieved an error %s", err.Error())
		}
		body, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			assert.Failf(t, "expected to get the body successfully, got %s", err.Error())
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			assert.Failf(t, "expected to receive response, got %s", err.Error())
		}
	})

	t.Run("put data", func(t *testing.T) {

		path := fmt.Sprintf("%s/%s", endpointPath, response[0].ID)
		client := createClient(t, url, path, token, headers)

		req.Name = "test2"
		reqBody, _ := json.Marshal(req)

		resp, err := client.Call(httpclient.PUT, nil, nil, reqBody)
		if err != nil {
			assert.Failf(t, "expected to call successfully, recieved an error %s", err.Error())
		}
		body, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			assert.Failf(t, "expected to get the body successfully, got %s", err.Error())
		}
		fmt.Println(string(body))
	})

	t.Run("delete data", func(t *testing.T) {
		path := fmt.Sprintf("%s/%s", endpointPath, response[0].ID)
		client := createClient(t, url, path, token, headers)

		resp, err := client.Call(httpclient.DELETE, nil, nil, nil)
		if err != nil {
			assert.Failf(t, "expected to call successfully, recieved an error %s", err.Error())
		}
		body, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			assert.Failf(t, "expected to get the body successfully, got %s", err.Error())
		}
		fmt.Println(string(body))
	})
}

func createClient(t *testing.T, url string, path string, token string, headers map[string]string) *httpclient.Caller {
	config, err := httpclient.NewConfig().
		WithHost(url).
		WithRoute(path).
		WithTimeout(1 * time.Second).
		WithHeaders(headers).
		Build()
	if err != nil {
		assert.Failf(t, "expected config creation to succeed, got error", err.Error())
	}
	auth := cauth.NewJWTAuth(token)
	client, err := httpclient.NewHTTPCaller(config).WithAuth(auth)
	if err != nil {
		assert.Failf(t, "expected to create client, recieved %s", err.Error())
	}
	return client
}
