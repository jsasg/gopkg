package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateTime
type DateTime struct {
	time.Time
}

func (d *DateTime) Scan(src any) error {
	t, ok := src.(time.Time)
	if ok {
		*d = DateTime{Time: t}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", src)
}

func (d DateTime) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Format("2006-01-02 15:04:05"))), nil
}
