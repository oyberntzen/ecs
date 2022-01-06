package ecs3

const (
	increaseFactor    = 2
	decreaseFactor    = 2
	decreaseThreshold = 3
)

type item struct {
	entity Entity
	data   interface{}
}

type pool struct {
	items    []item
	indicies map[Entity]uint32
}

func (p *pool) Add(e Entity, i item) {
	length := len(p.items)
	if cap(p.items)-length == 0 {
		if length == 0 {
			p.items = []item{i}
			p.indicies[e] = 0
			return
		}
		newItems := make([]item, length+1, length*increaseFactor)
		copy(newItems, p.items)
		p.items = newItems
		p.items[length] = i
		p.indicies[e] = uint32(length)
		return
	}
	p.items = p.items[:length+1]
	p.items[length] = i
	p.indicies[e] = uint32(length)
}

func (p *pool) Remove(entity Entity) bool {
	index, ok := p.indicies[entity]
	if !ok {
		return false
	}
	delete(p.indicies, entity)

	length := len(p.items)
	p.items[index] = p.items[length-1]
	p.items = p.items[:length-1]
	length--

	if uint32(length) > index {
		p.indicies[p.items[index].entity] = index
	}

	if length*decreaseThreshold < cap(p.items) {
		newItems := make([]item, length, length*decreaseFactor)
		copy(newItems, p.items)
		p.items = newItems
	}

	return true
}
