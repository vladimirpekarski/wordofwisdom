package gob

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/vladimirpekarski/wordofwisdom/internal/message"
)

func SendMessage[T message.Message](rw io.ReadWriter, message T) error {
	enc := gob.NewEncoder(rw)
	if err := enc.Encode(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func ReceiveMessage[T message.Message](rw io.ReadWriter, message *T) error {
	dec := gob.NewDecoder(rw)
	if err := dec.Decode(message); err != nil {
		return fmt.Errorf("failed to receive message: %w", err)
	}

	return nil
}
