package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/gapi"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDbMigrations(config.MigrationsURL, config.DBSource)

	store := db.NewStore(conn)
	// go runGRPCGatewayServer(config, store)
	runGinServer(config, store)
	runGRPCServer(config, store)
}

func runDbMigrations(migrationUrl string, databaseUrl string) {
	migration, err := migrate.New(migrationUrl, databaseUrl)

	if err != nil {
		log.Fatalln("Could not load migrations: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln("Could not run migrations: ", err)
	}

	log.Println("Migration run sucessfully")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {

		log.Fatalln("Could not start server", err)
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatalln("Could not create gRPC server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {

		log.Fatalln("Could not start server", err)
	}

	err = grpcServer.Serve(listener)

	if err != nil {

		log.Fatalln("Could not gRPC server", err)
	}

}

func runGRPCGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatalln("Could not create gRPC server", err)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatalln("Could not register handler server", err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./docs/swagger"))

	mux.Handle("/", grpcMux)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {

		log.Fatalln("Could not create listener", err)
	}

	err = http.Serve(listener, mux)

	if err != nil {

		log.Fatalln("Could not start gRPC gateway", err)
	}
}
