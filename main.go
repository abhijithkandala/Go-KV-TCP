package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/abhikred/go-kv/internal/server"
	"github.com/abhikred/go-kv/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	s := store.NewActorStore()
	defer s.Close()

	err := server.ListenAndServe(ctx, "127.0.0.1:6380", s)
	if err != nil {
		log.Fatal(err)
	}
}
