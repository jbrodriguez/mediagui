package cmd

import (
	"mediagui/domain"
	"mediagui/services"
)

type Boot struct {
}

func (b *Boot) Run(ctx *domain.Context) error {
	return services.CreateOrchestrator(ctx).Run()
}
