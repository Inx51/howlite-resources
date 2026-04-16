package types

import "time"

type ResourceCreated struct {
	CreatedUtc       time.Time
	ResourceIdentity string
}
