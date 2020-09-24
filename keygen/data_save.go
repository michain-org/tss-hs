package keygen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/vss"
	kg "github.com/binance-chain/tss-lib/eddsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/pkg/errors"
)

//StorageSavedata save a savedata which is used for signing after keygen into file
func StorageSavedata(sv *kg.LocalPartySaveData, path string) error {
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
func LoadSavedata(path string) (*kg.LocalPartySaveData, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sv kg.LocalPartySaveData
	if err := json.Unmarshal(b, &sv); err != nil {
		return nil, err
	}
	return &sv, nil
}

//CheckSaves check for errors of savedate
func CheckSaves(svs []*kg.LocalPartySaveData, threshold int) (bool, error) {
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
