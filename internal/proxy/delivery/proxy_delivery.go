package delivery

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yletamitlu/proxy/internal/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	if request.Method == http.MethodConnect {
		err := pd.proxyUcase.HandleHttpsRequest(writer, request)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		err := pd.proxyUcase.SaveReqToDB(request, "http")
		if err != nil {
			logrus.Error(err)
		}

		err = pd.proxyUcase.HandleHttpRequest(writer, request)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (pd *ProxyDelivery) GetAllRequestsHandler(writer http.ResponseWriter, request *http.Request) {
	requests, err := pd.proxyUcase.GetAllRequests()
	if err != nil {
		logrus.Error(err)
	}

	var result string
	for _, request := range requests {
		result += "\nId: " + strconv.FormatInt(request.Id, 16)  + "\nMethod: " + request.Method +
			"\nScheme: " + request.Scheme + "\nPath: " + request.Path + "\nHost: " + request.Host +
			"\nBody: " + request.Body + "\n"
	}
	_, err = io.Copy(writer, strings.NewReader(result))
	if err != nil {
		logrus.Info(err)
	}
}

func (pd *ProxyDelivery) GetRequestHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 32)
	foundReq, err := pd.proxyUcase.GetRequest(id)
	if err != nil {
		logrus.Error(err)
	}

	result := "Id: " + strconv.FormatInt(foundReq.Id, 16)  + "\nMethod: " + foundReq.Method +
		"\nScheme: " + foundReq.Scheme + "\nPath: " + foundReq.Path + "\nHost: " + foundReq.Host +
		"\nBody: " + foundReq.Body + "\n"

	_, err = io.Copy(writer, strings.NewReader(result))
	if err != nil {
		logrus.Info(err)
	}
}

func (pd *ProxyDelivery) RepeatRequestHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 32)
	foundReq, err := pd.proxyUcase.GetRequest(id)
	if err != nil {
		logrus.Error(err)
	}

	pd.HandleRequest(writer, &http.Request{
		Method: foundReq.Method,
		URL: &url.URL{
			Scheme: foundReq.Scheme,
			Host:   foundReq.Host,
			Path:   foundReq.Path,
		},
		Header: foundReq.Headers,
		Body:   ioutil.NopCloser(strings.NewReader(foundReq.Body)),
		Host:   request.Host,
	})
}
