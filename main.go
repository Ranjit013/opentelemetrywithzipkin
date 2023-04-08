package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"opentelemetrywithzipkin/config"
	"opentelemetrywithzipkin/utils"
	"os"
	"time"
)

var endpointUrl = "http://localhost:9411/api/v2/spans"
var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

// initTracer creates a new trace provider instance and registers it as global trace provider.

var provider trace.Tracer

func main() {

	tracer, err := config.InitTracer(endpointUrl)
	if err != nil {
		return
	}

	otel.SetTracerProvider(tracer)
	provider = tracer.Tracer("With Zipkin")
	router := mux.NewRouter()
	router.HandleFunc("/", zipkinCodeHandler)
	/*ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	url := flag.String("zipkin", endpointUrl, "zipkin url")
	flag.Parse()
	defer cancel()

	tracer, err := config.InitTracer(*url)
	if err != nil {
		return
	}
	defer func() {
		if err := tracer.Shutdown(ctx); err != nil {
			log.Fatal("failed to tracer TracerProvider: %w", err)
		}
	}()

	var tr = otel.GetTracerProvider().Tracer("component-main")
	ctx, span := tr.Start(ctx, "foo", trace.WithSpanKind(trace.SpanKindServer))
	<-time.After(6 * time.Millisecond)
	bar(ctx)
	<-time.After(6 * time.Millisecond)
	span.End()*/
	router.Use(otelmux.Middleware("service-name"))
	http.ListenAndServe(":8080", router)
}

func zipkinCodeHandler(writer http.ResponseWriter, request *http.Request) {
	bar(request.Context())
	response, err := utils.SendRequest(request.Context(), http.MethodGet, "http://localhost:8081/api/home", nil)
	if err != nil {
		return
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, _ = writer.Write(bytes)
}

func bar(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	<-time.After(6 * time.Millisecond)
	span.End()
}
