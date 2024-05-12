package common

type Message struct {
	Type   string `json:"type"`
	Key    string `json:"key"`
	Msg    int    `json:"msg"`
	Offset int    `json:"offset,omitempty"`
}
