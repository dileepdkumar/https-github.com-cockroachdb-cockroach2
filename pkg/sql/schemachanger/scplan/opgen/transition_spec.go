// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package opgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

// transitionSpec is used to describe a transition. It is used to build a
// transitions. Its fields are not validated.
type transitionSpec struct {
	from       scpb.Status
	to         scpb.Status
	revertible bool
	minPhase   scop.Phase
	emitFns    []interface{}
}

type transitionProperty interface {
	apply(spec *transitionSpec)
}

func to(to scpb.Status, properties ...transitionProperty) transitionSpec {
	ts := transitionSpec{
		to:         to,
		revertible: true,
	}
	for _, p := range properties {
		p.apply(&ts)
	}
	return ts
}

func revertible(b bool) transitionProperty {
	return revertibleProperty(b)
}

func minPhase(p scop.Phase) transitionProperty {
	return phaseProperty(p)
}

func emit(fn interface{}) transitionProperty {
	return emitFnSpec{fn}
}

type phaseProperty scop.Phase

func (p phaseProperty) apply(spec *transitionSpec) {
	spec.minPhase = scop.Phase(p)
}

type revertibleProperty bool

func (r revertibleProperty) apply(spec *transitionSpec) {
	spec.revertible = bool(r)
}

var _ transitionProperty = revertibleProperty(true)

type emitFnSpec struct {
	fn interface{}
}

func (e emitFnSpec) apply(spec *transitionSpec) {
	spec.emitFns = append(spec.emitFns, e.fn)
}
