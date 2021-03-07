package models

type Request struct {
	Id       int
	Method   string
	Protocol string
	Host     string
	Path     string
	Headers  map[string][]string
	Body     string
}
