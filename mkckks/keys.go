package mkckks

import "github.com/oksuman/snu-mghe/mkrlwe"

// NewKeyGenerator creates a rlwe.KeyGenerator instance from the CKKS parameters.
func NewKeyGenerator(params Parameters) *mkrlwe.KeyGenerator {
	return mkrlwe.NewKeyGenerator(params.Parameters)
}
