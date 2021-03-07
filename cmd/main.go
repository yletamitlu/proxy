package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"github.com/yletamitlu/internal/proxy/repository"
)

func main() {
	conn, err := sqlx.Connect("pgx", "postgres://proxyuser:techdb@localhost:5432/proxydb")
	if err != nil {
		log.Fatal(err)
	}

	conn.SetMaxOpenConns(8)
	conn.SetMaxIdleConns(8)

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	proxyR := proxyRepos.NewUserRepository(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Server started...")
	})

	http.ListenAndServe(":80", nil)
}

