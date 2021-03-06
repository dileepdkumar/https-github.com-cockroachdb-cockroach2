// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package kvserver

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/gogo/protobuf/proto"
)

const (
	unavailableRangesRuleName             = "UnavailableRanges"
	underreplicatedRangesRuleName         = "UnderreplicatedRanges"
	requestsStuckInRaftRuleName           = "RequestsStuckInRaft"
	highOpenFDCountRuleName               = "HighOpenFDCount"
	nodeCapacityRuleName                  = "node:capacity"
	clusterCapacityRuleName               = "cluster:capacity"
	nodeCapacityAvailableRuleName         = "node:capacity_available"
	clusterCapacityAvailableRuleName      = "cluster:capacity_available"
	capacityAvailableRatioRuleName        = "capacity_available:ratio"
	nodeCapacityAvailableRatioRuleName    = "node:capacity_available:ratio"
	clusterCapacityAvailableRatioRuleName = "cluster:capacity_available:ratio"
)

// CreateAndAddRules initializes all KV metric rules and adds them
// to the rule registry for tracking. All rules are exported in the
// YAML format.
func CreateAndAddRules(ctx context.Context, ruleRegistry *metric.RuleRegistry) {
	createAndRegisterUnavailableRangesRule(ctx, ruleRegistry)
	createAndRegisterUnderReplicatedRangesRule(ctx, ruleRegistry)
	createAndRegisterRequestsStuckInRaftRule(ctx, ruleRegistry)
	createAndRegisterHighOpenFDCountRule(ctx, ruleRegistry)
	createAndRegisterNodeCapacityRule(ctx, ruleRegistry)
	createAndRegisterClusterCapacityRule(ctx, ruleRegistry)
	createAndRegisterNodeCapacityAvailableRule(ctx, ruleRegistry)
	createAndRegisterClusterCapacityAvailableRule(ctx, ruleRegistry)
	createAndRegisterCapacityAvailableRatioRule(ctx, ruleRegistry)
	createAndRegisterNodeCapacityAvailableRatioRule(ctx, ruleRegistry)
	createAndRegisterClusterCapacityAvailableRatioRule(ctx, ruleRegistry)
}

func createAndRegisterUnavailableRangesRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "(sum by(instance, cluster) (ranges_unavailable)) > 0"
	var annotations []metric.LabelPair
	annotations = append(annotations, metric.LabelPair{
		Name:  proto.String("summary"),
		Value: proto.String("Instance {{ $labels.instance }} has {{ $value }} unavailable ranges"),
	})
	recommendedHoldDuration := 10 * time.Minute
	help := "This check detects when the number of ranges with less than quorum replicas live are non-zero for too long"

	unavailableRanges, err := metric.NewAlertingRule(
		unavailableRangesRuleName,
		expr,
		annotations,
		nil,
		recommendedHoldDuration,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, unavailableRangesRuleName, unavailableRanges, ruleRegistry)
}

func createAndRegisterUnderReplicatedRangesRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "(sum by(instance, cluster) (ranges_underreplicated)) > 0"
	var annotations []metric.LabelPair
	annotations = append(annotations, metric.LabelPair{
		Name:  proto.String("summary"),
		Value: proto.String("Instance {{ $labels.instance }} has {{ $value }} under-replicated ranges"),
	})
	recommendedHoldDuration := time.Hour
	help := "This check detects when the number of ranges with less than desired replicas live is non-zero for too long."

	underreplicatedRanges, err := metric.NewAlertingRule(
		underreplicatedRangesRuleName,
		expr,
		annotations,
		nil,
		recommendedHoldDuration,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, underreplicatedRangesRuleName, underreplicatedRanges, ruleRegistry)
}

func createAndRegisterRequestsStuckInRaftRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "requests_slow_raft > 0"
	var annotations []metric.LabelPair
	annotations = append(annotations, metric.LabelPair{
		Name:  proto.String("summary"),
		Value: proto.String("{{ $value }} requests stuck in raft on {{ $labels.instance }}"),
	})
	recommendedHoldDuration := 10 * time.Minute
	help := "This check detects when requests are taking a very long time in replication."

	requestsStuckInRaft, err := metric.NewAlertingRule(
		requestsStuckInRaftRuleName,
		expr,
		annotations,
		nil,
		recommendedHoldDuration,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, requestsStuckInRaftRuleName, requestsStuckInRaft, ruleRegistry)
}

func createAndRegisterHighOpenFDCountRule(ctx context.Context, ruleRegistry *metric.RuleRegistry) {
	expr := "sys_fd_open / sys_fd_softlimit > 0.8"
	var annotations []metric.LabelPair
	annotations = append(annotations, metric.LabelPair{
		Name:  proto.String("summary"),
		Value: proto.String("Too many open file descriptors on {{ $labels.instance }}: {{ $value }} fraction used"),
	})
	recommendedHoldDuration := 10 * time.Minute
	help := "This check detects when a cluster is getting close to the open file descriptor limit"

	highOpenFDCount, err := metric.NewAlertingRule(
		highOpenFDCountRuleName,
		expr,
		annotations,
		nil,
		recommendedHoldDuration,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, highOpenFDCountRuleName, highOpenFDCount, ruleRegistry)
}

func createAndRegisterNodeCapacityRule(ctx context.Context, ruleRegistry *metric.RuleRegistry) {
	expr := "sum without(store) (capacity)"
	help := "Aggregation expression to compute node capacity."
	nodeCapacity, err := metric.NewAggregationRule(
		nodeCapacityRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, nodeCapacityRuleName, nodeCapacity, ruleRegistry)
}

func createAndRegisterClusterCapacityRule(ctx context.Context, ruleRegistry *metric.RuleRegistry) {
	expr := "sum without(instance) (node:capacity)"
	help := "Aggregation expression to compute cluster capacity."

	clusterCapacity, err := metric.NewAggregationRule(
		clusterCapacityRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, clusterCapacityRuleName, clusterCapacity, ruleRegistry)
}

func createAndRegisterNodeCapacityAvailableRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "sum without(store) (capacity_available)"
	help := "Aggregation expression to compute available capacity for a node."

	var err error
	nodeCapacityAvailable, err := metric.NewAggregationRule(
		nodeCapacityAvailableRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, nodeCapacityAvailableRuleName, nodeCapacityAvailable, ruleRegistry)
}

func createAndRegisterClusterCapacityAvailableRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "sum without(instance) (node:capacity_available)"
	help := "Aggregation expression to compute available capacity for a cluster."

	clusterCapacityAvailable, err := metric.NewAggregationRule(
		clusterCapacityAvailableRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, clusterCapacityAvailableRuleName, clusterCapacityAvailable, ruleRegistry)
}

func createAndRegisterCapacityAvailableRatioRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "capacity_available / capacity"
	help := "Aggregation expression to compute available capacity ratio."

	capacityAvailableRatio, err := metric.NewAggregationRule(
		capacityAvailableRatioRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, capacityAvailableRatioRuleName, capacityAvailableRatio, ruleRegistry)
}

func createAndRegisterNodeCapacityAvailableRatioRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "node:capacity_available / node:capacity"
	help := "Aggregation expression to compute available capacity ratio for a node."

	nodeCapacityAvailableRatio, err := metric.NewAggregationRule(
		nodeCapacityAvailableRatioRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, nodeCapacityAvailableRatioRuleName, nodeCapacityAvailableRatio, ruleRegistry)
}

func createAndRegisterClusterCapacityAvailableRatioRule(
	ctx context.Context, ruleRegistry *metric.RuleRegistry,
) {
	expr := "cluster:capacity_available/cluster:capacity"
	help := "Aggregation expression to compute available capacity ratio for a cluster."

	clusterCapacityAvailableRatio, err := metric.NewAggregationRule(
		clusterCapacityAvailableRatioRuleName,
		expr,
		nil,
		help,
		true,
	)
	maybeAddRuleToRegistry(ctx, err, clusterCapacityAvailableRatioRuleName, clusterCapacityAvailableRatio, ruleRegistry)
}

func maybeAddRuleToRegistry(
	ctx context.Context, err error, name string, rule metric.Rule, ruleRegistry *metric.RuleRegistry,
) {
	if err != nil {
		log.Warningf(ctx, "unable to create kv rule %s: %s", name, err.Error())
	}
	if ruleRegistry == nil {
		log.Warningf(ctx, "unable to add kv rule %s: rule registry uninitialized", name)
	}
	ruleRegistry.AddRule(rule)
}
