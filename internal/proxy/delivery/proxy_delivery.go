package delivery

import (
	"github.com/yletamitlu/proxy/internal/proxy"
	"log"
	"net/http"
	"net/url"
)

type ProxyDelivery struct {
	proxyUcase proxy.ProxyUsecase
}

func NewProxyDelivery(proxyUsecase proxy.ProxyUsecase) *ProxyDelivery {
	return &ProxyDelivery{
		proxyUcase: proxyUsecase,
	}
}

func (pd *ProxyDelivery) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	parsedUrl, _ := url.Parse(request.RequestURI)
	if parsedUrl.Scheme == "http" {
		err := pd.proxyUcase.HandleHttpRequest(writer, request)
		if err != nil {
			log.Fatal(err)
		}
	} else {

	}
}
