// Copyright 2022 Øystein Berntzen

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecs_test

import (
	"testing"

	"github.com/oyberntzen/ecs"
	"github.com/stretchr/testify/assert"
)

type system1 struct {
	ecs.System
	updated bool
}

func (sys *system1) Update(dt float64) {
	sys.updated = true
}

type system2 struct {
	system1
	inited  bool
	deleted bool
}

func (sys *system2) Init() {
	sys.inited = true
}

func (sys *system2) Delete() {
	sys.deleted = true
}

func TestSystemUpdate(t *testing.T) {
	scene := ecs.Scene{}

	sys := &system1{}
	scene.AddSystem(sys)
	scene.Update(0)

	assert.Equal(t, true, sys.updated, "System should be updated")
}

func TestSystemInit(t *testing.T) {
	scene := ecs.Scene{}

	sys1 := &system1{}
	sys2 := &system2{}
	scene.AddSystem(sys1)
	scene.AddSystem(sys2)
	scene.Init()

	assert.Equal(t, true, sys2.inited, "System should be inited")
}

func TestSystemDelete(t *testing.T) {
	scene := ecs.Scene{}

	sys1 := &system1{}
	sys2 := &system2{}
	scene.AddSystem(sys1)
	scene.AddSystem(sys2)
	scene.Delete()

	assert.Equal(t, true, sys2.deleted, "System should be deleted")
}

type system3 struct {
	ecs.System
	allComponents []comp1
}

type comp1 struct {
	ecs.Component
	num int
}

func (sys *system3) Update(dt float64) {
	sys.allComponents = ecs.AllComponents[comp1](sys.Scene())
}

func TestSystemAllComponents(t *testing.T) {
	scene := ecs.Scene{}
	sys := &system3{}
	scene.AddSystem(sys)

	entities := make([]ecs.Entity, 100)
	for i := 0; i < len(entities); i++ {
		entities[i] = scene.NewEntity()
		ecs.AddComponent(&comp1{Component: ecs.NewComponent(entities[i]), num: i})
	}

	scene.Update(0)

	for i := 0; i < len(entities); i++ {
		comp := sys.allComponents[i]
		assert.Equal(t, i, comp.num, "Components should be equal")
	}
}

type system4 struct {
	ecs.System
}

func (sys *system4) Update(dt float64) {
	allComponents := ecs.AllComponents[comp1](sys.Scene())
	for i := 0; i < len(allComponents); i++ {
		allComponents[i].num++
	}
}

func BenchmarkSystemAllComponents(b *testing.B) {
	scene := ecs.Scene{}
	sys := &system4{}
	scene.AddSystem(sys)

	entities := make([]ecs.Entity, 100)
	for i := 0; i < len(entities); i++ {
		entities[i] = scene.NewEntity()
		ecs.AddComponent(&comp1{Component: ecs.NewComponent(entities[i]), num: i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			scene.Update(0)
		}
	}
}
