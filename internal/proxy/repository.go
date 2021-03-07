package proxy

import "github.com/yletamitlu/proxy/internal/models"

type ProxyRepository interface {
	InsertInto(request *models.Request) error
	GetRequest(id int) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
