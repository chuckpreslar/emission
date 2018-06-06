package emission

import (
	"testing"
)

func TestAddListener(t *testing.T) {
	event := "test"

	emitter := NewEmitter()
	emitter.AddListener(event, func() {})

	if 1 != len(emitter.events[event]) {
		t.Error("Failed to add listener to the emitter.")
	}
}

func TestEmit(t *testing.T) {
	event := "test"
	flag := true

	emitter := NewEmitter()
	emitter.AddListener(event, func() { flag = !flag })
	emitter.Emit(event)

	if flag {
		t.Error("Emit failed to call listener to unset flag.")
	}
}

func TestEmitSync(t *testing.T) {
	event := "test"
	flag := true

	emitter := NewEmitter()
	emitter.AddListener(event, func() { flag = !flag })
	emitter.EmitSync(event)

	if flag {
		t.Error("EmitSync failed to call listener to unset flag.")
	}
}

func TestEmitWithMultipleListeners(t *testing.T) {
	event := "test"
	invoked := 0

	emitter := NewEmitter()
	emitter.AddListener(event, func() {
		invoked = invoked + 1
	})
	emitter.AddListener(event, func() {
		invoked = invoked + 1
	})
	emitter.Emit(event)

	if invoked != 2 {
		t.Error("Emit failed to call all listeners.")
	}
}

func TestRemoveListener(t *testing.T) {
	event := "test"
	listener := func() {}

	emitter := NewEmitter()
	handle := emitter.AddListener(event, listener)
	emitter.RemoveListener(event, handle)

	if 0 != len(emitter.events[event]) {
		t.Error("Failed to remove listener from the emitter.")
	}
}

func TestOnce(t *testing.T) {
	event := "test"
	flag := true

	emitter := NewEmitter()
	emitter.Once(event, func() { flag = !flag })
	emitter.Emit(event)
	emitter.Emit(event)

	if flag {
		t.Error("Once called listener multiple times reseting the flag.")
	}
}

func TestRecoveryWith(t *testing.T) {
	event := "test"
	flag := true

	emitter := NewEmitter()
	emitter.AddListener(event, func() { panic(event) })
	emitter.RecoverWith(func(event, listener interface{}, err error) { flag = !flag })
	emitter.Emit(event)

	if flag {
		t.Error("Listener supplied to RecoverWith was not called to unset flag on panic.")
	}
}

func TestRemoveOnce(t *testing.T) {
	event := "test"
	flag := false
	fn := func() { flag = !flag }

	emitter := NewEmitter()
	handle := emitter.Once(event, fn)
	emitter.RemoveListener(event, handle)
	emitter.Emit(event)

	if flag {
		t.Error("Failed to remove Listener for Once")
	}
}

func TestCountListener(t *testing.T) {
	event := "test"

	emitter := NewEmitter()
	emitter.AddListener(event, func() {})

	if 1 != emitter.GetListenerCount(event) {
		t.Error("Failed to get listener count from emitter.")
	}

	if 0 != emitter.GetListenerCount("fake") {
		t.Error("Failed to get listener count from emitter.")
	}
}

type SomeType struct{}

func (*SomeType) Receiver(evt string) {}

func TestRemoveStructMethod(t *testing.T) {
	event := "test"
	listener := &SomeType{}
	emitter := NewEmitter()
	handle := emitter.AddListener(event, listener.Receiver)

	emitter.RemoveListener(event, handle)
	if 0 != emitter.GetListenerCount(event) {
		t.Error("Failed to remove listener from emitter.")
	}
}

func TestRemoveDoubleListener(t *testing.T) {
	event := "test"

	fn1 := func() {}

	emitter := NewEmitter()
	handle1 := emitter.On(event, fn1)
	handle2 := emitter.On(event, fn1)
	emitter.RemoveListener(event, handle1)
	if 1 != emitter.GetListenerCount(event) {
		t.Error("Should have removed just one listener.")
	}
	emitter.RemoveListener(event, handle2)
	if 0 != emitter.GetListenerCount(event) {
		t.Error("Should have removed both listeners.")
	}
}
