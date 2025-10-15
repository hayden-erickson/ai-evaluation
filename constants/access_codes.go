package constants

// Access code states
const (
	AccessCodeStateActive      = "active"
	AccessCodeStateSetup       = "setup"
	AccessCodeStatePending     = "pending"
	AccessCodeStateInactive    = "inactive"
	AccessCodeStateRemoved     = "removed"
	AccessCodeStateRemoving    = "removing"
	AccessCodeStateOverlocking = "overlocking"
	AccessCodeStateOverlocked  = "overlocked"
	AccessCodeStateRemove      = "remove"
)

// Lock states
const (
	LockStateOverlock = "overlock"
	LockStateGatelock = "gatelock"
	LockStatePrelet   = "prelet"
)

// Access code validation messages
const (
	AccessCodeMsgDuplicate = "duplicate_code"
)
