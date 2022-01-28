package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

//type (
//	Title string
//	Email string
//)

type InstanceNameSearch struct {
	ProjectID      int
	ProjectTitle   string
	InstanceName   string
	ServiceTitle   string
	ServiceAddress string
}

// search ищет всех сотрудников, email которых начинается с prefix.
// Из функции возвращается список InstanceNameSearch, отсортированный по Email.
// Размер возвращаемого списка ограничен значением limit.
func search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]InstanceNameSearch, error) {
	//const sql = `select title, owner_email from projects where owner_email like $1 order by owner_email asc limit $2;`
	const sql = `
SELECT projects.id AS "Project Id", 
	projects.title AS "Project Title", 
	instances.instance_name AS "Instance Name",
	services.service_title AS "Service Title",
	services.service_address AS "Service Address"
FROM projects 
	 JOIN instances
		 ON projects.id = instances.project_id
	 JOIN services
		 ON services.id = instances.service_id
WHERE instances.instance_name like $1 
limit $2;
`

	pattern := "%" + prefix + "%"
	rows, err := dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []InstanceNameSearch

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint InstanceNameSearch

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.ProjectID, &hint.ProjectTitle, &hint.InstanceName, &hint.ServiceTitle, &hint.ServiceAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		hints = append(hints, hint)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return hints, nil
}

func main() {
	ctx := context.Background()

	url := "postgres://techuser:techuser@localhost:5432/projects"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	limit := 5

	hints, err := search(ctx, dbpool, "jenkins", limit)
	if err != nil {
		log.Fatal(err)
	}

	for _, hint := range hints {
		fmt.Println(hint.ProjectID, hint.ProjectTitle, hint.InstanceName, hint.ServiceTitle, hint.ServiceAddress)
	}
}
