package main

import (
	"github.com/ori-shalom/http3-proxy/config"
	"github.com/ori-shalom/http3-proxy/http3proxy"
	"log"
)

func main() {
	if err := handleMain(); err != nil {
		log.Fatal(err)
	}
}

func handleMain() error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	return http3proxy.NewHttp3Proxy(conf)
}
