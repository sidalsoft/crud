package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sidalsoft/crud/cmd/app"
	"github.com/sidalsoft/crud/pkg/customers"
	"github.com/sidalsoft/crud/pkg/managers"
	"github.com/sidalsoft/crud/pkg/products"
	"github.com/sidalsoft/crud/pkg/salePositions"
	"github.com/sidalsoft/crud/pkg/sales"
	"github.com/sidalsoft/crud/pkg/security"
	"go.uber.org/dig"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	//keyGen.GenCertificate()
	/*err := http.ListenAndServeTLS(
		"go.alif.dev:9999",
		"server.crt",
		"server-private.key",
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("hello crypto"))
		}))
	log.Fatal(err)*/

	//keyGen.KeyGenInit()
	//keyGen.Init()
	//keyGen.HomeworkInit()
	host := "0.0.0.0"
	port := "8000"
	dsn := "postgres://postgres:postgres@localhost:5432/bankdb"

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customers.NewService,
		managers.NewManagersService,
		products.NewProductService,
		salePositions.NewSalePositionsService,
		sales.NewSalesService,
		security.NewAuthService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}
	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			return err
		}
	}
	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})

}
