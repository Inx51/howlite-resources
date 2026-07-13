package event

import "github.com/vmihailenco/msgpack/v5"

type Envelope struct {
	Data msgpack.RawMessage
	Type string
}

func NewEnvelope(eventType string, eventData any) (*Envelope, error) {
	raw, err := msgpack.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		Type: eventType,
		Data: raw,
	}, nil
}
