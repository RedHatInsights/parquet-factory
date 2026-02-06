// Copyright 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/parquet-factory/conf"
	"github.com/RedHatInsights/parquet-factory/metrics"
)

func initInfoLog(msg string) {
	log.Info().Str("type", "init").Msg(msg)
}

func printVersionInfo() {
	initInfoLog("Version: " + buildVersion)
	initInfoLog("Build time: " + buildTime)
	initInfoLog("Branch: " + buildBranch)
	initInfoLog("Commit: " + buildCommit)
}

func startMetrics() error {
	ctx := context.TODO()
	sentryConf := conf.GetSentryConfiguration()

	if err := metrics.InitMetrics("ccx-" + sentryConf.SentryEnvironment); err != nil {
		log.Error().Msgf("metrics cannot be loaded: %s", err)
		return err
	}
	metrics.State.Set(metrics.Init)

	metricsConf := conf.GetMetricsConfiguration()
	go push.SendMetricsInLoop(ctx, metricsConf.Job, metricsConf.GatewayURL, metricsConf.GatewayAuthToken, time.Duration(metricsConf.TimeBetweenPush)*time.Second)
	return nil
}
