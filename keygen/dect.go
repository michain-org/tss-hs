package keygen

import (
	"math/big"

	tssComm "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	kg "github.com/binance-chain/tss-lib/eddsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/pkg/errors"
)

func makePolynomial(ids []*big.Int, Threshold int) (*LocalPartySaveData, vss.Vs, vss.Shares, error) {
	partyCount := len(ids)
	s := kg.NewLocalPartySaveData(partyCount)
	save := &s
	ui := tssComm.GetRandomPositiveInt(tss.EC().Params().N)
	vs, shares, err := vss.Create(Threshold, ui, ids)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "create polynomial fail")
	}
	return convert(save), vs, shares, nil
}

func PubPoly(dimens int, vss ...vss.Vs) (vss.Vs, error) {
	if vss != nil && len(vss) != 0 {
		return nil, errors.New("vss should be suffitient")
	}
	for i := range vss {
		if len(vss[i]) != dimens {
			return nil, errors.New("vss should be with curtain dimension")
		}
	}
	poly := make([]*crypto.ECPoint, dimens)
	var e error
	for i := range vss {
		for j := range vss[i] {
			if i == 0 {
				poly[j] = vss[i][j]
			} else {
				poly[j], e = poly[j].Add(vss[i][j])
				if e != nil {
					return nil, errors.Errorf("could not add ECPoint at coordinate i:%s,j:%s", i, j)
				}
			}
		}
	}
	return poly, nil
}

func CalSecXi(Xis []*big.Int) *big.Int {
	Xi := new(big.Int)
	for i := range Xis {
		Xi = Xi.Add(Xi, Xis[i])
	}
	return Xi
}

func ValShare(poly []*crypto.ECPoint, Xi, share *big.Int) bool {
	var po *crypto.ECPoint
	var e error
	x := new(big.Int).SetInt64(1)
	for i := range poly {
		if i == 0 {
			po, _ = crypto.NewECPoint(tss.EC(), poly[0].X(), poly[0].Y())
		} else {
			x = x.Mul(x, x)
			p := poly[i].ScalarMult(x)
			po, e = po.Add(p)
			if e != nil {
				return false
			}
		}
	}
	Po := crypto.ScalarBaseMult(tss.EC(), share)
	return Po.Equals(po)
}
