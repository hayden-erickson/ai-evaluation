package main

type CommandCenterClient struct {
	center ICommandCenter
}

func (c *CommandCenterClient) RevokeAccessCodes(units []int, options map[string]struct{}) error {
	return c.center.RevokeAccessCodes(units, options)
}

func (c *CommandCenterClient) SetAccessCodes(units []int, options map[string]struct{}) error {
	return c.center.SetAccessCodes(units, options)
}
