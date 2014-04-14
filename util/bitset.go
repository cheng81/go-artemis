package util

import (
	"fmt"
)

type BitSet struct{ arr []uint64 }

func NewBitSet(size int) (out *BitSet) {
	sz := int(size/64) + 1
	out = new(BitSet)
	out.arr = make([]uint64, sz, sz+1)
	return
}

func (b *BitSet) ToStr() string {
	str := "BitSet["
	b.WithSetBitsFrom(0, func(idx uint) {
		str = fmt.Sprint(str, idx, ",")
	})
	return fmt.Sprint(str, "]")
}

func (b *BitSet) Intersects(o *BitSet) bool {
	min := len(b.arr)
	if len(o.arr) < min {
		min = len(o.arr)
	}
	for i := 0; i < min; i++ {
		if (b.arr[i] & o.arr[i]) != 0 {
			return true
		}
	}
	return false
}

func (b *BitSet) Empty() bool {
	for _, c := range b.arr {
		if c != 0 {
			return false
		}
	}
	return true
}

func (b *BitSet) Clear() {
	for i, _ := range b.arr {
		b.arr[i] = 0
	}
}

func (b *BitSet) indexOf(index uint) (uint, uint) {
	if index < 64 {
		return 0, index
	}
	return uint(index / 64), uint(index % 64)
}
func (b *BitSet) ensureSize(sz uint) {
	if uint(len(b.arr)) < sz+1 {
		grow := int(sz+(sz*3/2)) - len(b.arr)
		b.arr = append(b.arr, make([]uint64, grow, grow)...)
	}
}
func (b *BitSet) SetV(index uint, val bool) {
	chunk, pos := b.indexOf(index)
	b.ensureSize(chunk)
	if val {
		b.arr[chunk] |= (1 << pos)
	} else {
		b.arr[chunk] &^= (1 << pos)
	}
}

func (b *BitSet) Set(index uint) {
	b.SetV(index, true)
}
func (b *BitSet) Unset(index uint) {
	b.SetV(index, false)
}

func (b *BitSet) Get(index uint) bool {
	if uint(len(b.arr)*64) < index {
		return false
	}
	chunk, pos := b.indexOf(index)
	return b.get(b.arr[chunk], pos)
}

func (b *BitSet) get(chunk uint64, pos uint) bool {
	return (chunk & (1 << pos)) != 0
}

func (b *BitSet) unindex(chunk, pos uint) uint {
	return (chunk * 64) + pos
}

func (b *BitSet) NextSetBit(index uint) (uint, bool) {
	chunk_index, pos := b.indexOf(index)
	alen := uint(len(b.arr))
	for chunk_index < alen {
		chunk := b.arr[chunk_index]
		for pos < 64 {
			if b.get(chunk, pos) {
				return b.unindex(chunk_index, pos), true
			}
			pos += 1
		}
		chunk_index += 1
		pos = 0
	}
	return 0, false
}

func (b *BitSet) WithSetBitsFrom(index uint, fn func(uint)) {
	chunk_index, pos := b.indexOf(index)
	alen := uint(len(b.arr))
	for ; chunk_index < alen; chunk_index += 1 {
		chunk := b.arr[chunk_index]
		for ; pos < 64; pos += 1 {
			if b.get(chunk, pos) {
				fn(b.unindex(chunk_index, pos))
			}
		}
		pos = 0
	}
}
