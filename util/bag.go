package util

type ImmutableBag interface {
	GetAt(uint) interface{}
	GetAtSafe(uint) (interface{}, bool)
	Size() uint
	Empty() bool
	Contains(interface{}) bool
	ForEach(fn func(int, interface{}))
}

var empty = emptyBag(struct{}{})

type emptyBag struct{}

func (_ emptyBag) GetAt(_ uint) interface{}             { return nil }
func (_ emptyBag) GetAtSafe(_ uint) (interface{}, bool) { return nil, false }
func (_ emptyBag) Size() uint                           { return 0 }
func (_ emptyBag) Empty() bool                          { return true }
func (_ emptyBag) Contains(_ interface{}) bool          { return false }
func (_ emptyBag) ForEach(_ func(int, interface{}))     {}

func EmptyBag() ImmutableBag {
	return empty
}

type Bag struct {
	data []interface{}
	size uint
}

func NewBag(size int) (out *Bag) {
	out = new(Bag)
	out.data = make([]interface{}, size, size)
	out.size = 0
	return
}

func (b *Bag) RemoveAt(index uint) (out interface{}) {
	out = b.data[index]
	b.size -= 1
	b.data[index] = b.data[b.size]
	b.data[b.size] = nil
	return
}

func (b *Bag) RemoveLast() (out interface{}) {
	if b.size >= 0 {
		b.size -= 1
		out = b.data[b.size]
		b.data[b.size] = nil
		return
	}
	return nil
}

func (b *Bag) RemoveElem(e interface{}) bool {
	for i := uint(0); i < b.size; i++ {
		if b.data[i] == e {
			b.size -= 1
			b.data[i] = b.data[b.size]
			b.data[b.size] = nil
			return true
		}
	}
	return false
}

func (b *Bag) RemoveAll(elems ...interface{}) (modified bool) {
	modified = false
	for _, el := range elems {
		for i := uint(0); i < b.size; i++ {
			if b.data[i] == el {
				b.RemoveElem(el)
				i -= 1
				modified = true
				break
			}
		}
	}
	return
}

func (b *Bag) Add(el interface{}) {
	if b.size == uint(len(b.data)) {
		b.grow()
	}
	b.data[b.size] = el
	b.size += 1
}
func (b *Bag) SetAt(index uint, el interface{}) {
	if index >= uint(len(b.data)) {
		b.growTo(int(index) * 2)
		b.size = index + 1
	}
	b.data[index] = el
}

func (b *Bag) GetAt(index uint) interface{} {
	if index >= uint(len(b.data)) {
		return nil
	}
	return b.data[index]
}
func (b *Bag) GetAtSafe(index uint) (interface{}, bool) {
	res := b.GetAt(index)
	return res, res != nil
}

func (b *Bag) Size() uint    { return b.size }
func (b *Bag) Capacity() int { return cap(b.data) }
func (b *Bag) Empty() bool   { return b.size == 0 }

func (b *Bag) grow() {
	b.growTo(int((len(b.data) * 3 / 2) + 1))
}
func (b *Bag) growTo(newCapacity int) {
	grow := newCapacity - len(b.data)
	b.data = append(b.data, make([]interface{}, grow, grow)...)
}

func (b *Bag) Contains(item interface{}) bool {
	for _, el := range b.data {
		if el == item {
			return true
		}
	}
	return false
}

func (b *Bag) Clear() {
	for i, _ := range b.data {
		b.data[i] = nil
	}
}

func (b *Bag) AddAll(elems ...interface{}) {
	for _, el := range elems {
		b.Add(el)
	}
}

func (b *Bag) EnsureCapacity(index uint) {
	if index >= uint(len(b.data)) {
		b.growTo(int(index) * 2)
		b.size = index + 1
	}
}

func (b *Bag) ForEach(fn func(int, interface{})) {
	for i, el := range b.data {
		if el != nil {
			fn(i, el)
		}
	}
}
