package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type GeneratorPostgres struct {
	db *sqlx.DB
}

func NewGeneratorPostgres(db *sqlx.DB) *GeneratorPostgres {
	return &GeneratorPostgres{
		db: db,
	}
}

func (p *GeneratorPostgres) CountUsedCodes() (int, error) {
	var availableKeysCount int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE is_used", codeWordsTable)
	row := p.db.QueryRow(query)
	if err := row.Scan(&availableKeysCount); err != nil {
		return 0, err
	}
	return availableKeysCount, nil
}

func (p *GeneratorPostgres) SaveToDB(code, passwordHash string) error {
	query := fmt.Sprintf("INSERT INTO %s (word, password, is_used) VALUES ($1, $2, $3) ON CONFLICT (word) DO NOTHING " , codeWordsTable)
	_, err := p.db.Exec(query, code, passwordHash, false)
	if err != nil {
		return err
	}

	return nil
}
