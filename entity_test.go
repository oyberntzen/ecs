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

	entity.AddComponent(&comp1{num: 10})
	entity.AddComponent(&comp2{num: 21})

	entity.Remove()

	c1, err := entity.GetComponent(&comp1{})
	assert.NotNil(t, err, "there should be error")
	assert.Equal(t, nil, c1, "component should be nil")

	c2, err := entity.GetComponent(&comp2{})
	assert.NotNil(t, err, "there should be error")
	assert.Equal(t, nil, c2, "component should be nil")
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

	expected1 := comp1{num: 10}
	expected2 := comp2{num: 21}

	err1 := entity.AddComponent(&expected1)
	err2 := entity.AddComponent(&expected2)

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	result1, err1 := entity.GetComponent(&comp1{})
	result2, err2 := entity.GetComponent(&comp2{})

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	assert.Equal(t, expected1.num, result1.(*comp1).num, "Components should be equal")
	assert.Equal(t, expected2.num, result2.(*comp2).num, "Components should be equal")

	err1 = entity.RemoveComponent(&comp1{})
	err2 = entity.RemoveComponent(&comp2{})

	assert.Nil(t, err1, "Error should be nil")
	assert.Nil(t, err2, "Error should be nil")

	_, err1 = entity.GetComponent(&comp1{})
	_, err2 = entity.GetComponent(&comp2{})

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
		err := entities[n].AddComponent(&comp{num: n})
		assert.Nil(t, err, "Error should be nil")
	}

	for n := 0; n < 5; n++ {
		result, err := entities[n].GetComponent(&comp{})
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, n, result.(*comp).num, "Components should be equal")
	}

	for n := 0; n < 5; n++ {
		err := entities[n].RemoveComponent(&comp{})
		assert.Nil(t, err, "Error should be nil")
	}

	for n := 0; n < 5; n++ {
		_, err := entities[n].GetComponent(&comp{})
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
			entities[n].AddComponent(&comp1{num: 1})
			entities[n].AddComponent(&comp2{num: 2})
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
			entities[n].AddComponent(&comp{num: n})
		}
		b.StopTimer()

		for n := 0; n < 100; n++ {
			entities[n].RemoveComponent(&comp{})
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
			entities[n].AddComponent(&comp{num: n})
		}

		b.StartTimer()
		for n := 0; n < 100; n++ {
			entities[n].RemoveComponent(&comp{})
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
		entities[n].AddComponent(&comp{num: n})
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for n := 0; n < 100; n++ {
			entities[n].GetComponent(&comp{})
		}
		b.StopTimer()
	}
}
