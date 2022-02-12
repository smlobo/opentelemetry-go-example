# opentelemetry-go-example
An example app where two Golang services collaborate on an http request. The timing of these requests are recorded using the OpenTelemetry SDK.

## Setup
* `make module build`

## Jaeger Exporter - view traces on the Jaeger UI
* `docker run -d -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one`
* Execute the `frontend` and `backend` services on separate terminals.
  For example:
  * `./bin/opentelemetry-go-example backend`
  * `./bin/opentelemetry-go-example frontend`
* Generate traces:
  * `curl http://localhost:8081`
* ![Example Jaeger Trace](images/jaeger-screenshot.png?raw=true "Example Jaeger Trace")

## Zipkin Exporter - view traces on the Zipkin UI
* `docker run -d -p 9411:9411 openzipkin/zipkin`
* Execute the `frontend` and `backend` services on separate terminals.
  For example:
    * `EXPORTER=zipkin ./bin/opentelemetry-go-example backend`
    * `EXPORTER=zipkin ./bin/opentelemetry-go-example frontend`
* Generate traces:
    * `curl http://localhost:8081`
* ![Example Zipkin Trace](images/zipkin-screenshot.png?raw=true "Example Zipkin Trace")

# References
* https://github.com/smlobo/zipkin-go-example
* https://pkg.go.dev/go.opentelemetry.io
