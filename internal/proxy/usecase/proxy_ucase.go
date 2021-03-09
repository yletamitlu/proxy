package usecase

import (
	"github.com/yletamitlu/proxy/internal/proxy"
	"io"
	"log"
	"net/http"
	"strings"
)

type ProxyUcase struct {
	proxyRepo proxy.ProxyRepository
}

func NewProxyUcase(repos proxy.ProxyRepository) proxy.ProxyUsecase {
	return &ProxyUcase{
		proxyRepo: repos,
	}
}

func (pu *ProxyUcase) HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) error {
	proxyResponse, err := pu.DoHttpRequest(interceptedHttpRequest)

	if err != nil {
		log.Fatal(err)
	}

	var headers string
	for header, values := range proxyResponse.Header {
		for _, value := range values {
			headers += header  + ": " + value + "\n"
		}
	}

	_, err = io.Copy(writer, io.MultiReader(strings.NewReader(headers + "\n"), proxyResponse.Body))

	if err != nil {
		log.Fatal(err)
	}

	defer proxyResponse.Body.Close()
	return nil
}

func (pu *ProxyUcase) DoHttpRequest(HttpRequest *http.Request) (*http.Response, error) {
	proxyRequest, err := http.NewRequest(HttpRequest.Method, HttpRequest.URL.String(), HttpRequest.Body)
	if err != nil {
		return nil, err
	}

	proxyRequest.Header = HttpRequest.Header

	proxyClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	proxyResponse, err := proxyClient.Do(proxyRequest)
	if err != nil {
		return nil, err
	}

	return proxyResponse, nil
}
