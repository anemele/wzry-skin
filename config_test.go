package main

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fmt.Println(Config)
	fmt.Println(Config["savepath"])
	fmt.Println(Config["none"] == "")
}
