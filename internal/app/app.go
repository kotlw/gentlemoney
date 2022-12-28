package app

import (
	"fmt"
	"github.com/kotlw/gentlemoney/config"
)

func Run() {
	// Configuration
	cfg := config.Default()
	fmt.Println(cfg)

}
