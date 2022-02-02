package mypg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	dbpool *pgxpool.Pool
}

func NewPG(dbpool *pgxpool.Pool) *PG {
	return &PG{dbpool}
}

// В рамках каждого слоя желательно работать только со структурами,
// принадлежащими этому слою.
// Это многословно, но код получается явным и легко изменяемым.
type InstanceNameSearch struct {
	ProjectID      int
	ProjectTitle   string
	InstanceName   string
	ServiceTitle   string
	ServiceAddress string
}

func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]InstanceNameSearch, error) {
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
WHERE instances.instance_name like '$1' 
limit $2;
`

	pattern := "%" + prefix + "%"
	rows, err := s.dbpool.Query(ctx, sql, pattern, limit)
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
