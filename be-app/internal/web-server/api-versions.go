package webserver

type ApiVersion uint16

// array containing supported api versions starting from 1
const (
	v1 ApiVersion = 1 << iota
)
