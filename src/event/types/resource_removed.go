package types

import "time"

type ResourceRemoved struct {
	RemovedUtc       time.Time
	ResourceIdentity string
}
