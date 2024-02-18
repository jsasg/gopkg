package types

import (
	"database/sql/driver"
	"time"
)

type Time time.Time

func (t *Time) Scan(src interface{}) error {
	if tt, ok := src.(time.Time); ok {
		parseTime, err := time.Parse("15:04:05", tt.String())
		if err != nil {
			return err
		}
		*t = Time(parseTime)
	}
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05"), nil
}
