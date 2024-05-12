package common

type Message struct {
	Type   string `json:"type"`
	Key    string `json:"key,omitempty"`
	Msg    int    `json:"msg,omitempty"`
	Offset int    `json:"offset,omitempty"`
}
