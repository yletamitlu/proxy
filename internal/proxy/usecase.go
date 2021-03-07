package proxy

import "net/http"

type ProxyUsecase interface {
	HandleRequest(httpRequest *http.Request) (*http.Response, error)
}
