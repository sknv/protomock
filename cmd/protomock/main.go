package main

import (
	"flag"
	stdlog "log"

	_ "go.uber.org/automaxprocs"

	"github.com/sknv/protomock/internal/config"
)

func main() {
	configPath := config.FilePathFlag()
	flag.Parse() //nolint:wsl // process a variable above

	cfg, err := config.Parse(*configPath)
	fatalIfError(err)

	err = run(cfg)
	fatalIfError(err)
}

func run(cfg *config.Config) error {
	stdlog.Printf("config is: %+v", cfg)

	return nil
}

func fatalIfError(err error) {
	if err != nil {
		stdlog.Fatal(err)
	}
}
