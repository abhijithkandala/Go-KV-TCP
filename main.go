package main

import (
	"log"

	"github.com/abhikred/go-kv/internal/server"
	"github.com/abhikred/go-kv/internal/store"
)

func main() {
	s := store.NewActorStore()
	err := server.ListenAndServe("127.0.0.1:6380", s)
	if err != nil {
		log.Fatal(err)
	}
}
