package mvcc

import "errors"

type StateMachineStatus int

const (
	StateMachineStatusIsDumping StateMachineStatus = iota
)

var ErrStateMachineIsDumping = errors.New("state machine is dumping, do not repeat the operation")
