package main

import (
	"github.com/ori-shalom/http3-proxy/proxy"
	"log"
)

func main() {
	if err := handleMain(); err != nil {
		log.Fatal(err)
	}
}

func handleMain() error {
	conf, err := proxy.LoadConfig()
	if err != nil {
		return err
	}

	return proxy.NewHttp3Proxy(conf)
}
