package signing

import (
	"crypto/sha512"
	"math/big"

	"github.com/binance-chain/edwards25519/edwards25519"
	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/decred/dcrd/dcrec/edwards"
	kg "github.com/michain-org/tss-hs/keygen"
	"github.com/pkg/errors"
)

type (
	TempDate struct {
		wi, r, sumr *big.Int
		R           *crypto.ECPoint
		msg         []byte
		si          *[32]byte
	}

	Base struct {
		*kg.LocalPartySaveData
		index int
		temp  *TempDate
	}
)

func prepare(sv *kg.LocalPartySaveData) int {
	ks := sv.Ks
	id := sv.ShareID
	for i, v := range ks {
		if v.Cmp(id) == 0 {
			return i
		}
	}
	panic("Use a bad SaveData for preparing")
}

func randr() (*big.Int, *crypto.ECPoint) {
	ri := common.GetRandomPositiveInt(tss.EC().Params().N)
	pointRi := crypto.ScalarBaseMult(tss.EC(), ri)
	return ri, pointRi
}

func (t *TempDate) reset() {
	r, point := randr()
	t.r = r
	t.R = point
}

func newBase(sv *kg.LocalPartySaveData) *Base {
	return &Base{LocalPartySaveData: sv,
		index: prepare(sv),
		temp:  (*TempDate)(nil)}
}

// make a completely new temp for base
func (b *Base) newRound(msg []byte) {
	m := make([]byte, len(msg))
	copy(m, msg)
	b.temp = new(TempDate)
	b.temp.msg = msg
	b.temp.reset()
}

func (b *Base) calcWi(idxs []int) error {
	ks := b.Ks
	shareID := b.ShareID
	modQ := common.ModInt(tss.EC().Params().N)
	maxIdx := len(ks) - 1
	wi := b.Xi
	for _, v := range idxs {
		if v > maxIdx {
			return errors.New("Get a out-of-range index when calcWi")
		}
		if v == b.index {
			continue
		}
		coef := modQ.Mul(ks[v], modQ.ModInverse(new(big.Int).Sub(ks[v], shareID)))
		wi = modQ.Mul(wi, coef)
	}
	b.temp.wi = wi
	return nil
}

func (b *Base) calcSi(rs []*crypto.ECPoint) {
	var R edwards25519.ExtendedGroupElement
	riBytes := bigIntToEncodedBytes(b.temp.r)
	edwards25519.GeScalarMultBase(&R, riBytes)
	for _, Rj := range rs {
		if b.temp.R.Equals(Rj) {
			// fmt.Println("break")
			continue
		}
		extendedRj := ecPointToExtendedElement(Rj.X(), Rj.Y())
		R = addExtendedElements(R, extendedRj)
	}
	var encodedR [32]byte
	R.ToBytes(&encodedR)
	encodedPubKey := ecPointToEncodedBytes(b.EDDSAPub.X(), b.EDDSAPub.Y())

	// h = hash512(k || A || M)
	h := sha512.New()
	h.Reset()
	h.Write(encodedR[:])
	h.Write(encodedPubKey[:])
	h.Write(b.temp.msg)

	var lambda [64]byte
	h.Sum(lambda[:0])
	var lambdaReduced [32]byte
	edwards25519.ScReduce(&lambdaReduced, &lambda)

	// 8. compute si
	var localS [32]byte
	edwards25519.ScMulAdd(&localS, &lambdaReduced, bigIntToEncodedBytes(b.temp.wi), riBytes)

	b.temp.si = &localS
	b.temp.sumr = encodedBytesToBigInt(&encodedR)
}

func (b *Base) Verify(s *[32]byte) bool {
	pk := edwards.PublicKey{
		Curve: tss.EC(),
		X:     b.EDDSAPub.X(),
		Y:     b.EDDSAPub.Y(),
	}

	return edwards.Verify(&pk,
		b.temp.msg,
		b.temp.sumr,
		encodedBytesToBigInt(s),
	)
}

func SumSis(sis []*[32]byte) (sumS *[32]byte) {
	switch len(sis) {
	case 0:
		break
	case 1:
		sumS = sis[0]
	default:
		sumS = sis[0]
		for j := range sis {
			if j == 0 {
				continue
			}
			sjBytes := sis[j]
			var tmpSumS [32]byte
			edwards25519.ScMulAdd(&tmpSumS, sumS, bigIntToEncodedBytes(big.NewInt(1)), sjBytes)
			sumS = &tmpSumS
		}
	}
	return
}
