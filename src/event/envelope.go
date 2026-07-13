package event

import "encoding/json"

type Envelope struct {
	Data json.RawMessage `json:"data"`
	Type string          `json:"type"`
}

func NewEnvelope(eventType string, eventData any) (*Envelope, error) {
	raw, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		Type: eventType,
		Data: raw,
	}, nil
}
