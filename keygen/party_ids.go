package keygen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/binance-chain/tss-lib/common"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/pkg/errors"
)

type (
	//Paramer contains each party's storage state for init keygen
	Paramer struct {
		Pids      tss.SortedPartyIDs
		Index     int
		Threshold int
	}
)

//GenPartyIDs generate a specific unique ID for each party
func GenPartyIDs(Participants int, startAt ...int) (tss.SortedPartyIDs, error) {
	if Participants <= 0 {
		return nil, errors.New("Participants must be positive")
	}
	ids := make(tss.UnSortedPartyIDs, 0, Participants)
	key := common.MustGetRandomInt(256)
	frm := 0
	i := 0 // default `i`
	if len(startAt) > 0 {
		frm = startAt[0]
		i = startAt[0]
	}
	for ; i < Participants+frm; i++ {
		ids = append(ids, &tss.PartyID{
			MessageWrapper_PartyID: &tss.MessageWrapper_PartyID{
				Id:      fmt.Sprintf("%d", i+1),
				Moniker: fmt.Sprintf("P[%d]", i+1),
				Key:     new(big.Int).Sub(key, big.NewInt(int64(Participants)-int64(i))).Bytes(),
			}, Index: i})
	}
	return tss.SortPartyIDs(ids, startAt...), nil
}

//LoadPids load paryID
func LoadPids(path string) (*Paramer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var paramer *Paramer
	if err := json.Unmarshal(b, paramer); err != nil {
		return nil, err
	}
	return paramer, nil
}

//StoragePids storage partyID into file
func StoragePids(pids tss.SortedPartyIDs, index int, threshold int, path string) error {
	perm := os.FileMode(0600)
	paramer := &Paramer{Pids: pids, Index: index, Threshold: threshold}
	b, err := json.Marshal(paramer)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, b, perm); err != nil {
		return err
	}
	return nil
}
