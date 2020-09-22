package util

import (
	messages "github.com/galenliu/addon-ipc-messages"
	json "github.com/json-iterator/go"
)

func MarshalMessage(message messages.BaseMessage) []byte {
	b, err := json.MarshalIndent(message, "", "")
	if err != nil {
		return b
	}
	return nil
}
