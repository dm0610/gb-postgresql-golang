package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//********************************************************RUN SELECT********************************************************************
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

func runSearch() {
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

	limit := 2

	hints, err := search(ctx, dbpool, "jenkins", limit)
	if err != nil {
		log.Fatal(err)
	}

	for _, hint := range hints {
		fmt.Println(hint.ProjectID, hint.ProjectTitle, hint.InstanceName, hint.ServiceTitle, hint.ServiceAddress)
	}
}

//********************************************************RUN INSERT********************************************************************

type (
	ProjectID int
)
type Project struct {
	Title      string
	OwnerEmail string
}

func runInsert() {
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

	project := Project{
		Title:      "My Old Corp",
		OwnerEmail: "v.petrov@mail.ru",
	}

	id, err := insert(ctx, dbpool, project)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)
}

func insert(ctx context.Context, dbpool *pgxpool.Pool, project Project) (ProjectID, error) {
	const sql = `
insert into projects (title, owner_email) values
	($1, $2)
returning id;
`

	// При insert разумно использовать метод dbpool.Exec,
	// который не требует возврата данных из запроса.
	// В данном случае после вставки строки мы получаем её идентификатор.
	// Идентификатор вставленной строки может быть использован
	// в интерфейсе приложения.

	var id ProjectID
	err := dbpool.QueryRow(ctx, sql,
		// Параметры должны передаваться в том порядке,
		// в котором перечислены столбцы в SQL запросе.
		project.Title,
		project.OwnerEmail,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert project: %w", err)
	}

	return id, nil
}

//********************************************************RUN UPDATE********************************************************************

type TransactionFunc func(context.Context, pgx.Tx) error

// inTx создает транзакцию и передает её для использования в функцию f
// если в функции f происходит ошибка, транзакция откатывается
func inTx(
	ctx context.Context,
	dbpool *pgxpool.Pool,
	f TransactionFunc,
) error {
	transaction, err := dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	err = f(ctx, transaction)
	if err != nil {
		rbErr := transaction.Rollback(ctx)

		if rbErr != nil {
			log.Print(rbErr)
		}

		return err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		rbErr := transaction.Rollback(ctx)

		if rbErr != nil {
			log.Print(rbErr)
		}

		return err
	}

	return nil
}

func update(
	ctx context.Context, dbpool *pgxpool.Pool, project Project,
) error {
	err := inTx(ctx, dbpool, func(ctx context.Context, tx pgx.Tx) error {
		const sql = `update projects set owner_email = $1 where title = $2;`

		_, err := tx.Exec(ctx, sql, project.OwnerEmail, project.Title)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func runUpdate() {
	project := Project{
		Title:      "New Origin",
		OwnerEmail: "dmv.strelnikov@mail.ru",
	}
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
	err = update(ctx, dbpool, project)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("project updated")
}

//*********************************************************RUN MAIN*********************************************************************
func main() {
	runSearch()
	runUpdate()
	runInsert()
}
