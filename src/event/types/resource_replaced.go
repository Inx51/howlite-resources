package types

import "time"

var ResourceRepalcedEventType = "ResourceReplaced"

type ResourceReplaced struct {
	ReplacedUtc      time.Time
	ResourceIdentity string
}
