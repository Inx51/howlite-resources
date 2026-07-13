package types

import "time"

var ResourceCreatedEventType = "ResourceCreated"

type ResourceCreated struct {
	CreatedUtc       time.Time
	ResourceIdentity string
}
