package repository

import (
	"encoding/json"
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
	bs, err := json.Marshal(request.Headers)
	if err != nil {
		return err
	}

	if err := pr.conn.QueryRow(
		`INSERT INTO requests(method, scheme, host, path, headers, body) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		request.Method, request.Scheme, request.Host, request.Path, string(bs), request.Body).Scan(&request.Id);
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

	rows, err := pr.conn.Queryx(`SELECT * from requests`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		reqMap := make(map[string]interface{})
		err = rows.MapScan(reqMap)

		headersRaw := reqMap["headers"].([]byte)

		var headers map[string][]string
		err := json.Unmarshal(headersRaw, &headers)
		if err != nil {
			return nil, err
		}

		requests = append(requests, &models.Request{
			Id:      reqMap["id"].(int64),
			Method:  reqMap["method"].(string),
			Scheme:  reqMap["scheme"].(string),
			Host:    reqMap["host"].(string),
			Path:    reqMap["path"].(string),
			Headers: headers,
			Body:    reqMap["body"].(string),
		})
	}

	return requests, nil
}
