package api

import (
	"bytes"
	dbDriver "database/sql/driver"
	"strconv"
	"time"
)

/*
 * Custom Date Object
 * The dates coming in to the API have the form YYYY-MM-DD, and are stored likewise in the DB
 * There is scope for adding times to the object
 */

type Date time.Time

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(bs []byte) error {
	if len(bs) == 2 {
		*d = Date{}
		return nil
	}

	newDate, err := time.Parse(`"`+dateFormat+`"`, string(bs))
	if err != nil {
		return err
	}
	
    *d = Date(newDate)
    return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateFormat)+2)
	b = append(b, '"')
	b = time.Time(*d).AppendFormat(b, dateFormat)
	b = append(b, '"')
	return b, nil
}

func (d Date) Value() (dbDriver.Value, error) {
    if d.String() == "0001-01-01" {
        return nil, nil
    }
    return []byte(time.Time(d).Format(dateFormat)), nil
}

func (d *Date) Scan(v interface{}) error {
    dDate, err := time.Parse("2006-01-02 00:00:00 +0000 +0000", v.(time.Time).String())
	if err != nil {
		return err
	}
    *d = Date(dDate)
    return nil
}

func (d Date) String() string {
    return time.Time(d).Format(dateFormat)
}

/*
 * Custom Prediction Array
 * Used for combining multiple predictions to a user in one query
 * See prediction.go -> updateUserPredictionsByUserId
 * See prediction.go for "Prediction" struct
 */
type PredictionArray []Prediction

/*
 * Appends an int pointer to the prediction array buffer
 * Nil pointers are equivalent to NULL values in the DB
 */
func appendIntPointerToSqlBuffer(b []byte, i *int) []byte {
	if i == nil {
		return append(b, 'N', 'U', 'L', 'L')
	}
	return strconv.AppendInt(b, int64(*i), 10)
}

// Appends a Prediction object to the array buffer
func appendPredictionToArrayBuffer(b []byte, pred Prediction) []byte {
	b = append(b, '(')
	b = append(b, []byte(pred.PredictionId)...)
	b = append(b, ',')
	b = appendIntPointerToSqlBuffer(b, pred.HomeGoals)
	b = append(b, ',')
	b = appendIntPointerToSqlBuffer(b, pred.AwayGoals)
	return append(b, ')')
}

/* 
 * Creates a string of the prediction array
 * String value takes form "(id, h, a),(id, h, a),..."
 * If there are no predictions in the array, the string is empty
 */
func (p PredictionArray) String() string {
	if n := len(p); n > 0 {
		b := make([]byte, 0, 8*n-2)
		b = appendPredictionToArrayBuffer(b, p[0])
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendPredictionToArrayBuffer(b, p[i])
		}		

		return string(b)
	}
	return ""
}

/*
 * Arrays for unnesting into insert queries
 * General structure taken from pq.Arrays
 * Arrays take form UNNEST(ARRAY[val1, val2, ...]) or NULL if their length is 0
 */
type UnnestStringArray []string
type UnnestInt32Array []int32
type UnnestInt64Array []int64
type UnnestDateArray []Date

const NULL string = "NULL"

func UnnestArray(a interface{}) interface {
	String() string
} {
	switch a := a.(type) {
	case []string:
		return (*UnnestStringArray)(&a)
	case []int32:
		return (*UnnestInt32Array)(&a)
	case []int64:
		return (*UnnestInt64Array)(&a)
	case []Date:
		return (*UnnestDateArray)(&a)
	}
	return nil
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '\'')
	for {
		i := bytes.IndexAny(v, `"\\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '\'')
}

func (a UnnestStringArray) String() string {
	if n := len(a); n > 0 {
		b := make([]byte, 0, 14+3*n)
		b = append(b, []byte("UNNEST(ARRAY[")...)
		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}

		return string(append(b, ']', ')'))
	}

	return NULL
}

func (a UnnestInt32Array) String() string {
	if n := len(a); n > 0 {
		b := make([]byte, 0, 14+3*n)
		b := []byte("UNNEST(ARRAY[")

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}
		
		return string(append(b, ']', ')'))
	}

	return NULL
}

func (a UnnestInt64Array) String() string {
	if n := len(a); n > 0 {
		b := make([]byte, 0, 14+3*n)
		b = []byte("UNNEST(ARRAY[")

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',', ' ')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return string(append(b, ']', ')'))
	}

	return NULL
}

func (a UnnestDateArray) String() string {
	if n := len(a); n > 0 {
		b := make([]byte, 0, 14+3*n)
		b = append(b, []byte("UNNEST(ARRAY[")...)
		b = appendArrayQuotedBytes(b, []byte(a[0].String()))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i].String()))
		}

		return string(append(b, ']', ')'))
	}

	return NULL
}
