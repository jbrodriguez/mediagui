package core

import "mediagui/domain"

func (c *Core) GetConfig() *domain.Context {
	return c.ctx
}
