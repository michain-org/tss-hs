package common

import "math/big"

func Ints2bytes(bigNums []*big.Int) [][]byte {
	b := make([][]byte, len(bigNums))
	for i, s := range bigNums {
		if s == nil {
			continue
		}
		b[i] = s.Bytes()
	}
	return b
}

func Bytes2ints(bs [][]byte) []*big.Int {
	bn := make([]*big.Int, len(bs))
	for i := range bs {
		bn[i] = new(big.Int).SetBytes(bs[i])
	}
	return bn
}

func NonEmptyMulBytes(bs [][]byte) bool {
	if bs == nil || len(bs) == 0 {
		return false
	}
	for i := range bs {
		if !NonEmptyBytes(bs[i]) {
			return false
		}
	}
	return true
}

func NonEmptyBytes(bs []byte) bool {
	return bs != nil && len(bs) != 0
}
