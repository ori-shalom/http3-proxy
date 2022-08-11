package http3proxy

import (
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/ori-shalom/http3-proxy/config"
	"io"
	"log"
	"net/http"
)

var http3Client = http.Client{Transport: &http3.RoundTripper{}}

func NewHttp3Proxy(conf config.Config) error {
	handler := proxyHandler(conf.TargetHost)

	return http.ListenAndServe(":"+conf.Port, handler)
}

func proxyHandler(targetHost string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := prepareProxyRequest(r, targetHost)

		response, err := http3Client.Do(request)
		if err != nil {
			log.Println(err)
			return
		}
		writeResponse(w, response)
	}
}

func prepareProxyRequest(r *http.Request, targetHost string) *http.Request {
	request := r.Clone(r.Context())
	request.URL.Scheme = "https"
	request.URL.Host = targetHost
	request.Host = targetHost
	request.RequestURI = ""
	return request
}

func writeResponse(w http.ResponseWriter, response *http.Response) {
	// ignore error (nothing to do if connection to server close)
	defer func() { _ = response.Body.Close() }()

	// copy response headers
	for header, values := range response.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(response.StatusCode)

	// can error only if connection is closed on either end. no point of printing such error.
	_, _ = io.Copy(w, response.Body)
}
