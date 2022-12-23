package app

import (
	"fmt"
	"gentlemoney/config"
)

func Run() {
	// Configuration
	cfg := config.Default()
	fmt.Println(cfg)

}
