package emission

import (
  "fmt"
  "testing"
  "time"
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
  count := 0
  pre := len((*emitter.events[e]).listeners)
  fn := func(args ...interface{}) {
    count++
  }

  emitter.Once(e, fn).
    Emit(e, nil).
    Emit(e, nil)

  if count != 1 {
    t.Errorf("Listner was called %d times, expected to be called only once.\n", count)
  } else if post := len((*emitter.events[e]).listeners); pre != post {
    t.Errorf("Expected %d event handler(s), found %d.\n", pre, post)
  }
}

func TestThrottled(t *testing.T) {
  flag := 0
  interval := 1
  duration := time.Duration(interval) * time.Millisecond
  fn := func(args ...interface{}) {
    flag++
  }

  emitter.Throttled(e, duration, fn)

  for i := 0; i < 3; i++ {
    emitter.
      Emit(e).
      Emit(e).
      Emit(e)
    if flag != i {
      t.Logf("Expected only one call to throttled"+
        "function within %d ns, was called %v times.\n", interval, flag)
    }
    time.Sleep(duration * 2)
  }

  fmt.Printf("Throttled function was called %v times in %v.\n", flag, 3*duration)

  if flag != 3 {
    t.Errorf("Expected the throtted function"+
      " to be called a minimum of 3 times, was called %v times.\n", flag)
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
  emitter.SetMaxListeners(DEFAULT_MAX_LISTENERS)
}
