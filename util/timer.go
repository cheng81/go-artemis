package util

func NewTimer(delay float64, repeat bool) (out *Timer) {
	out = new(Timer)
	out.delay = delay
	out.repeat = repeat
	return
}

type Timer struct {
	delay, acc float64

	repeat, done, stopped bool

	callback func()
}

func (t *Timer) SetCallback(c func()) {
	t.callback = c
}

func (t *Timer) Update(delta float64) {
	if t.done || t.stopped {
		return
	}

	t.acc += delta
	if t.acc >= t.delay {
		t.acc -= t.delay

		if t.repeat {
			t.Reset()
		} else {
			t.done = true
		}

		t.callback()
	}
}

func (t *Timer) Reset() {
	t.stopped = false
	t.done = false
	t.acc = 0
}

func (t *Timer) Done() bool     { return t.done }
func (t *Timer) Running() bool  { return !t.done && !t.stopped && (t.acc < t.delay) }
func (t *Timer) Delay() float64 { return t.delay }

func (t *Timer) Stop() {
	t.stopped = true
}
func (t *Timer) SetDelay(delay float64) { t.delay = delay }

func (t *Timer) PercentageRemaning() float64 {
	if t.done {
		return 100.
	} else if t.stopped {
		return 0.
	} else {
		return 1. - (t.delay-t.acc)/t.delay
	}
}
