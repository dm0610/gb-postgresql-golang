package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

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

//*********************************************************RUN LOAD GENERATOR*********************************************************************

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			_, err := search(ctx, dbpool, "alex", 5)
			if err != nil {
				log.Fatal(err)
			}

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func runLoadGen() {
	ctx := context.Background()

	url := "postgres://techuser:techuser@localhost:5432/projects"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 100
	cfg.MinConns = 50

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	duration := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Println("start attack")
	res := attack(ctx, duration, threads, dbpool)

	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
}

//*********************************************************RUN MAIN*********************************************************************
func main() {
	//runSearch()
	//runUpdate()
	//runInsert()
	runLoadGen()
	/*
	   projects=> select * from projects;
	    id |       title        |      owner_email
	   ----+--------------------+------------------------
	     1 | Customers Feedback | a.ivanov1@mail.ru
	     2 | New Cloud          | s.atremonov2@mail.ru
	     3 | Shared Data        | n.semenov3@mail.ru
	     7 | My Old Corp        | v.petrov@mail.ru
	     4 | New Origin         | dmv.strelnikov@mail.ru
	   (5 rows)
	*/

	/*
		cfg.MaxConns = 1
		cfg.MinConns = 1
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			duration: 10.288708042s
			threads: 1000
			queries: 30315
			QPS: 3031
		cfg.MaxConns = 10
		cfg.MinConns = 5
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			duration: 10.093315301s
			threads: 1000
			queries: 106583
			QPS: 10658
		cfg.MaxConns = 20
		cfg.MinConns = 10
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			duration: 10.084377254s
			threads: 1000
			queries: 117238
			QPS: 11723
		cfg.MaxConns = 40
		cfg.MinConns = 20
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			duration: 10.077240155s
			threads: 1000
			queries: 118083
			QPS: 11808
		cfg.MaxConns = 80
		cfg.MinConns = 40
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			duration: 10.104022697s
			threads: 1000
			queries: 106220
			QPS: 10622
		cfg.MaxConns = 100
		cfg.MinConns = 50
			dmvstrelnikov@dmvstrelnikov-VirtualBox:~/Documents/GeekBrains/GB-Postgres/gb-postgresql-golang/HW04$ go run .
			start attack
			2022/01/28 18:14:16 failed to query data: failed to connect to `host=localhost user=techuser database=projects`: server error (FATAL: remaining connection slots are reserved for non-replication superuser connections (SQLSTATE 53300))
			exit status 1
	*/

}
