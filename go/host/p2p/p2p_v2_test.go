package p2p

import (
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
	validator, err := NewValidator("ws://" + addr)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Validator sends a "Ping" message to the Sequencer
	// (The response "Pong" should be logged by the Validator's handleMessages method)
	err = validator.Ping()
	if err != nil {
		t.Fatalf("Failed to send ping: %v", err)
	}

	// Sequencer broadcasts a "Hello" message
	// (The "Hello" message should be logged by the Validator's handleMessages method)
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
			validator, err := NewValidator("ws://" + addr)
			if err != nil {
				t.Errorf("Validator %d: failed to create: %v", index, err)
				return
			}

			// Validator sends a "Ping" message to the Sequencer
			// (The response "Pong" should be logged by the Validator's handleMessages method)
			err = validator.Ping()
			if err != nil {
				t.Errorf("Validator %d: failed to send ping: %v", index, err)
			}

		}(i)
	}

	// Allow enough time for all validators to send their ping messages and get responses
	time.Sleep(interval * time.Duration(validatorCount+1))

	// Sequencer broadcasts a "Hello" message
	// (The "Hello" message should be logged by each Validator's handleMessages method)
	seq.BroadcastHello()

	// A short delay to let asynchronous actions complete
	time.Sleep(time.Millisecond * 100)
}
