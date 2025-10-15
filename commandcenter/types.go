package commandcenter

import "context"

// Client represents a command center client for managing access codes
type Client struct {
	siteID int
	ctx    context.Context
}

// ClientConfig holds configuration for command center clients
type ClientConfig struct {
	SiteID  int
	Context context.Context
	// Add other configuration fields as needed
	BaseURL string
	APIKey  string
}

// AccessCodeOptions represents options for access code operations
type AccessCodeOptions map[string]struct{}

// AccessCodeRequest represents a request to modify access codes
type AccessCodeRequest struct {
	Units   []int
	Options AccessCodeOptions
}

// AccessCodeResponse represents the response from access code operations
type AccessCodeResponse struct {
	Success bool
	Message string
	Errors  []string
}
