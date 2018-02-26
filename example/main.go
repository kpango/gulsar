package main

import (
	"context"
	"time"

	"github.com/kpango/glg"
	"github.com/kpango/gulsar"
)

func main() {
	g := gulsar.New("http://localhost", "wss://pulsar:port/path/topic").
		SetTlsCredential("sample.cert", "sample.key")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err := g.ConnectWithHeader(ctx, map[string][]string{
		"Auth-Header-Key": {
			"Token",
		},
	})

	if err != nil {
		glg.Fatalln(err)
	}

	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := g.Produce(&gulsar.ProduceMessage{})
				if err != nil {
					glg.Error(err)
					continue
				}

				ackMsg, err := g.ReceiveACK()
				if err != nil {
					glg.Error(err)
					continue
				}

				glg.Info(ackMsg)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := g.Consume()
			if err != nil {
				glg.Error(err)
				continue
			}

			glg.Info(string(msg.Body))

			g.SendACK(msg.MessageID)

		}
	}
}
