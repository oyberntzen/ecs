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

func TestNewEntity(t *testing.T) {
	scene := ecs.Scene{}
	e1 := scene.NewEntity()
	e2 := scene.NewEntity()

	assert.Equal(t, false, e1 == e2, "e1 and e2 should not be equal")
}

func TestAddGetRemoveComponent(t *testing.T) {
	scene := ecs.Scene{}
	e := scene.NewEntity()

	type comp1 struct {
		num int
	}

	_, ok := scene.GetComponent(e, comp1{})
	assert.Equal(t, false, ok, "component should not exist")

	c1 := comp1{45}
	scene.AddComponent(e, c1)
	c2, ok := scene.GetComponent(e, comp1{})
	assert.Equal(t, true, ok, "component should exist")
	assert.Equal(t, c1, c2, "components should be equal")

	scene.RemoveComponent(e, comp1{})
	_, ok = scene.GetComponent(e, comp1{})
	assert.Equal(t, false, ok, "component should not exist")
}

func TestAddComponentMany(t *testing.T) {
	num := 100
	scene := ecs.Scene{}

	type comp1 struct {
		num int
	}

	entities := make([]ecs.Entity, num)
	for i := 0; i < num; i++ {
		entities[i] = scene.NewEntity()
		scene.AddComponent(entities[i], comp1{i})
	}

	for i := 0; i < num; i++ {
		c, ok := scene.GetComponent(entities[i], comp1{})
		assert.Equal(t, true, ok, "component should exist")
		assert.Equal(t, c, comp1{i}, "components should be equal")
	}
}
