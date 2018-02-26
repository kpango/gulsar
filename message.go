package gulsar

import "encoding/base64"

// ConsumeMessage is consumed message struct
type ConsumeMessage struct {
	MessageID   string `json:"messageId"`
	PublishTime string `json:"publishTime,omitempty"`
	Message
}

// ProduceMessage is produced message struct
type ProduceMessage struct {
	ReplicationClusters []string `json:"replicationClusters"`
	Message
}

// Message is base message struct
type Message struct {
	Payload    string                 `json:"payload"`
	Properties map[string]interface{} `json:"properties"`
	Context    string                 `json:"context,omitempty"`
	Body       []byte
}

// ACKPayload ack payload struct
type ACKPayload struct {
	Context   string `json:"context"`
	MessageID string `json:"messageId,omitempty"`
	Result    string `json:"result"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
}

func (m *Message) encodePayload() {
	if m.Payload == "" {
		m.Payload = base64.StdEncoding.EncodeToString(m.Body)
	}
}

func (m *Message) decodePayload() error {
	var err error
	m.Body, err = base64.StdEncoding.DecodeString(m.Payload)
	return err
}
