package models

// BUser represents a business user in the system
type BUser struct {
	CompanyUUID string
	Id          int
	Sites       []string
}

// Claims represents the authentication claims for a user
type Claims struct {
	CompanyUUID     string
	CurrentSite     int
	CurrentSiteUUID string
	UserID          int
}
