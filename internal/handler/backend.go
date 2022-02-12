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
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	appconfig "opentelemetry-go-example/internal/config"
	"time"
)

func BackendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received backend request to:", r.Host, r.URL.Path, "::", r.Method)

		// Common sleep time
		duration := time.Duration(appconfig.Sleep) * time.Millisecond

		// Make spans look pretty
		time.Sleep(duration)

		// Now start a child span
		_, span := otel.Tracer("otel-go-backend").Start(r.Context(), "backend-work")

		time.Sleep(duration)
		span.AddEvent("backend-job")

		w.WriteHeader(http.StatusOK)
		responseBody := fmt.Sprintf("From backend: %s", time.Now().Local().Format("15:04:05.000"))
		w.Write([]byte(responseBody))

		time.Sleep(duration)
		span.End()
		time.Sleep(duration)
	}
}
