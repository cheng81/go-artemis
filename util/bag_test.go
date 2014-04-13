package util

import (
	// "fmt"
	"testing"
)

func TestBagCreate(t *testing.T) {
	b := NewBag(64)
	if b.Capacity() != 64 {
		t.Error("Bag capacity should be 64", b.Capacity())
	}
	if b.Size() != 0 {
		t.Error("Bag size should be 0", b.Size())
	}

	b.Add("foo")
	if b.Size() != 1 {
		t.Error("Bag size should be 1", b.Size())
	}

	el := b.GetAt(0)
	if el != "foo" {
		t.Error("Bag.GetAt(0) should be 'foo'", el)
	}

	b.SetAt(100, "bar")
	if b.Capacity() < 100 {
		t.Error("Bag.Capacity should be >= 100", b.Capacity())
	}
}

func TestBagSize(t *testing.T) {
	b := NewBag(64)
	if b.Size() != 0 {
		t.Error("Bag.Size should be 0")
	}

	b.Add(1)
	if b.Size() != 1 {
		t.Error("Bag.Size should be 1")
	}

	last := b.RemoveLast().(int)

	if last != 1 {
		t.Error("Bag.RemoveLast should be 1")
	}
	if b.Size() != 0 {
		t.Error("Bag.Size should be 0")
	}

}
