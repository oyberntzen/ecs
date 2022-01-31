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
	entities := [10]ecs.Entity{}
	for i := 0; i < 10; i++ {
		entities[i] = scene.NewEntity()
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i != j {
				assert.Equal(t, false, entities[i] == entities[j], "e1 and e2 should not be equal")
			}
		}
	}
}

func TestEntityRemove(t *testing.T) {
	scene := ecs.Scene{}
	entity := scene.NewEntity()

	type comp1 struct {
		ecs.Component
		num int
	}

	type comp2 struct {
		ecs.Component
		num int
	}

	ecs.AddComponent(&comp1{Component: ecs.NewComponent(entity), num: 10}, &entity)
	ecs.AddComponent(&comp2{Component: ecs.NewComponent(entity), num: 21}, &entity)

	entity.Remove()

	_, err := ecs.GetComponent[comp1](&entity)
	assert.NotNil(t, err, "there should be error")

	_, err = ecs.GetComponent[comp2](&entity)
	assert.NotNil(t, err, "there should be error")
}

func TestEntityAddGetRemoveComponent(t *testing.T) {
	scene := ecs.Scene{}
	entity := scene.NewEntity()

	type comp1 struct {
		ecs.Component
		num int
	}

	type comp2 struct {
		ecs.Component
		num int
	}

	expected1 := comp1{Component: ecs.NewComponent(entity), num: 10}
	expected2 := comp2{Component: ecs.NewComponent(entity), num: 21}

	err1 := ecs.AddComponent(&expected1, &entity)
	err2 := ecs.AddComponent(&expected2, &entity)

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	result1, err1 := ecs.GetComponent[comp1](&entity)
	result2, err2 := ecs.GetComponent[comp2](&entity)

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	assert.Equal(t, expected1.num, result1.num, "Components should be equal")
	assert.Equal(t, expected2.num, result2.num, "Components should be equal")

	err1 = ecs.RemoveComponent[comp1](&entity)
	err2 = ecs.RemoveComponent[comp2](&entity)

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	_, err1 = ecs.GetComponent[comp1](&entity)
	_, err2 = ecs.GetComponent[comp2](&entity)

	assert.NotNil(t, err1, "Error should not be nil")
	assert.NotNil(t, err2, "Error should not be nil")
}

func TestEntityAddGetRemoveComponentMany(t *testing.T) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 5)

	type comp struct {
		ecs.Component
		num int
	}

	for n := 0; n < 5; n++ {
		entities[n] = scene.NewEntity()
		err := ecs.AddComponent(&comp{Component: ecs.NewComponent(entities[n]), num: n}, &entities[n])
		assert.Nil(t, err, "Error should be nil")
	}

	for n := 0; n < 5; n++ {
		result, err := ecs.GetComponent[comp](&entities[n])
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, n, result.num, "Components should be equal")
	}

	for n := 0; n < 5; n++ {
		err := ecs.RemoveComponent[comp](&entities[n])
		assert.Nil(t, err, "Error should be nil")
	}

	for n := 0; n < 5; n++ {
		_, err := ecs.GetComponent[comp](&entities[n])
		assert.NotNil(t, err, "Error should not be nil")
	}
}

func BenchmarkNewEntity(b *testing.B) {
	scene := ecs.Scene{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			scene.NewEntity()
		}
	}
}

func BenchmarkEntityRemove(b *testing.B) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 100)
	for n := 0; n < 100; n++ {
		entities[n] = scene.NewEntity()
	}

	type comp1 struct {
		ecs.Component
		num int
	}

	type comp2 struct {
		ecs.Component
		num int
	}

	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&comp1{Component: ecs.NewComponent(entities[n]), num: 1}, &entities[n])
			ecs.AddComponent(&comp2{Component: ecs.NewComponent(entities[n]), num: 2}, &entities[n])
		}
		b.StartTimer()
		for n := 0; n < 100; n++ {
			entities[n].Remove()
		}
		b.StopTimer()
	}
}

func BenchmarkEntityAddComponent(b *testing.B) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 100)
	for n := 0; n < 100; n++ {
		entities[n] = scene.NewEntity()
	}

	type comp struct {
		ecs.Component
		num int
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&comp{Component: ecs.NewComponent(entities[n]), num: n}, &entities[n])
		}
		b.StopTimer()

		for n := 0; n < 100; n++ {
			ecs.RemoveComponent[comp](&entities[n])
		}
	}
}

func BenchmarkEntityRemoveComponent(b *testing.B) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 100)
	for n := 0; n < 100; n++ {
		entities[n] = scene.NewEntity()
	}

	type comp struct {
		ecs.Component
		num int
	}

	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&comp{Component: ecs.NewComponent(entities[n]), num: n}, &entities[n])
		}

		b.StartTimer()
		for n := 0; n < 100; n++ {
			ecs.RemoveComponent[comp](&entities[n])
		}
		b.StopTimer()
	}
}

func BenchmarkEntityGetComponent(b *testing.B) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 100)

	type comp struct {
		ecs.Component
		num int
	}

	for n := 0; n < 100; n++ {
		entities[n] = scene.NewEntity()
		ecs.AddComponent(&comp{Component: ecs.NewComponent(entities[n]), num: n}, &entities[n])
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for n := 0; n < 100; n++ {
			ecs.GetComponent[comp](&entities[n])
		}
		b.StopTimer()
	}
}
