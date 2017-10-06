// https://gist.github.com/bsphere/8369aca6dde3e7b4392c
package timestamp

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}

// GetBSON Transform timestamp ms to time.Time
func (t *Timestamp) GetBSON() (interface{}, error) {
	if time.Time(*t).IsZero() {
		return nil, nil
	}
	return time.Time(*t), nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var tm int64
	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}
	*t = Timestamp(time.Unix(0, tm*int64(time.Millisecond)))
	return nil
}

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}
