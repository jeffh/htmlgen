// Package internal holds helpers shared by hx and ds.
package internal

import (
	"strconv"
	"time"
)

// FormatDuration emits an HTMX/Datastar interval: "Ns" when d is a whole
// number of seconds, otherwise "Nms". time.Duration.String() is not used
// because HTMX parseInterval and Datastar tagToMs misread values such as
// "1m0s".
func FormatDuration(d time.Duration) string {
	if d%time.Second == 0 {
		return strconv.FormatInt(int64(d/time.Second), 10) + "s"
	}
	return strconv.FormatInt(int64(d/time.Millisecond), 10) + "ms"
}
