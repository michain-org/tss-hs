package sign

import (
	"crypto/sha512"
	"math/big"

	"github.com/binance-chain/edwards25519/edwards25519"
	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/decred/dcrd/dcrec/edwards"
	"github.com/michain-org/tss-hs/keygen"
	kg "github.com/michain-org/tss-hs/keygen"
	"github.com/pkg/errors"
)

type (
	TempDate struct {
		Wi, R, SumR *big.Int
		Ri          *crypto.ECPoint
		Msg         []byte
		Si          *[32]byte
	}
)

func GetIndex(sv *keygen.LocalPartySaveData) (int, error) {
	ks := sv.Ks
	id := sv.ShareID
	for i, v := range ks {
		if v.Cmp(id) == 0 {
			return i, nil
		}
	}
	return -1, errors.New("Use a bad SaveData for preparing")
}

func Rand() (*big.Int, *crypto.ECPoint) {
	ri := common.GetRandomPositiveInt(tss.EC().Params().N)
	pointRi := crypto.ScalarBaseMult(tss.EC(), ri)
	return ri, pointRi
}

func CalcWi(idx int, idxs []int, sv *kg.LocalPartySaveData) (*big.Int, error) {
	ks := sv.Ks
	shareID := sv.ShareID
	modQ := common.ModInt(tss.EC().Params().N)
	maxIdx := len(ks) - 1
	wi := sv.Xi
	for _, v := range idxs {
		if v > maxIdx || v < 0 {
			return nil, errors.New("Get a out-of-range index when calcWi")
		}
		if v == idx {
			continue
		}
		coef := modQ.Mul(ks[v], modQ.ModInverse(new(big.Int).Sub(ks[v], shareID)))
		wi = modQ.Mul(wi, coef)
	}
	return wi, nil
}

func CalcSi(ris []*crypto.ECPoint, tmp *TempDate, sv *kg.LocalPartySaveData) (*[32]byte, *big.Int, error) {
	if tmp.Wi == nil ||
		tmp.R == nil ||
		tmp.Msg == nil ||
		tmp.Ri == nil {
		return nil, nil, errors.New("a insufficient paramters temp")
	}
	var R edwards25519.ExtendedGroupElement
	riBytes := bigIntToEncodedBytes(tmp.R)
	edwards25519.GeScalarMultBase(&R, riBytes)
	for _, Rj := range ris {
		if tmp.Ri.Equals(Rj) {
			continue
		}
		extendedRj := ecPointToExtendedElement(Rj.X(), Rj.Y())
		R = addExtendedElements(R, extendedRj)
	}
	var encodedR [32]byte
	R.ToBytes(&encodedR)
	encodedPubKey := ecPointToEncodedBytes(sv.EDDSAPub.X(), sv.EDDSAPub.Y())

	// h = hash512(k || A || M)
	h := sha512.New()
	h.Reset()
	h.Write(encodedR[:])
	h.Write(encodedPubKey[:])
	h.Write(tmp.Msg)

	var lambda [64]byte
	h.Sum(lambda[:0])
	var lambdaReduced [32]byte
	edwards25519.ScReduce(&lambdaReduced, &lambda)

	// 8. compute si
	var localS [32]byte
	edwards25519.ScMulAdd(&localS, &lambdaReduced, bigIntToEncodedBytes(tmp.Wi), riBytes)

	Si := &localS
	SumR := encodedBytesToBigInt(&encodedR)
	return Si, SumR, nil
}

func Mpk(sv *kg.LocalPartySaveData) *edwards.PublicKey {
	return &edwards.PublicKey{
		Curve: tss.EC(),
		X:     sv.EDDSAPub.X(),
		Y:     sv.EDDSAPub.Y(),
	}
}

func Verify(mpk *edwards.PublicKey, s *[32]byte, r *big.Int, msg []byte) bool {
	return edwards.Verify(mpk, msg, r, encodedBytesToBigInt(s))
}

func Sums(ss []*[32]byte) (sum *[32]byte) {
	switch len(ss) {
	case 0:
		break
	case 1:
		sum = ss[0]
	default:
		sum = ss[0]
		for j := range ss {
			if j == 0 {
				continue
			}
			sjBytes := ss[j]
			var tmp [32]byte
			edwards25519.ScMulAdd(&tmp, sum, bigIntToEncodedBytes(big.NewInt(1)), sjBytes)
			sum = &tmp
		}
	}
	return
}
