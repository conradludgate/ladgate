package modules

import "time"

type Module struct {
	Server string

	Refresh time.Duration
}
