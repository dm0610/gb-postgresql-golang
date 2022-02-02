package mypg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dm0610/gb-postgresql-golang/hw05/mypg"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()
	var pg mypg.PG
	hints, err := pg.Search(ctx, dbpool, "jenkins", 5)
	if err != nil {
		log.Fatal(err)
	}

	for _, hint := range hints {
		fmt.Println(hint.ProjectID, hint.ProjectTitle, hint.InstanceName, hint.ServiceTitle, hint.ServiceAddress)
	}
}

func connect(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return dbpool
}
