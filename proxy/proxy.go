package proxy

import (
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"log"
	"net/http"
)

var http3Client = http.Client{Transport: &http3.RoundTripper{}}

func NewHttp3Proxy(conf Config) error {
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
		// ignore error (nothing to do if connection to server close)
		defer func() { _ = response.Body.Close() }()

		writeResponse(w, response)
	}
}

func prepareProxyRequest(r *http.Request, targetHost string) *http.Request {
	request := r.Clone(r.Context())
	request.URL.Scheme = "https"
	request.URL.Host = targetHost
	request.Host = targetHost
	request.RequestURI = ""
	request.Proto = ""
	return request
}

func writeResponse(w http.ResponseWriter, response *http.Response) {

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
