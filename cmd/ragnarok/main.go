package main

import (
	"context"
	"log"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/attack/memory"
)

func main() {
	opts := attack.Opts{
		"size": 256 * memory.MiB,
	}

	a, err := attack.New(memory.AllocID, opts)

	if err != nil {
		log.Fatal(err)
	}

	if err := a.Apply(context.TODO()); err != nil {
		log.Fatal(err)
	}
	defer a.Revert()
	time.Sleep(30 * time.Second)
}
