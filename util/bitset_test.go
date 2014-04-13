package util

import (
	"fmt"
	"testing"
)

func TestBitSet_0(t *testing.T) {
	bs := NewBitSet(100)
	for i := uint(0); i < 100; i++ {
		if bs.Get(i) {
			t.Error("BitSet[", i, "] should be false")
		}
	}
}

func TestBitSet_str(t *testing.T) {
	b := NewBitSet(64)
	if b.Get(0) {
		t.Error("BitSet[0] should be false")
	}
	fmt.Println(b.ToStr())
	b.Set(0)
	fmt.Println(b.ToStr())
}

func TestBitSet_Clear(t *testing.T) {
	bs := NewBitSet(30)
	bs.Set(15)

	bs.Clear()

	if bs.Get(15) {
		t.Error("BitSet[15] after Clear should be false")
	}
}

func TestBitSet_SetUnset(t *testing.T) {
	bs := NewBitSet(10)
	bs.Set(5)
	if !bs.Get(5) {
		t.Error("BitSet[5] should be true")
	}
	bs.Unset(5)
	if bs.Get(5) {
		t.Error("BitSet[5] should be false")
	}
}

func TestBitSet_Growing(t *testing.T) {
	bs := NewBitSet(20)
	bs.Set(100)
	if !bs.Get(100) {
		t.Error("BitSet[100] should be true")
	}
	bs.Set(1000)
	if !bs.Get(1000) {
		t.Error("BitSet[1000] should be true")
	}
}

func TestBitSet_NextSet(t *testing.T) {
	bs := NewBitSet(50)
	bs.Set(10)
	bs.Set(30)

	if index, found := bs.NextSetBit(0); index != 10 || !found {
		t.Error("BitSet.NextSetBit(0) should be 10, true")
	}

	if index, found := bs.NextSetBit(11); index != 30 || !found {
		t.Error("BitSet.NextSetBit(11) should be 30, true")
	}

	if index, found := bs.NextSetBit(31); index != 0 || found {
		t.Error("BitSet.NextSetBit(31) should be 0, false")
	}

	bs.WithSetBitsFrom(0, func(idx uint) {
		fmt.Println("Found bit set at", idx)
	})
}

func TestBitSet_Intersect(t *testing.T) {
	b1 := NewBitSet(64)
	b2 := NewBitSet(256)

	b1.Set(10)

	if b1.Intersects(b2) {
		t.Error("b1 should not intersect b2")
	}

	b2.Set(130)
	if b2.Intersects(b1) {
		t.Error("b2 should not intersect b1")
	}

	b2.Set(10)
	if !b1.Intersects(b2) {
		t.Error("b1 should intersect b2")
	}

}
