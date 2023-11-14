package services

import (
	"mediagui/domain"
	"mediagui/logger"
	"mediagui/services/core"
	"mediagui/services/server"
)

type Orchestrator struct {
	ctx *domain.Context
}

func CreateOrchestrator(ctx *domain.Context) *Orchestrator {
	return &Orchestrator{
		ctx: ctx,
	}
}

func (o *Orchestrator) Run() error {
	logger.Blue("starting mediagui v%s (DataDir: %s) ...", o.ctx.Version, o.ctx.DataDir)

	core := core.Create(o.ctx)
	server := server.Create(o.ctx, core)

	err := server.Start()
	if err != nil {
		return err
	}

	err = core.Start()
	if err != nil {
		return err
	}

	return nil
}
