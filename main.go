package main

import (
	s "github.com/krishnakantha1/expenseTrackerIngestion/server"
)

func main() {
	server := s.NewServer()
	server.Init()
}
