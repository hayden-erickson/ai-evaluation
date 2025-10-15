package commandcenter

// AccessCodeManager defines the interface for managing access codes
type AccessCodeManager interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

// ClientFactory defines the interface for creating command center clients
type ClientFactory interface {
	NewClient(siteID int) AccessCodeManager
}
