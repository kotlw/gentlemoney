package app

import (
	"fmt"
	"github.com/kotlw/gentlemoney/config"
	"github.com/kotlw/gentlemoney/internal/tui"
)

func Run() {
	// Configuration
	cfg := config.Default()
	fmt.Println(cfg)
	
	// Terminal user interface.
	t := tui.New()
	if err := t.Run(); err != nil {
		t.Stop()
	}
}
