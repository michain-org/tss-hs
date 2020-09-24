package keygen

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/pkg/errors"
)

type (
	LocalSecrets struct {
		// secret fields (not shared, but stored locally)
		Xi, ShareID *big.Int // xi, kj
	}

	// Everything in LocalPartySaveData is saved locally to user's HD when done
	LocalPartySaveData struct {
		LocalSecrets

		// original indexes (ki in signing preparation phase)
		Ks []*big.Int

		// public keys (Xj = uj*G for each Pj)
		BigXj []*crypto.ECPoint // Xj

		// used for test assertions (may be discarded)
		EDDSAPub *crypto.ECPoint // y
	}
)



//StorageSavedata save a savedata which is used for signing after keygen into file
func StorageSavedata(sv *LocalPartySaveData, path string) error {
	perm := os.FileMode(0600)
	b, err := json.Marshal(sv)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, b, perm); err != nil {
		return err
	}
	return nil
}

//LoadSavedata load a savedata from file
func LoadSavedata(path string) (*LocalPartySaveData, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sv LocalPartySaveData
	if err := json.Unmarshal(b, &sv); err != nil {
		return nil, err
	}
	return &sv, nil
}

//CheckSaves check for errors of savedate
func CheckSaves(svs []*LocalPartySaveData, threshold int) (bool, error) {
	if threshold > len(svs) {
		return false, errors.New("number of threshold should less than savedates")
	}
	EDDSAPub := svs[0].EDDSAPub
	// svs = svs[:threshold+1]
	shares := make(vss.Shares, len(svs))
	for i, s := range svs {
		shares[i] = &vss.Share{Threshold: threshold, ID: s.ShareID, Share: s.Xi}
	}
	secret, err := shares.ReConstruct()
	if err != nil {
		return false, errors.Wrap(err, "ReConstruct failed")
	}
	p := crypto.ScalarBaseMult(tss.EC(), secret)
	return p.Equals(EDDSAPub), nil
}
