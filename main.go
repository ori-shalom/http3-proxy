package main

import (
	"encoding/json"
	"fmt"
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const targetHostEnvVarName = "TARGET_HOST"

func main() {
	if err := handleMain(); err != nil {
		log.Println("Failed to stat proxy.")
		log.Fatal(err)
	}
}

func handleMain() error {
	client := http.Client{
		Transport: &http3.RoundTripper{},
	}
	targetHost := os.Getenv(targetHostEnvVarName)
	if targetHost == "" {
		return fmt.Errorf("environment variable %s is empty", targetHostEnvVarName)
	}

	handler := proxyHandler(targetHost, client)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		return err
	}
	return nil
}

func proxyHandler(targetHost string, client http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		request := r.Clone(r.Context())
		request.URL.Host = targetHost
		request.URL.Scheme = "https"
		request.Host = targetHost
		request.RequestURI = ""

		response, err := client.Do(request)
		if err != nil {
			logError(err)
			return
		}
		now := time.Now()
		logMap(map[string]any{
			"@timestamp": now,
			"duration":   now.Sub(start).String(),
			"status":     response.StatusCode,
			"method":     request.Method,
			"path":       request.URL.Path,
		})
		for header, values := range response.Header {
			for _, value := range values {
				w.Header().Add(header, value)
			}
		}
		w.WriteHeader(response.StatusCode)
		defer func() {
			if err := response.Body.Close(); err != nil {
				logError(err)
			}
		}()

		if _, err := io.Copy(w, response.Body); err != nil {
			logError(err)
		}
	}
}

func logMap(jsonMap map[string]any) {
	if err := json.NewEncoder(log.Writer()).Encode(jsonMap); err != nil {
		log.Println(err)
	}
}

func logError(err error) {
	logMap(map[string]any{
		"error": err.Error(),
	})
}
