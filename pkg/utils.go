package pkg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// return a lexicographical sortable unique ID
func getULID() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func checkTimeValid(id ulid.ULID) error {
	curTime := time.Now()
	time2 := time.Unix(int64(id.Time()), 0)

	if curTime.Sub(time2).Minutes() > 5 {
		return fmt.Errorf(fmt.Sprintf("time has lapsed, diff=%v", curTime.Sub(time2).Seconds()))
	}
	return nil
}
