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

func TestRemoveListener(t *testing.T) {
  flag := false
  pre := len((*emitter.events[e]).listeners)
  fn := func(args ...interface{}) {
    flag = true
  }

  emitter.AddListener(e, fn).
    RemoveListener(fn).
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
