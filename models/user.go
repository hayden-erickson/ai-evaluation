package models

// BUser represents a business user in the system
type BUser struct {
	CompanyUUID string
	Id          int
	Sites       []string
}

// Claims represents JWT claims for authentication
type Claims struct {
	CompanyUUID     string
	CurrentSite     int
	CurrentSiteUUID string
	UserID          int
}
