package timestamp

import (
	"math"
	"time"
)
import "strconv"

type Timestamp int64

func (t Timestamp) String() string {
	return strconv.FormatInt(int64(t), 10)
}

func (t Timestamp) ToInt64() int64 {
	return int64(t)
}

func GetTime() Timestamp {
	return Timestamp(time.Now().UnixMilli())
}

// ToTimestamp converts a string to a Timestamp.
// If the string is not a valid number, it returns 0.
func ToTimestamp(timestamp string) Timestamp {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return Timestamp(0)
	}
	return Timestamp(ts)
}

func CheckTimestampClose(t1, t2 Timestamp, number float64) bool {
	return math.Abs(float64(t1-t2)) < number
}

func CheckTimestampClose5000(t1, t2 Timestamp) bool {
	return CheckTimestampClose(t1, t2, 50000)
}
