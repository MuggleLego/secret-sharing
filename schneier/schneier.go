package schneier

import (
	"crypto/rand"
	"fmt"
)

const (
	// ShareOverhead is the byte size overhead of each share
	// when using Split on a secret. This is caused by appending
	// a one byte tag to the share.
	ShareOverhead = 1
)

// Split using a XOR-based algorithm. It takes an arbitrarily long secret
// and generates a `parts` number of shares, `threshold` of which are required to
// reconstruct the secret. The parts and threshold must be equal, and must be at least 2,
// and less than 256
func Split(secret []byte, parts, threshold int) ([][]byte, error) {
	// Sanity check the input
	if parts < threshold {
		return nil, fmt.Errorf("parts cannot be less than threshold")
	}
	if parts > 255 {
		return nil, fmt.Errorf("parts cannot exceed 255")
	}
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}
	if threshold > 255 {
		return nil, fmt.Errorf("threshold cannot exceed 255")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("cannot split an empty secret")
	}
	if parts != threshold {
		return nil, fmt.Errorf("parts and threshold must be equal")
	}

	//Initialization of returned sub-secret
	b := len(secret)
	shares := make([][]byte, parts)
	shares[parts-1] = secret

	for i := 0; i < parts-1; i++ {
		//Randomly sample r_0,r_1,...,r_{n-1}
		//where |r_i| == |s|
		shares[i] = make([]byte, b)
		if _, err := rand.Read(shares[i]); err != nil {
			return nil, fmt.Errorf("initialize failed")
		}
		for j := range shares[i] {
			//r_n = s ^ r_0 ^ r_1 ^ ...	^ r_{n-1}
			shares[parts-1][j] ^= shares[i][j]
		}
	}

	return shares, nil
}

func Combine(shares [][]byte) ([]byte, error) {
	//Sanity check for the shares
	if shares == nil {
		return nil, fmt.Errorf("cannot combine nil shares")
	}
	if len(shares) < 2 {
		return nil, fmt.Errorf("cannot combine less than 2 shares")
	}
	for i := range shares {
		if shares[i] == nil {
			return nil, fmt.Errorf("the %vth share invalid,cannot combine nil shares", i)
		}
		if len(shares[i]) != len(shares[0]) {
			return nil, fmt.Errorf("invalid shares provided:party %v", i)
		}
	}

	b := len(shares)
	for i := 0; i < b-1; i++ {
		//Secret = (s ^ r_0 ^ r_1 ^ ... ^ r_{n-1}) ^ r_0 ^ r_1 ^ ... ^ r_{n-1} = s
		for j := range shares[0] {
			shares[b-1][j] ^= shares[i][j]
		}
	}
	return shares[b-1], nil
}
