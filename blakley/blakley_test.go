package blakley

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_Mult(t *testing.T) {
	for i := 1; i < 256; i++ {
		assert.Equal(t, mult(uint8(i), invTabel[i]), uint8(0x1))
	}
	assert.Equal(t, mult(0, 4), uint8(0))
	if out := mult(3, 3); out != 3 {
		t.Fatalf("Bad: %v 3", out)
	}
}

func TestField_Div(t *testing.T) {
	if out := div(9, 7); out != 3 {
		t.Fatalf("Bad: %v 3", out)
	}
	if out := div(9, 3); out != 7 {
		t.Fatalf("Bad: %v 7", out)
	}
	if out := div(0, 3); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}
	if out := div(2, 2); out != 1 {
		t.Fatalf("Bad: %v 1", out)
	}
	/*if out := div(7, 0); out == 0 {
		t.Fatalf("Very Bad")
	}*/
}

func TestEvaluate(t *testing.T) {
	a := []uint8{3, 3, 3}
	b := []uint8{7, 7, 7}
	c := []uint8{5, 6, 2}
	d := []uint8{0x52, 0x7B, 0}
	if out := evaluate(a, b, 3); out != 7 {
		t.Fatalf("Bad: %v 7", out)
	}
	if out := evaluate(d, c, 3); out != 2 {
		t.Fatalf("Bad: %v 2", out)
	}
}

func TestSolve_invalid(t *testing.T) {
	matrixA := [][]byte{
		{3, 3},
		{2, 1, 1},
		{7, 7, 3},
	}

	matrixB := [][]byte{
		{2, 1, 1},
		{7, 7, 3},
	}

	b := []byte{3, 2, 3}
	_, err := solve(matrixA, b, 3)
	if err == nil {
		t.Fatalf("Bad")
	}
	_, err = solve(matrixB, b, 3)
	if err == nil {
		t.Fatalf("Bad")
	}
}

func TestSolve(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
	}

	b := []byte{0, 5, 2}
	s, _ := solve(matrix, b, 3)
	assert.Equal(t, uint8(0x1), s)
}

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
	out, err := Split(secret, 8, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(out) != 8 {
		t.Fatalf("bad: %v", len(out))
	}
	for i := range out {
		if len(out[i]) != len(out[0]) {
			t.Fatalf("bad: %v,%v,%v", len(out[i]), len(out[0]), i)
		}
		if len(out[i]) != len(secret) {
			t.Fatalf("bad:%v,%v", len(out[i]), len(secret))
		}
		for j := range out[i] {
			if len(out[i][j]) != 3 {
				t.Fatalf("bad: %v", len(out[i][j]))
			}
		}
	}
}

func TestCombine_invalid(t *testing.T) {
	// Not enough parts
	if _, err := Combine(nil); err == nil {
		t.Fatalf("should err")
	}

	partA := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{1, 2, 3},
	}
	partB := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{1, 2, 3},
		{2, 3},
	}
	partC := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{1, 2, 3, 4},
		{2, 3, 4, 5},
	}
	partsA := [][][]byte{
		partA,
		partB,
	}
	partsB := [][][]byte{
		partA,
	}
	partsC := [][][]byte{
		partB,
		partC,
	}

	//Mismatch in the length of secret
	if _, err := Combine(partsA); err == nil {
		t.Fatalf("should err")
	}

	//Mismatch in shares
	if _, err := Combine(partsC); err == nil {
		t.Fatalf("should err")
	}

	//The number of parties < 2
	if _, err := Combine(partsB); err == nil {
		t.Fatalf("should err")
	}
}

func TestCombine(t *testing.T) {
	secret := []byte("VGhpcyBpcyBhIHNpbXBsZSB0ZXN0IQpBbmQgSSB3YW50IHRvIGRyaW5rIGEgY3VwIG9mIHBvcDop")
	shares, err := Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	s, _ := Combine(shares)
	if len(s) != len(secret) {
		t.Fatalf("Bad: %v", s)
	}
	assert.Equal(t, secret, s)
}
