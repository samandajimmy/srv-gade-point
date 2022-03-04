package database

import (
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type DbBun struct {
	*bun.DB
}

func (b *DbBun) QueryThenScan(c echo.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := b.QueryContext(c.Request().Context(), query, args...)

	if err != nil {
		return err
	}

	err = b.ScanRows(c.Request().Context(), rows, dest)

	if err != nil {
		return err
	}

	return nil
}
