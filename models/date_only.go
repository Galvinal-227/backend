package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"pkl-attendance/backend/utils"
)

// DateOnly adalah tipe tanggal (tanpa jam) yang selalu diwakili YYYY-MM-DD.
type DateOnly time.Time

func DateOnlyFromTime(t time.Time) DateOnly {
	jt := t.In(utils.JakartaLoc())
	return DateOnly(time.Date(jt.Year(), jt.Month(), jt.Day(), 0, 0, 0, 0, time.UTC))
}

func (d DateOnly) Time() time.Time { return time.Time(d) }
func (d DateOnly) String() string  { return time.Time(d).Format("2006-01-02") }

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d *DateOnly) Scan(v interface{}) error {
	switch val := v.(type) {
	case time.Time:
		*d = DateOnly(val)
		return nil
	case string:
		t, err := time.Parse("2006-01-02", val)
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(val))
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	case nil:
		return nil
	}
	return errors.New("cannot scan into DateOnly")
}

func ParseDateOnly(s string) (DateOnly, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return DateOnly{}, err
	}
	return DateOnly(t), nil
}