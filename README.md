# gulsar [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![release](https://img.shields.io/github/release/kpango/gulsar.svg)](https://github.com/kpango/gulsar/releases/latest) [![CircleCI](https://circleci.com/gh/kpango/gulsar.svg?style=shield)](https://circleci.com/gh/kpango/gulsar) [![codecov](https://codecov.io/gh/kpango/gulsar/branch/master/graph/badge.svg)](https://codecov.io/gh/kpango/gulsar) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a6e544eee7bc49e08a000bb10ba3deed)](https://www.codacy.com/app/i.can.feel.gravity/gulsar?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kpango/gulsar&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/kpango/gulsar)](https://goreportcard.com/report/github.com/kpango/gulsar) [![GoDoc](http://godoc.org/github.com/kpango/gulsar?status.svg)](http://godoc.org/github.com/kpango/gulsar) [![Join the chat at https://gitter.im/kpango/gulsar](https://badges.gitter.im/kpango/gulsar.svg)](https://gitter.im/kpango/gulsar?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Pulsar WebSocket Client Library for Go.

## Requirement
Go 1.8

## Installation
```shell
go get github.com/kpango/gulsar
```

## Example
```go
package main

import (
	"context"
	"time"

	"github.com/kpango/glg"
	"github.com/kpango/gulsar"
)

func main() {
	g := gulsar.New("http://localhost", "wss://pulsar:port/path/topic")
	g.SetTlsCredential("sample.cert", "sample.key")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	g.ConnectWithHeader(ctx, map[string][]string{
		"Auth-Header-Key": {
			"Token",
		},
	})

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
```

## Contribution
1. Fork it ( https://github.com/kpango/gulsar/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author
[kpango](https://github.com/kpango)

## LICENSE
gulsar released under MIT license, refer [LICENSE](https://github.com/kpango/gulsar/blob/master/LICENSE) file.
