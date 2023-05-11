package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

const (
	AND      = " AND "
	AS       = " AS "
	EQUALS   = " = "
	FROM     = " FROM "
	JOIN     = " JOIN "
	LIMIT    = " LIMIT "
	OFFSET   = " OFFSET "
	ON       = " ON "
	OR       = " OR "
	ORDER_BY = " ORDER BY "
	SELECT   = " SELECT "
	WHERE    = " WHERE "
)

type Query struct {
	Columns       string
	Source        string
	WhereClause   []string
	OrderClause   []string
	JoinStatement string
	Limit         *int
	Offset        *int
	Name          string
}

func NewQuery(source interface{}, columns ...string) *Query {
	if len(columns) == 0 {
		columns = []string{"*"}
	}

	if sourceString, err := CreateSourceString(source); err != nil {
		return &Query{}
	} else {
		return &Query{
			Columns: strings.Join(columns[:], ", "),
			Source:  sourceString,
		}
	}
}

func (q *Query) addToWhereClause(statements ...string) {
	q.WhereClause = append(q.WhereClause, statements...)
}

func (q *Query) All() (*sql.Rows, error) {
	if queryString, err := q.Compile(); err != nil {
		return nil, err
	} else {
		return driver.Query(queryString)
	}
}

func (q *Query) First(dest ...any) error {
	if queryString, err := q.Compile(); err != nil {
		return err
	} else {
		return driver.QueryRow(queryString).Scan(dest...)
	}
}

func (q *Query) Compile() (string, error) {
	queryString := SELECT + q.Columns + FROM + q.Source

	if q.JoinStatement != "" {
		queryString += q.JoinStatement
	}

	if len(q.WhereClause) > 0 {
		queryString += WHERE + strings.Join(q.WhereClause[:], AND)
	}
	if len(q.OrderClause) > 0 {
		queryString += ORDER_BY + strings.Join(q.OrderClause[:], ", ")
	}
	if q.Limit != nil {
		queryString += fmt.Sprintf("%s%d", LIMIT, *q.Limit)
	}
	if q.Offset != nil {
		queryString += fmt.Sprintf("%s%d ", OFFSET, *q.Offset)
	}

	return queryString, nil
}

func (q *Query) Filter(column string, operator string, value interface{}) *Query {
	var valueString string
	switch v := value.(type) {
	case int:
		valueString = strconv.Itoa(v)
	case float64:
		valueString = strconv.FormatFloat(v, 'E', -1, 64)
	case string:
		valueString = v
	case *Date:
		valueString = v.String()
	case bool:
		valueString = strconv.FormatBool(v)
	default:
		return q
	}

	q.addToWhereClause(column + " " + operator + " '" + valueString + "'")
	return q
}

func (q *Query) NameQuery(name string) *Query {
	q.Name = name
	return q
}

func (q *Query) Join(source interface{}, joinType string, leftOn string, rightOn string) *Query {
	switch s := source.(type) {
	case string:
		q.JoinStatement = " " + strings.ToUpper(joinType) + JOIN + s + ON + leftOn + EQUALS + rightOn
	case *Query:
		if sourceString, err := CreateSourceString(s); err != nil {
			return q
		} else {
			q.JoinStatement = " " + strings.ToUpper(joinType) + JOIN + sourceString + ON + leftOn + EQUALS + rightOn
		}
	}

	return q
}

func (q *Query) Where(whereStatement string, extraStatements ...string) *Query {
	q.addToWhereClause(whereStatement)
	q.addToWhereClause(extraStatements...)

	return q
}

func CreateSourceString(source interface{}) (string, error) {
	switch s := source.(type) {
	case string:
		return s, nil
	case *Query:
		if queryString, err := s.Compile(); err != nil {
			return "", err
		} else {
			if s.Name != "" {
				return "(" + queryString + " ) " + s.Name, nil
			}
			return "(" + queryString + " ) " + "subq", nil
		}
	default:
		return "", nil
	}
}

func And(statements ...string) string {
	return strings.Join(statements[:], AND)
}

func Or(statements ...string) string {
	return strings.Join(statements[:], OR)
}