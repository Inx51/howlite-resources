package types

import "time"

var ResourceRemoavedEventType = "ResourceRemoved"

type ResourceRemoved struct {
	RemovedUtc       time.Time
	ResourceIdentity string
}
