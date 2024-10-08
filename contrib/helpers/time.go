package helpers

import (
	"fmt"
	"time"
)

func ISOWeek(t time.Time) string {
	y, w := t.ISOWeek()
	return fmt.Sprintf("%d-%.2d", y, w)
}
