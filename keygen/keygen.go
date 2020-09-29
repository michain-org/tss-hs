package keygen

import (
	"math/big"

	tssComm "github.com/binance-chain/tss-lib/common"
	"github.com/decred/dcrd/dcrec/edwards"
	"github.com/michain-org/tss-hs/common"
	"github.com/pkg/errors"

	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	kg "github.com/binance-chain/tss-lib/eddsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

const (
	path = "./paramer"
)

func makeParty(
	para *Paramer,
	out chan<- tss.Message,
	end chan<- kg.LocalPartySaveData,
) *kg.LocalParty {
	p2pCtx := tss.NewPeerContext(para.Pids)
	params := tss.NewParameters(p2pCtx, para.Pids[para.Index], len(para.Pids), para.Threshold)
	P := kg.NewLocalParty(params, out, end).(*kg.LocalParty)
	return P
}

//Decentralized generate savedata in a decentralized way
func Decentralized(
	para *Paramer,
	in <-chan tss.Message,
	out chan<- tss.Message,
) (*LocalPartySaveData, error) {

	errCh := make(chan *tss.Error)
	outCh := make(chan tss.Message)
	endCh := make(chan kg.LocalPartySaveData)
	party := makeParty(para, outCh, endCh)
	go func(P *kg.LocalParty) {
		if err := P.Start(); err != nil {
			errCh <- err
		}
	}(party)

	for {
		select {
		case inmsg := <-in:
			//TODO erros msg from others,stop keygen
			go common.MsginHandle(party, inmsg, errCh)

		case msg := <-outCh:
			out <- msg

		case errMsg := <-errCh:
			//TODO send error msg
			return nil, errMsg

		case save := <-endCh:
			return convert(&save), nil
		}

	}

}

// Local generate partySaveDate Centrally
// Threshold = Participants - (Participants - 1)/3 -1
func Local(Threshold int, Participants int) ([]*LocalPartySaveData, error) {
	tss.SetCurve(edwards.Edwards())
	pIDs, err := GenPartyIDs(Participants)
	if err != nil {
		return nil, err
	}
	ids := pIDs.Keys()
	idsLen := len(ids)
	vc := make(vss.Vs, Threshold+1)
	Xi := make([]*big.Int, idsLen)
	allsv := make([]*LocalPartySaveData, idsLen)

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
	modQ := tssComm.ModInt(tss.EC().Params().N)

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

func convert(sv *kg.LocalPartySaveData) *LocalPartySaveData {
	svLocal := new(LocalPartySaveData)
	svLocal.Xi = sv.Xi
	svLocal.ShareID = sv.ShareID
	svLocal.Ks = sv.Ks
	svLocal.BigXj = sv.BigXj
	svLocal.EDDSAPub = sv.EDDSAPub
	return svLocal
}
