// Copyright 2022 Sheldon Lobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	appconfig "opentelemetry-go-example/internal/config"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var backendClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

func FrontendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received frontend request to:", r.Host, r.URL.Path, "::", r.Method)

		// Common sleep time
		duration := time.Duration(appconfig.Sleep) * time.Millisecond

		// Make spans look pretty
		time.Sleep(duration)

		// Now start a child span
		_, span := otel.Tracer("otel-go-frontend").Start(r.Context(), "frontend-work")
		defer span.End()

		// Add an event
		time.Sleep(duration)
		span.AddEvent("frontend-job")

		time.Sleep(duration)

		// Make wrapped call to backend
		backendURL := fmt.Sprintf("http://%s:%s/", appconfig.Config["BACKEND_SERVER"],
			appconfig.Config["BACKEND_PORT"])
		backendRequest, err := http.NewRequestWithContext(r.Context(), http.MethodPost, backendURL, nil)
		if err != nil {
			http.Error(w, "Bad backend request", http.StatusInternalServerError)
			return
		}
		backendResponse, err := backendClient.Do(backendRequest)
		if err != nil {
			log.Printf("Bad backend request: %s; %v", backendURL, err)
			http.Error(w, "Bad backend response", http.StatusBadGateway)
			return
		}
		defer backendResponse.Body.Close()

		backendBody := "Bad backend"
		if backendResponse.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(backendResponse.Body)
			if err != nil {
				log.Printf("Could not read backend response: %v", err)
				http.Error(w, "Bad backend response", http.StatusBadGateway)
				return
			}
			backendBody = string(bodyBytes)
		}

		time.Sleep(duration)

		w.WriteHeader(http.StatusOK)
		responseBody := fmt.Sprintf("From frontend: %s [%s]\n",
			time.Now().Local().Format("15:04:05.000"), backendBody)
		w.Write([]byte(responseBody))

		time.Sleep(duration)
	}
}
