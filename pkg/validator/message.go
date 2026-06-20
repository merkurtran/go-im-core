package validator

import "errors"

var validMessageType = map[string]bool{
	"text":  true,
	"image": true,
	"file":  true,
}

func ValidateMessage(senderID, receiverID, content, messageType string) error {
	if senderID == "" {
		return errors.New("sender id is required")
	}
	if receiverID == "" {
		return errors.New("receiver id is required")
	}
	if senderID == receiverID {
		return errors.New("sender and receiver cannot be the same")
	}
	if !validMessageType[messageType] {
		return errors.New("invalid message type")
	}
	if messageType == "text" && content == "" {
		return errors.New("text message cannot be empty")
	}
	if len(content) > 4096 {
		return errors.New("message content is too long")
	}
	return nil
}
