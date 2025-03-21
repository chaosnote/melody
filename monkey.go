package melody

import (
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Extend the current architecture to include external connection capabilities.

type Monkey struct {
	Melody
}

func (m *Monkey) Dial(setting url.URL, keys map[string]any) error {
	conn, _, e := websocket.DefaultDialer.Dial(setting.String(), nil)

	if e != nil {
		return e
	}

	session := &Session{
		Request:    nil,
		Keys:       keys,
		conn:       conn,
		output:     make(chan envelope, m.Config.MessageBufferSize),
		outputDone: make(chan struct{}),
		melody:     &m.Melody,
		open:       true,
		rwmutex:    &sync.RWMutex{},
	}

	m.hub.register <- session

	m.connectHandler(session)

	go session.writePump()

	session.readPump()

	if !m.hub.closed() {
		m.hub.unregister <- session
	}

	session.close()

	m.disconnectHandler(session)

	return nil
}
