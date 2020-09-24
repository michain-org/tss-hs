package keygen

import (
	"math/big"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	kg "github.com/binance-chain/tss-lib/eddsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/decred/dcrd/dcrec/edwards"
	"github.com/pkg/errors"
)

func makePolynomial(ids []*big.Int, Threshold int) (*kg.LocalPartySaveData, vss.Vs, vss.Shares, error) {
	partyCount := len(ids)
	s := kg.NewLocalPartySaveData(partyCount)
	save := &s
	ui := common.GetRandomPositiveInt(tss.EC().Params().N)
	vs, shares, err := vss.Create(Threshold, ui, ids)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "create polynomial fail")
	}
	return save, vs, shares, nil

}

// Local generate partySaveDate Centrally
// Threshold = Participants - (Participants - 1)/3 -1
func Local(Threshold int, Participants int) ([]*kg.LocalPartySaveData, error) {
	tss.SetCurve(edwards.Edwards())
	pIDs, err := GenPartyIDs(Participants)
	if err != nil {
		return nil, err
	}
	ids := pIDs.Keys()
	idsLen := len(ids)
	vc := make(vss.Vs, Threshold+1)
	Xi := make([]*big.Int, idsLen)
	allsv := make([]*kg.LocalPartySaveData, idsLen)

	for i := 0; i < idsLen; i++ {
		sv, vs, shares, err := makePolynomial(ids, Threshold)
		if err != nil {
			return nil, errors.Wrap(err, "make polynomial fail")
		}
		allsv[i] = sv
		for j := 0; j < len(shares); j++ {
			if i == 0 {
				Xi[j] = shares[j].Share
				continue
			}
			Xi[j] = new(big.Int).Add(Xi[j], shares[j].Share)
		}
		for j := 0; j < len(vs); j++ {
			if i == 0 {
				vc[j] = vs[j]
				continue
			}
			if vc[j], err = vc[j].Add(vs[j]); err != nil {
				return nil, errors.Wrapf(err, "vs can not Add for point %v", vs[j])
			}
		}
	}
	_x := vc[0].X()
	_y := vc[0].Y()
	EddsaPub, _ := crypto.NewECPoint(tss.EC(), _x, _y)

	//calc BigXi
	bigXi := make([]*crypto.ECPoint, idsLen)
	modQ := common.ModInt(tss.EC().Params().N)

	for idx, id := range ids {
		bigX, _ := crypto.NewECPoint(tss.EC(), vc[0].X(), vc[0].Y())
		z := new(big.Int).SetInt64(int64(1))
		for i := 0; i < len(vc); i++ {
			z = modQ.Mul(z, id)
			bigX, err = bigX.Add(vc[i].ScalarMult(z))
			if err != nil {
				return nil, errors.Wrapf(err, "calculate bigXi err id %d vc %v", id, vc[i])
			}
		}
		bigXi[idx] = bigX
	}

	//set saveData
	for i, sv := range allsv {
		sv.ShareID = ids[i]
		sv.Xi = Xi[i]
		sv.Ks = ids
		sv.BigXj = bigXi[:]
		sv.EDDSAPub = EddsaPub
	}

	return allsv, nil
}
