package api

import (
	sqlDriver "github.com/dangwn/points-are-bad-tooling/pabsqldriver"
)

var driver sqlDriver.SqlDriver = func() sqlDriver.SqlDriver {
	d, err := sqlDriver.NewSqlDriverV2(POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, "5432", POSTGRES_DB)
	if err != nil {
		panic("could not start postgres driver due to following error: " + err.Error())
	}
	return *d
}()
