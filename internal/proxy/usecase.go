package proxy

import (
	"github.com/yletamitlu/proxy/internal/models"
	"net/http"
)

type ProxyUsecase interface {
	HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) error
	HandleHttpsRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) error
	DoHttpRequest(httpRequest *http.Request) (*http.Response, error)

	SaveReqToDB(request *http.Request, scheme string) error
	GetRequest(id int64) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
