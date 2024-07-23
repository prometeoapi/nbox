package main

import (
	"context"
	"flag"
	"log"
	"nbox/internal/adapters/aws"
	"nbox/internal/application"
	"nbox/internal/entrypoints/api"
	"nbox/internal/entrypoints/api/handlers"
	"nbox/internal/entrypoints/api/health"
	"nbox/internal/usecases"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func main() {
	var port string
	var address string

	flag.StringVar(&port, "port", "7337", "--port=7337")
	flag.StringVar(&address, "address", "", "--address=0.0.0.0")
	flag.Parse()

	fx.New(
		fx.Provide(aws.NewAwsConfig),
		fx.Provide(aws.NewS3Client),
		fx.Provide(aws.NewDynamodbClient),
		fx.Provide(aws.NewStoreAdapter),
		fx.Provide(aws.NewDynamodbBackend),
		fx.Provide(handlers.NewEntryHandler),
		fx.Provide(handlers.NewBoxHandler),
		fx.Provide(usecases.NewPathUseCase),
		fx.Provide(usecases.NewBox),
		fx.Provide(application.NewConfig),
		fx.Provide(api.NewApi),
		fx.Provide(health.NewHealthy),
		fx.Invoke(func(api *api.Api, config *application.Config) {
			done := make(chan error)
			ctx := context.Background()

			server := &http.Server{
				Addr:              net.JoinHostPort(address, port),
				Handler:           api.Engine,
				ReadHeaderTimeout: 30 * time.Second,
			}

			go func() {
				log.Printf("starting server on %s%s\n", address, port)
				done <- server.ListenAndServe()
			}()

			select {
			case err := <-done:
				log.Fatal(err)
			case <-ctx.Done():
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal(err)
				}
			}
		}),
	)

}
