package main

import (
	"github.com/spf13/viper"
	"github.com/xh/gorder/common/config"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("%v", viper.Get("order"))

}
