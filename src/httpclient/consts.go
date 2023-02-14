package httpclient

import (
	"net/http"
)

type HttpMethod string

const (
	HEAD   HttpMethod = http.MethodHead
	GET    HttpMethod = http.MethodGet
	POST   HttpMethod = http.MethodPost
	PUT    HttpMethod = http.MethodPut
	DELETE HttpMethod = http.MethodDelete
)
