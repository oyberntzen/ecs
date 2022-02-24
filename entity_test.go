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

func TestNewEntity(t *testing.T) {
	scene := ecs.Scene{}
	entities := [10]ecs.Entity{}
	for i := 0; i < 10; i++ {
		entities[i] = scene.NewEntity()
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i != j {
				t.Run("Expected correct result", subx.Test(subx.Value(entities[i]), subx.CompareNotEqual(entities[j])))
				// assert.Equal(t, false, entities[i] == entities[j], "e1 and e2 should not be equal")
			}
		}
	}
}

func TestEntityRemove(t *testing.T) {
	scene := ecs.Scene{}
	entity := scene.NewEntity()

	type comp1 struct {
		num int
	}

	type comp2 struct {
		num int
	}

	ecs.AddComponent(&entity, &comp1{num: 10})
	ecs.AddComponent(&entity, &comp2{num: 21})

	entity.Remove()

	_, err := ecs.GetComponent[comp1](&entity)
	t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareNotEqual[error](nil)))
	// assert.NotNil(t, err, "there should be error")

	_, err = ecs.GetComponent[comp2](&entity)
	// assert.NotNil(t, err, "there should be error")
	t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareNotEqual[error](nil)))
}

func TestEntityAddGetRemoveComponent(t *testing.T) {
	scene := ecs.Scene{}
	entity := scene.NewEntity()

	type comp1 struct {
		num int
	}

	type comp2 struct {
		num int
	}

	expected1 := comp1{num: 10}
	expected2 := comp2{num: 21}

	err1 := ecs.AddComponent(&entity, &expected1)
	err2 := ecs.AddComponent(&entity, &expected2)

	t.Run("Expected correct result", subx.Test(subx.Value(err1), subx.CompareEqual[error](nil)))
	t.Run("Expected correct result", subx.Test(subx.Value(err2), subx.CompareEqual[error](nil)))

	result1, err1 := ecs.GetComponent[comp1](&entity)
	result2, err2 := ecs.GetComponent[comp2](&entity)

	t.Run("Expected correct result", subx.Test(subx.Value(err1), subx.CompareEqual[error](nil)))
	t.Run("Expected correct result", subx.Test(subx.Value(err2), subx.CompareEqual[error](nil)))

	t.Run("Expected correct result", subx.Test(subx.Value(result1.num), subx.CompareEqual(expected1.num)))
	t.Run("Expected correct result", subx.Test(subx.Value(result2.num), subx.CompareEqual(expected2.num)))

	err1 = ecs.RemoveComponent[comp1](&entity)
	err2 = ecs.RemoveComponent[comp2](&entity)

	t.Run("Expected correct result", subx.Test(subx.Value(err1), subx.CompareEqual[error](nil)))
	t.Run("Expected correct result", subx.Test(subx.Value(err2), subx.CompareEqual[error](nil)))

	_, err1 = ecs.GetComponent[comp1](&entity)
	_, err2 = ecs.GetComponent[comp2](&entity)

	t.Run("Expected correct result", subx.Test(subx.Value(err1), subx.CompareNotEqual[error](nil)))
	t.Run("Expected correct result", subx.Test(subx.Value(err2), subx.CompareNotEqual[error](nil)))
}

func TestEntityAddGetRemoveComponentMany(t *testing.T) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 5)

	type comp struct {
		num int
	}

	for n := 0; n < 5; n++ {
		entities[n] = scene.NewEntity()
		err := ecs.AddComponent(&entities[n], &comp{num: n})
		t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareEqual[error](nil)))
	}

	for n := 0; n < 5; n++ {
		result, err := ecs.GetComponent[comp](&entities[n])
		t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareEqual[error](nil)))
		t.Run("Expected correct result", subx.Test(subx.Value(result.num), subx.CompareEqual(n)))
	}

	for n := 0; n < 5; n++ {
		err := ecs.RemoveComponent[comp](&entities[n])
		t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareEqual[error](nil)))
	}

	for n := 0; n < 5; n++ {
		_, err := ecs.GetComponent[comp](&entities[n])
		t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareNotEqual[error](nil)))
	}
}

func TestAllComponents(t *testing.T) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 5)

	type comp struct {
		num int
	}

	for n := 0; n < 5; n++ {
		entities[n] = scene.NewEntity()
		err := ecs.AddComponent(&entities[n], &comp{num: n})
		t.Run("Expected correct result", subx.Test(subx.Value(err), subx.CompareEqual[error](nil)))

	}
	allComponents := ecs.AllComponents[comp](&scene)
	for n := 0; n < 5; n++ {
		num := allComponents[n].Component().num
		t.Run("Expected correct result", subx.Test(subx.Value(num), subx.CompareEqual(n)))
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
		num int
	}

	type comp2 struct {
		num int
	}

	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&entities[n], &comp1{num: 1})
			ecs.AddComponent(&entities[n], &comp2{num: 2})
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
		num int
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&entities[n], &comp{num: n})
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
		num int
	}

	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			ecs.AddComponent(&entities[n], &comp{num: n})
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
		num int
	}

	for n := 0; n < 100; n++ {
		entities[n] = scene.NewEntity()
		ecs.AddComponent(&entities[n], &comp{num: n})
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for n := 0; n < 100; n++ {
			ecs.GetComponent[comp](&entities[n])
		}
		b.StopTimer()
	}
}

func BenchmarkAllComponents(b *testing.B) {
	scene := ecs.Scene{}
	entities := make([]ecs.Entity, 100)

	type comp struct {
		num int
	}

	for i := 0; i < len(entities); i++ {
		entities[i] = scene.NewEntity()
		ecs.AddComponent(&entities[i], &comp{num: i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < 100; n++ {
			allComponents := ecs.AllComponents[comp](&scene)
			for i := 0; i < len(allComponents); i++ {
				allComponents[i].Component().num++
			}
		}
	}
}
