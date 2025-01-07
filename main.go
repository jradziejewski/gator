package main

import (
	"fmt"

	"github.com/jradziejewski/gator/internal/config"
)

func main() {
	conf, _ := config.Read()
	conf.SetUser("kubson")

	updatedConf, _ := config.Read()
	fmt.Println(updatedConf)
}
