package commandcenter

import "context"

// Factory implements ClientFactory for creating command center clients
type Factory struct{}

// NewFactory creates a new command center client factory
func NewFactory() ClientFactory {
	return &Factory{}
}

// NewClient creates a new command center client for the given site
func (f *Factory) NewClient(siteID int) AccessCodeManager {
	return &Client{
		siteID: siteID,
		ctx:    context.Background(),
	}
}

// NewClientWithContext creates a new command center client with context
func (f *Factory) NewClientWithContext(siteID int, ctx context.Context) AccessCodeManager {
	return &Client{
		siteID: siteID,
		ctx:    ctx,
	}
}

// NewClientWithConfig creates a new command center client with full configuration
func (f *Factory) NewClientWithConfig(config ClientConfig) AccessCodeManager {
	return &Client{
		siteID: config.SiteID,
		ctx:    config.Context,
	}
}
