package schneier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit_invalid(t *testing.T) {
	secret := []byte("test")

	if _, err := Split(secret, 0, 0); err == nil {
		t.Fatalf("expect error,case 1")
	}

	if _, err := Split(secret, 2, 3); err == nil {
		t.Fatalf("expect error,case 2")
	}

	if _, err := Split(secret, 1000, 3); err == nil {
		t.Fatalf("expect error,case 3")
	}

	if _, err := Split(secret, 10, 1); err == nil {
		t.Fatalf("expect error,case 4")
	}

	if _, err := Split(nil, 3, 2); err == nil {
		t.Fatalf("expect error,case 5")
	}
}

func TestSplit(t *testing.T) {
	secret := []byte("just a simple test")
	out, err := Split(secret, 3, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(out) != 3 {
		t.Fatalf("bad: %v", len(out))
	}
	for i := range out {
		if len(out[i]) != len(secret) {
			t.Fatalf("bad:%v,%v", len(out[i]), len(secret))
		}
	}
}

func TestCombine_invalid(t *testing.T) {
	// Not enough parts
	if _, err := Combine(nil); err == nil {
		t.Fatalf("should err")
	}

	// Mismatch in length
	partA := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{1, 2, 3},
	}

	partB := [][]byte{
		{1, 2, 3, 4},
	}
	//Mismatch in the length of secret
	if _, err := Combine(partA); err == nil {
		t.Fatalf("should err")
	}
	//The number of parties < 2
	if _, err := Combine(partB); err == nil {
		t.Fatalf("should err")
	}
}

func TestCombine(t *testing.T) {
	secret := []byte("This is a simple test")
	shares, err := Split(secret, 5, 5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	s, _ := Combine(shares)
	if len(s) != len(secret) {
		t.Fatalf("Bad: %v", s)
	}
	assert.Equal(t, secret, s)
}
