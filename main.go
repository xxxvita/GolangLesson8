package main

import (
	"fmt"
	"lesson8/configurator/configurator"
	"os"
)

var s *configurator.WebConfig

func main() {
	fmt.Println("Программа загружает конфигурацию с отдельным модулем")

	var err error

	s, err = configurator.Load("config.json")
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Загружена конфигурация: %#v\n", s)
}
