package types

import "time"

type ResourceReplaced struct {
	ReplacedUtc      time.Time
	ResourceIdentity string
}
