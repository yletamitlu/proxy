package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/proxy/internal/models"
	"github.com/yletamitlu/proxy/internal/proxy"
)

type ProxyRepos struct {
	conn *sqlx.DB
}

func NewProxyRepos(conn *sqlx.DB) proxy.ProxyRepository {
	return &ProxyRepos{
		conn: conn,
	}
}

func (pr *ProxyRepos) InsertInto(request *models.Request) error {
	if _, err := pr.conn.Exec(
		`INSERT INTO requests(method, protocol, host, path, headers, body) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
				request.Method, request.Protocol, request.Host, request.Path, request.Headers, request.Body);
		err != nil {
		return err
	}

	return nil
}

func (pr *ProxyRepos) GetRequest(id int) (*models.Request, error) {
	req := &models.Request{}

	if err := pr.conn.Get(req,
		`SELECT * from requests where id = $1`, id);
		err != nil {
		return nil, err
	}

	return req, nil
}
func (pr *ProxyRepos) GetAllRequests() ([]*models.Request, error) {
	var requests []*models.Request

	if err := pr.conn.Select(requests, `SELECT * from requests`); err != nil {
		return nil, err
	}

	return requests, nil
}
