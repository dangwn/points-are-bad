package api

import (
	dbDriver "database/sql/driver"
	"encoding/json"
	"time"
)

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