package main

import (
	log "github.com/sirupsen/logrus"
)

func panicOn(err error) {
	if err != nil {
		log.Fatalf("Panic: %v", err)
		panic(err)
	}
}
