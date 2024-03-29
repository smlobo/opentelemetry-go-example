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
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"io/ioutil"
	"log"
	"net/http"
	appconfig "opentelemetry-go-example/internal/config"
	"time"
)

func FrontendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received frontend request to:", r.Host, r.URL.Path, "::", r.Method)

		// Common sleep time
		duration := time.Duration(appconfig.Sleep) * time.Millisecond

		// Make spans look pretty
		time.Sleep(duration)

		// Now start a child span
		_, span := otel.Tracer("otel-go-frontend").Start(r.Context(), "frontend-work")

		// Add an event
		time.Sleep(duration)
		span.AddEvent("frontend-job")

		time.Sleep(duration)

		// Make wrapped call to backend
		backendURL := fmt.Sprintf("http://%s:%s/", appconfig.Config["BACKEND_SERVER"],
			appconfig.Config["BACKEND_PORT"])
		backendResponse, err := otelhttp.Post(r.Context(), backendURL, "text/html", nil)
		if err != nil {
			log.Fatal("Bad backend request:", backendURL, ";", err)
		}

		backendBody := "Bad backend"
		if backendResponse.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(backendResponse.Body)
			if err != nil {
				log.Fatal(err)
			}
			backendBody = string(bodyBytes)
		}
		backendResponse.Body.Close()

		time.Sleep(duration)

		w.WriteHeader(http.StatusOK)
		responseBody := fmt.Sprintf("From frontend: %s [%s]\n",
			time.Now().Local().Format("15:04:05.000"), backendBody)
		w.Write([]byte(responseBody))

		span.End()
		time.Sleep(duration)
	}
}
