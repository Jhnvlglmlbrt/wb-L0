package main

import (
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error at starting: %v", err)
	}
	fmt.Println(cfg)
}
