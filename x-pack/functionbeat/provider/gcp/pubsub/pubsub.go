// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package pubsub

import (
	"context"
	"fmt"

	gpubsub "cloud.google.com/go/pubsub"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
	"github.com/elastic/beats/x-pack/functionbeat/config"
	"github.com/elastic/beats/x-pack/functionbeat/function/beater"
	prov "github.com/elastic/beats/x-pack/functionbeat/provider/gcp/gcp"
	_ "github.com/elastic/beats/x-pack/functionbeat/provider/gcp/include"
)

func RunPubSub(ctx context.Context, m gpubsub.Message) error {
	cfgwarn.Beta("Google Cloud Platform support is in beta")
	settings := instance.Settings{
		Name:            "functionbeat",
		IndexPrefix:     "functionbeat",
		ConfigOverrides: config.FunctionOverrides,
	}

	cfgfile.ChangeDefaultCfgfileFlag(settings.Name)

	return instance.Run(settings, initFunctionbeat(ctx, m))
}

func initFunctionbeat(ctx context.Context, m gpubsub.Message) func(*beat.Beat, *common.Config) (beat.Beater, error) {
	return func(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
		bt, err := beater.New(b, cfg)
		if err != nil {
			return nil, err
		}

		fnbeat, ok := bt.(*beater.Functionbeat)
		if !ok {
			return nil, fmt.Errorf("not Functionbeat")
		}

		fnbeat.Ctx = prov.NewPubSubContext(fnbeat.Ctx, ctx, m)

		return fnbeat, nil
	}
}