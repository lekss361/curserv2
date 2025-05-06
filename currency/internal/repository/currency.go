package repository

import (
	"database/sql"
	"github.com/lekss361/curserv2/currency/internal/dto"
	"time"
)

type sqlRatesRepo struct {
	db *sql.DB
}

func NewRatesRepo(db *sql.DB) dto.RatesRepo {
	return &sqlRatesRepo{db: db}
}

func (r *sqlRatesRepo) Save(date time.Time, rates map[string]float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
        INSERT INTO currency_rates (date, code, rate)
        VALUES ($1, $2, $3)
        ON CONFLICT (date, code) DO UPDATE SET rate = EXCLUDED.rate
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for code, rate := range rates {
		prefixedCode := "rub" + code
		if _, err := stmt.Exec(date, prefixedCode, rate); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *sqlRatesRepo) Get(date time.Time) (map[string]float64, error) {
	rows, err := r.db.Query(`
        SELECT code, rate
        FROM currency_rates
        WHERE date = $1
    `, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var code string
		var rate float64
		if err := rows.Scan(&code, &rate); err != nil {
			return nil, err
		}
		result[code] = rate
	}
	return result, rows.Err()
}
