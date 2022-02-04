package main

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
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return dbpool
}

/*
dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW05$ go run .
1 Customers Feedback cf_jenkins jenkins jenkins.mycompany.ru
2 New Cloud nc_jenkins jenkins jenkins.mycompany.ru
3 Shared Data sd_jenkins jenkins jenkins.mycompany.ru
*/
