package services

import (
	sqlDriver "points-areb-bad/api-client/sql-driver"
)

var driver sqlDriver.SqlDriver = *sqlDriver.NewSqlDriver(
	POSTGRES_USER,
	POSTGRES_PASSWORD,
	POSTGRES_DB,
)
