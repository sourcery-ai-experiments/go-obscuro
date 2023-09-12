package p2p

import (
	"fmt"
	"testing"
	"time"
)

func TestValidatorAndSequencerCommunication(t *testing.T) {
	// Starting a Sequencer on a test address
	addr := "localhost:8080"
	seq := NewSequencer(addr)
	go seq.Start()                     // start in a separate goroutine
	time.Sleep(time.Millisecond * 100) // give a short time for server to start

	// Creating a Validator that connects to the sequencer
	validator := NewValidator("ws://" + addr)
	err := validator.Start()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	err = validator.SendMessage([]byte("ping from the validator"))
	if err != nil {
		t.Fatalf("Failed to send ping: %v", err)
	}

	seq.BroadcastHello()

	// A short delay to let asynchronous actions complete
	time.Sleep(time.Millisecond * 100)
}

func TestMultipleValidatorsAndSequencerCommunication(t *testing.T) {
	// Starting a Sequencer on a test address
	addr := "localhost:8081"
	seq := NewSequencer(addr)
	go seq.Start()                     // start in a separate goroutine
	time.Sleep(time.Millisecond * 100) // give a short time for server to start

	validatorCount := 10
	interval := time.Millisecond * 200

	for i := 0; i < validatorCount; i++ {
		go func(index int) {
			// Wait for the interval * index before creating and pinging with the validator
			time.Sleep(interval * time.Duration(index))

			// Creating a Validator that connects to the sequencer
			validator := NewValidator("ws://" + addr)
			err := validator.Start()
			if err != nil {
				t.Errorf("Validator %d: failed to create: %v", index, err)
				return
			}

			err = validator.SendMessage([]byte(fmt.Sprintf("Ping from validator - %d", index)))
			if err != nil {
				t.Errorf("Validator %d: failed to send ping: %v", index, err)
			}

		}(i)
	}

	// Allow enough time for all validators to send their ping messages and get responses
	time.Sleep(interval * time.Duration(validatorCount+1))

	seq.BroadcastHello()

	// A short delay to let asynchronous actions complete
	time.Sleep(time.Millisecond * 100)
}
