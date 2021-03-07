package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yletamitlu/proxy/internal/proxy"
)

type ProxyRepos struct {
	conn *sqlx.DB
}

func NewProxyRepos(conn *sqlx.DB) *proxy.ProxyRepository {
	return &ProxyRepos{
		conn: conn,
	}
}
