package main

import (
	"log"
	"time"

	"github.com/NuZard84/go_microservices/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}

	var r account.Repository
	retry.ForeverSleep(time.Second*2, func(attempt int) error {
		var err error
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("Linsting on port 8080 ..")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))

}
