// Command shop запускает HTTP-сервис магазина мерча Авито.
package main

import (
	"fmt"
	"os"

	"github.com/ealia28/avito-merch2025internship/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("server port: %s\n", cfg.Server.Port)
}
