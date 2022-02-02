package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dm0610/gb-postgresql-golang/hw05/mypg"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()
	pg := mypg.NewPG(dbpool)

	//pg := mypg.PG.

	hints, err := pg.Search(ctx, "jenkins", 5)
	if err != nil {
		log.Fatal(err)
	}

	for _, hint := range hints {
		fmt.Println(hint.ProjectID, hint.ProjectTitle, hint.InstanceName, hint.ServiceTitle, hint.ServiceAddress)
	}
}

func connect(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(ctx, "postgres://techuser:techuser@localhost:5432/projects")
	if err != nil {
		panic(err)
	}

	return dbpool
}
