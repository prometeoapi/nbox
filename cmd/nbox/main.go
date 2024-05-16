package main

import (
	"context"
	"flag"
	"go.uber.org/fx"
	"log"
	"nbox/internal/adapters/aws"
	"nbox/internal/application"
	"nbox/internal/entrypoints/api"
	"net"
	"net/http"
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
		fx.Provide(api.NewEntryHandler),
		fx.Provide(api.NewBoxHandler),
		fx.Provide(application.NewConfig),
		fx.Provide(api.NewApi),
		fx.Invoke(func(api *api.Api, config *application.Config) {
			done := make(chan error)
			ctx := context.Background()

			server := &http.Server{
				Addr:    net.JoinHostPort(address, port),
				Handler: api.Engine,
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
