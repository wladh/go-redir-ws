package statmsg

import (
	"time"
)

type Statmsg struct {
	Time *time.Time
	Key string
	IP string
	Referer string
	UA string
}
