package engine_test

// Values these tests build fixtures from. They live here rather than
// inline because each is used in four or more places, and something
// that means the same thing in every one of them should say so once.
const (
	hostWeb01   = "web-01"
	taskA       = "a"
	taskB       = "b"
	taskChild   = "child"
	taskFailing = "failing"
	taskOK      = "ok"
)
