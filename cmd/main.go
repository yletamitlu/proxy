package main

import (
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	proxyDelivery "github.com/yletamitlu/proxy/internal/proxy/delivery"
	proxyRepos "github.com/yletamitlu/proxy/internal/proxy/repository"
	proxyUcase "github.com/yletamitlu/proxy/internal/proxy/usecase"
	"net/http"
	"os"
)

func main() {
	conn, err := sqlx.Connect("pgx",
		"postgres://" + os.Getenv("DB_USER") + ":techdb@localhost:5432/" + os.Getenv("DB_NAME"))
	if err != nil {
		logrus.Info(err)
	}

	conn.SetMaxOpenConns(8)
	conn.SetMaxIdleConns(8)

	if err := conn.Ping(); err != nil {
		logrus.Info(err)
	}

	defer conn.Close()

	proxyR := proxyRepos.NewProxyRepos(conn)
	proxyU := proxyUcase.NewProxyUcase(proxyR)

	proxyD := proxyDelivery.NewProxyDelivery(proxyU)

	server := &http.Server{
		Addr: ":" + os.Getenv("PROXY_PORT"),
		Handler: http.HandlerFunc(proxyD.HandleRequest),
	}

	go func() {
		logrus.Info(server.ListenAndServe())
	}()

	router := mux.NewRouter()
	router.HandleFunc("/requests", proxyD.GetAllRequestsHandler)
	router.HandleFunc("/requests/{id:[0-9]+}", proxyD.GetRequestHandler)
	router.HandleFunc("/repeat/{id:[0-9]+}", proxyD.RepeatRequestHandler)
	http.Handle("/", router)
	logrus.Info(http.ListenAndServe(":" + os.Getenv("REPEATER_PORT"), nil))
}
