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

package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	appconfig "opentelemetry-go-example/internal/config"
	apphandler "opentelemetry-go-example/internal/handler"
	apptracer "opentelemetry-go-example/internal/tracer"
	"os"
	"strconv"
)

func setup(name string, handlerFunction http.HandlerFunc, port uint16) {
	tp := apptracer.InitTracerProvider("go-" + name + "-service")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	wrappedHandler := otelhttp.NewHandler(handlerFunction, name + "-handler")

	http.Handle("/", wrappedHandler)

	log.Println("Starting ", name, " at :", port)
	portString := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(portString, nil))
}

func main() {

	// Usage
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <frontend|backend>", os.Args[0])
		os.Exit(1)
	}

	// Do config
	appconfig.SetupConfig()

	if os.Args[1] == "frontend" {
		feP, _ := strconv.Atoi(appconfig.Config["FRONTEND_PORT"])
		setup("frontend", apphandler.FrontendHandler(), uint16(feP))
	} else if os.Args[1] == "backend" {
		beP, _ := strconv.Atoi(appconfig.Config["BACKEND_PORT"])
		setup("backend", apphandler.BackendHandler(), uint16(beP))
	}
}
