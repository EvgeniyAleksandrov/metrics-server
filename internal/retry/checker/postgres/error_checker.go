package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGErrorChecker struct {
	pgErrorCodes map[string]struct{}
}

func NewPGErrorChecker() *PGErrorChecker {
	pgErrorCodes := []string{pgerrcode.ConnectionException}

	classificatorsCodesMap := make(map[string]struct{}, len(pgErrorCodes))

	for _, pgErrorCode := range pgErrorCodes {
		if len(pgErrorCode) < 2 {
			panic("Invalid postgres error class")
		}
		classificatorCode := pgErrorCode[:2]

		classificatorsCodesMap[classificatorCode] = struct{}{}
	}

	return &PGErrorChecker{
		pgErrorCodes: classificatorsCodesMap,
	}
}

func (c *PGErrorChecker) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if len(pgErr.Code) < 2 {
			return false
		}

		if _, ok := c.pgErrorCodes[pgErr.Code[:2]]; ok {
			return true
		}

		return false
	}

	return false
}
