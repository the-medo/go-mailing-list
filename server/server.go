package main

import (
	"database/sql"
	"github.com/alexflint/go-arg"
	"go-mailing-list/grpcapi"
	"go-mailing-list/jsonapi"
	"go-mailing-list/mdb"
	"log"
	"sync"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
	BindGrpc string `arg:"env:MAILINGLIST_BIND_GRPC"`
}

func main() {
	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}

	if args.BindGrpc == "" {
		args.BindGrpc = ":8081"
	}

	log.Printf("using database '%v' \n", args.DbPath)

	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mdb.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting json server on %v\n", args.BindJson)
		jsonapi.Serve(db, args.BindJson)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Printf("starting gRPC server on %v\n", args.BindGrpc)
		grpcapi.Serve(db, args.BindGrpc)
		wg.Done()
	}()

	wg.Wait()
}
