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
	"github.com/smyrman/subx"
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

	t.Run("Expected correct result", subx.Test(subx.Value(sys.updated), subx.CompareEqual(true)))
}

func TestSystemInit(t *testing.T) {
	scene := ecs.Scene{}

	sys1 := &system1{}
	sys2 := &system2{}
	scene.AddSystem(sys1)
	scene.AddSystem(sys2)
	scene.Init()

	t.Run("Expected correct result", subx.Test(subx.Value(sys2.inited), subx.CompareEqual(true)))

}

func TestSystemDelete(t *testing.T) {
	scene := ecs.Scene{}

	sys1 := &system1{}
	sys2 := &system2{}
	scene.AddSystem(sys1)
	scene.AddSystem(sys2)
	scene.Delete()

	t.Run("Expected correct result", subx.Test(subx.Value(sys2.deleted), subx.CompareEqual(true)))

}
