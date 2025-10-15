package constants

// Access code state constants
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

// Lock state constants
const (
	LOCK_STATE_OVERLOCK = "overlock"
	LOCK_STATE_GATELOCK = "gatelock"
	LOCK_STATE_PRELET   = "prelet"
)

// Access code validation messages
const (
	AccessCodeMsgDuplicate = "duplicate_code"
)
