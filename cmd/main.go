package main

import (
	"multi-site-dashboard-go/internal/rest"
)

func main() {
	s, err := rest.NewServer()
	if err != nil {
		panic(err.Error())
	}
	if err := s.Run(); err != nil {
		s.Logger.Fatal(err.Error())
	}
}