package proxy

import (
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"log"
	"net/http"
)

func NewHttp3Proxy(conf Config) error {
	handler := proxyHandler(conf.TargetHost)

	return http.ListenAndServe(":"+conf.Port, handler)
}

func proxyHandler(targetHost string) http.HandlerFunc {
	http3Client := http.Client{
		Transport: &http3.RoundTripper{},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := prepareProxyRequest(r, targetHost)
		log.Println(request)

		response, err := http3Client.Do(request)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(
			response.StatusCode,
			response.ContentLength,
			response.Proto)

		// ignore error (nothing to do if connection to server close)
		defer func() {
			if err = response.Body.Close(); err != nil {
				log.Println(err)
			}
		}()

		// copy response headers
		for header, values := range response.Header {
			for _, value := range values {
				w.Header().Add(header, value)
			}
		}

		// can error only if connection is closed on either end. no point of printing such error.
		if _, err = io.Copy(w, response.Body); err != nil {
			log.Println(err)
		}

		w.WriteHeader(response.StatusCode)
	}
}

func prepareProxyRequest(r *http.Request, targetHost string) *http.Request {
	request := r.Clone(r.Context())
	request.RemoteAddr = ""
	request.Proto = ""
	request.ProtoMajor = 0
	request.ProtoMinor = 0
	request.RequestURI = ""
	request.TLS = nil
	request.Close = false
	request.ContentLength = 0
	request.Header.Set("Host", targetHost)
	request.Header.Del("X-Forwarded-For")
	request.Header.Del("X-Forwarded-Proto")
	request.URL.Scheme = "https"
	request.URL.Host = targetHost
	request.Host = targetHost
	return request
}
