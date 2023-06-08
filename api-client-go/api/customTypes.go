package api

import (
	"bytes"
	dbDriver "database/sql/driver"
	"encoding/json"
	"strconv"
	"time"
)

/*
 * Custom Date Object
 */

type Date time.Time

var _ json.Unmarshaler = &Date{}

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
 * Arrays for unnesting into insert queries
 * General structure taken from pq.Arrays
 * Arrays take form UNNEST(ARRAY[val1, val2, ...])
 */
type UnnestStringArray []string
type UnnestInt32Array []int32
type UnnestInt64Array []int64

const defaultUnnestArray string = "UNNEST(ARRAY[])"
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
	if a == nil {
		return NULL
	}

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

	return defaultUnnestArray
}

func (a UnnestInt32Array) String() string {
	if len(a) == 0 {
		return NULL
	}

	if n := len(a); n > 0 {
		b := []byte("UNNEST(ARRAY[")

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}
		
		return string(append(b, ']', ')'))
	}

	return defaultUnnestArray
}

func (a UnnestInt64Array) String() string {
	if len(a) == 0 {
		return NULL
	}

	if n := len(a); n > 0 {
		b := []byte("UNNEST(ARRAY[")

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',', ' ')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return string(append(b, ']', ')'))
	}

	return defaultUnnestArray
}