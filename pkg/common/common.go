package common

import (
	"encoding/json"
)

// Message represents a message sent from Src node to Dest node.
// The body is stored as unparsed JSON so the handler can parse it itself.
type Message struct {
	Src  string          `json:"src,omitempty"`
	Dest string          `json:"dest,omitempty"`
	Body json.RawMessage `json:"body,omitempty"`
}

type MessageBody struct {
	Type   string `json:"type"`
	Key    string `json:"key,omitempty"`
	Msg    int    `json:"msg,omitempty"`
	Offset int    `json:"offset,omitempty"`
}
