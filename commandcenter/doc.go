// Package commandcenter provides interfaces and implementations for managing
// access codes through the command center system.
//
// The package is organized into several files:
//
// - interfaces.go: Defines the core interfaces (AccessCodeManager, ClientFactory)
// - types.go: Contains all data types and structures
// - factory.go: Implements the factory pattern for creating clients
// - client.go: Contains the concrete implementation of AccessCodeManager
//
// Usage:
//
//	factory := commandcenter.NewFactory()
//	client := factory.NewClientWithContext(siteID, ctx)
//	err := client.SetAccessCodes(units, options)
//
package commandcenter
