package emission

import (
	"testing"
)

func TestAddListener(t *testing.T) {
	event := "test"

	emitter := NewEmitter().
		AddListener(event, func() {})

	if 1 != len(emitter.events[event]) {
		t.Error("Failed to add listener to the emitter.")
	}
}

func TestEmit(t *testing.T) {
	event := "test"
	flag := true

	NewEmitter().
		AddListener(event, func() { flag = !flag }).
		Emit(event)

	if flag {
		t.Error("Emit failed to call listener to unset flag.")
	}
}

func TestEmitWithMultipleListeners(t *testing.T) {
	event := "test"
	invoked := 0

	NewEmitter().
		AddListener(event, func() {
			invoked = invoked + 1
		}).
		AddListener(event, func() {
			invoked = invoked + 1
		}).
		Emit(event)

	if invoked != 2{
		t.Error("Emit failed to call all listeners.")
	}
}

func TestRemoveListener(t *testing.T) {
	event := "test"
	listener := func() {}

	emitter := NewEmitter().
		AddListener(event, listener).
		RemoveListener(event, listener)

	if 0 != len(emitter.events[event]) {
		t.Error("Failed to remove listener from the emitter.")
	}
}

func TestOnce(t *testing.T) {
	event := "test"
	flag := true

	NewEmitter().
		Once(event, func() { flag = !flag }).
		Emit(event).
		Emit(event)

	if flag {
		t.Error("Once called listener multiple times reseting the flag.")
	}
}

func TestRecoveryWith(t *testing.T) {
	event := "test"
	flag := true

	NewEmitter().
		AddListener(event, func() { panic(event) }).
		RecoverWith(func(event, listener interface{}, err error) { flag = !flag }).
		Emit(event)

	if flag {
		t.Error("Listener supplied to RecoverWith was not called to unset flag on panic.")
	}
}
