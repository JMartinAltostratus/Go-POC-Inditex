package main

import (
	bootstrap "github.com/JMartinAltostratus/Go-POC-Inditex/cmd/api/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
