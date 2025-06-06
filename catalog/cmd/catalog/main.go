package main

import (
	"log"
	"time"

	"github.com/NuZard84/go_microservices/catalog"
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

	var r catalog.Repository
	retry.ForeverSleep(time.Second*2, func(attempt int) error {
		var err error
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("Linsting on port 8080 ..")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))

}
