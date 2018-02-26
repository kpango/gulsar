package gulsar

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

// Pulsar is pulsar client struct
type Pulsar struct {
	Origin    string
	URL       string
	Cert      string
	Key       string
	Conn      *websocket.Conn
	mu        sync.Mutex
	connected bool
}

var (
	// ErrConnectionNotFound is not connected error value
	ErrConnectionNotFound = errors.New("Connection Not Found")
)

// New is Pulsar constructor
func New(origin, url string) *Pulsar {
	return &Pulsar{
		Origin:    origin,
		URL:       url,
		connected: false,
	}
}

// Connect connect to pulsar
func (p *Pulsar) Connect(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := websocket.NewConfig(p.URL, p.Origin)
	if err != nil {
		return err
	}

	if strings.HasPrefix(p.URL, "wss") && p.Cert != "" && p.Key != "" {
		cer, err := tls.LoadX509KeyPair(p.Cert, p.Key)
		if err != nil {
			return err
		}

		cfg.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cer},
		}
	}

	p.Conn, err = websocket.DialConfig(cfg)
	if err != nil {
		return err
	}

	p.connected = true

	return nil
}

// ConnectWithHeader connect to pulsar with header
func (p *Pulsar) ConnectWithHeader(ctx context.Context, header http.Header) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := websocket.NewConfig(p.URL, p.Origin)
	if err != nil {
		return err
	}

	cfg.Header = header

	if strings.HasPrefix(p.URL, "wss") && p.Cert != "" && p.Key != "" {
		cer, err := tls.LoadX509KeyPair(p.Cert, p.Key)
		if err != nil {
			return err
		}

		cfg.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cer},
		}
	}

	p.Conn, err = websocket.DialConfig(cfg)
	if err != nil {
		return err
	}

	p.connected = true

	return nil
}

// Consume receives pulsar produced raw message
func (p *Pulsar) Consume() (*ConsumeMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.connected || p.Conn == nil {
		return nil, ErrConnectionNotFound
	}

	var msg *ConsumeMessage
	err := websocket.JSON.Receive(p.Conn, &msg)
	if err != nil {
		return nil, err
	}
	err = msg.decodePayload()
	return msg, err
}

// ConsumeWithJSONDecode receives pulsar produced decoded JSON message
func (p *Pulsar) ConsumeWithJSONDecode(v interface{}) error {
	msg, err := p.Consume()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg.Body, &v)
}

// Produce sends message to pulsar
func (p *Pulsar) Produce(msg *ProduceMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.connected || p.Conn == nil {
		return ErrConnectionNotFound
	}
	return websocket.JSON.Send(p.Conn, msg)
}

// ReceiveACK receives ack
func (p *Pulsar) ReceiveACK() (*ACKPayload, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.connected || p.Conn == nil {
		return nil, ErrConnectionNotFound
	}

	var msg *ACKPayload
	err := websocket.JSON.Receive(p.Conn, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// SendACK sends ack to pulsar
func (p *Pulsar) SendACK(msgID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.connected || p.Conn == nil {
		return ErrConnectionNotFound
	}
	return websocket.JSON.Send(p.Conn, &ACKPayload{
		MessageID: msgID,
	})
}

// SetTlsCredential set credential to PulsarClient
func (p *Pulsar) SetTlsCredential(cert, key string) *Pulsar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Cert = cert
	p.Key = key

	return p
}
