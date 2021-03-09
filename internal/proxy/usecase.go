package proxy

import "net/http"

type ProxyUsecase interface {
	HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) error
	HandleHttpsRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) error
	DoHttpRequest(httpRequest *http.Request) (*http.Response, error)
}
