package engine

import (
	"context"

	"github.com/bitrainforest/PandaAgent/inside/checker"
	"github.com/bitrainforest/PandaAgent/inside/config"
	"github.com/bitrainforest/PandaAgent/inside/deal"
	"github.com/bitrainforest/PandaAgent/inside/downloader"
	"github.com/bitrainforest/PandaAgent/inside/types"
	"github.com/rs/zerolog/log"
)

type Engine struct {
	DealTransformer *deal.DealTransform
	Transformer     *downloader.Transformer
	Checker         checker.Checker
	Buf             chan types.Sector
	ctx             context.Context
	cancle          context.CancelFunc
}

func InitEngine(conf config.Config, ctx context.Context) Engine {
	var engine Engine
	engine.Transformer = downloader.InitTransformer(conf, ctx)
	engine.Checker = checker.InitChecker(conf, ctx)
	engine.Buf = make(chan types.Sector, 1024)
	engine.DealTransformer = deal.InitDealTransform(conf, ctx)
	engine.ctx, engine.cancle = context.WithCancel(ctx)
	return engine
}

func (eg Engine) Run() error {
	log.Info().Msgf("[Engine] Engine Start.")
	eg.Checker.Ping()
	eg.Checker.Check(eg.Buf)
	eg.Transformer.Run(eg.Buf)
	eg.DealTransformer.Run()
	return nil
}

func (eg Engine) Stop() {
	log.Info().Msgf("[Engine] Engine Stop.")
	eg.cancle()
}
