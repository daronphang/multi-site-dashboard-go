package main

import (
	"multi-site-dashboard-go/internal/rest"
)

func main() {
	s, err := rest.NewServer()
	if err != nil {
		panic(err.Error())
	}
	s.Run()
}