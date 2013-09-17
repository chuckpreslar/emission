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
    Emit("test")

  if flag {
    t.Error("Emit failed to call listener to unset flag.")
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
    Emit("test").
    Emit("test")

  if flag {
    t.Error("Once called listener multiple times reseting the flag.")
  }
}
