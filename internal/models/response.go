package models

type Response struct {
	Status  int
	Headers map[string][]string
	Body    string
}
