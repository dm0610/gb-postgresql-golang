package main

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

type (
	Phone string
	Email string
)

// В рамках каждого слоя желательно работать только со структурами,
// принадлежащими этому слою.
// Это многословно, но код получается явным и легко изменяемым.
type EmailSearchHint struct {
	Phone Phone
	Email Email
}

func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]EmailSearchHint, error) {
	const sql = `
select
	email,
	phone
from employees
where email like $1
order by email asc
limit $2;
`

	pattern := prefix + "%"
	rows, err := s.dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []EmailSearchHint

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint EmailSearchHint

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Email, &hint.Phone)
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
