package repository

// CommandCenterClient interface for command center operations
type CommandCenterClient interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

// CommandCenter implements CommandCenterClient
type CommandCenter struct{}

// RevokeAccessCodes revokes access codes for given units
func (cc *CommandCenter) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	return nil
}

// SetAccessCodes sets access codes for given units
func (cc *CommandCenter) SetAccessCodes(units []int, options map[string]struct{}) error {
	return nil
}
