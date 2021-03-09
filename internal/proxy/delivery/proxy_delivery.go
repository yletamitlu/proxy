package delivery

import (
	"github.com/sirupsen/logrus"
	"github.com/yletamitlu/proxy/internal/proxy"
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
			logrus.Info(err)
		}
	} else {
		err := pd.proxyUcase.HandleHttpsRequest(writer, request)
		if err != nil {
			logrus.Info(err)
		}
	}
}
