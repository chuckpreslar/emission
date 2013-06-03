package emission

import (
  "testing"
)

var emitter *Emitter = NewEmitter()
var e string = "test"

func TestAddListener(t *testing.T) {
  emitter.AddListener(e, func(args ...interface{}) {
    return
  })

  event, ok := emitter.events[e]

  if !ok {
    t.Errorf("Expected AddListener to establish event %s.\n", e)
  } else if 1 != len((*event).listeners) {
    t.Errorf("Expected AddListener to add only 1 handler to registered event %s.\n", e)
  }
}

func TestEmit(t *testing.T) {
  flag := true
  fn := func(args ...interface{}) {
    flag = false
  }

  emitter.AddListener(e, fn).
    Emit(e, nil)

  if flag {
    t.Errorf("Emit failed to call listener function.")
  }
}

func TestRemoveListener(t *testing.T) {
  flag := false
  pre := len((*emitter.events[e]).listeners)
  fn := func(args ...interface{}) {
    flag = true
  }

  emitter.AddListener(e, fn).
    RemoveListener(e, fn).
    Emit(e, nil)

  if flag {
    t.Errorf("Unremoved listener modified flag variable, set to %v.\n", flag)
  } else if post := len((*emitter.events[e]).listeners); pre != post {
    t.Errorf("Expected %d event handler(s), found %d.\n", pre, post)
  }
}

func TestOnce(t *testing.T) {
  flag := true
  pre := len((*emitter.events[e]).listeners)
  fn := func(args ...interface{}) {
    flag = !flag
  }

  emitter.Once(e, fn).
    Emit(e, nil).
    Emit(e, nil)

  if flag {
    t.Errorf("Listner was called twice, reset to %v.\n", flag)
  } else if post := len((*emitter.events[e]).listeners); pre != post {
    t.Errorf("Expected %d event handler(s), found %d.\n", pre, post)
  }
}

func TestSetMaxListeners(t *testing.T) {
  flag := false
  pre := len((*emitter.events[e]).listeners)
  fn := func(args ...interface{}) {
    flag = true
  }

  emitter.SetMaxListeners(0).
    AddListener(e, fn).
    Emit(e, nil)

  if flag {
    t.Errorf("Listner was successfully added after lowering maxListeners.\n")
  } else if post := len((*emitter.events[e]).listeners); pre != post {
    t.Errorf("Expected %d event handler(s), found %d.\n", pre, post)
  }
}
