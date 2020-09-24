package keygen

import (
	"fmt"
	"sync"
	"testing"

	kg "github.com/binance-chain/tss-lib/eddsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/stretchr/testify/assert"
)

const (
	TestParticipants = 10
	TestThreshold    = 6
)

func TestLocal(t *testing.T) {
	var Threshold int = TestThreshold
	var Participants int = TestParticipants
	saves, err := Local(Threshold, Participants)
	if err != nil {
		fmt.Printf("Local gen err %s", err)
		return
	}

	for _, save := range saves {
		fmt.Println(save)
	}

	chk, _ := CheckSaves(saves, Threshold)
	assert.True(t, chk)
}

func TestPidStorage(t *testing.T) {
	path := "./pid.json"
	Participants := 3
	pids, err := GenPartyIDs(Participants)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v", pids)
	StoragePids(pids, 0, 2, path)
	pid, err := LoadPids(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pid)
}

func TestDecentralized(t *testing.T) {
	type testSvs struct {
		Svs []*kg.LocalPartySaveData
		Mtx sync.Mutex
	}
	saves := new(testSvs)
	pids, err := GenPartyIDs(TestParticipants)
	if err != nil {
		fmt.Println(err)
		return
	}
	allsds := make([]chan tss.Message, TestParticipants)
	paras := make([]*Paramer, TestParticipants)
	for i := 0; i < TestParticipants; i++ {
		paras[i] = &Paramer{Pids: pids, Index: i, Threshold: TestThreshold}
		allsds[i] = make(chan tss.Message, TestParticipants)
	}
	allrecvs := make(chan tss.Message, TestParticipants)
	allsvs := make(chan kg.LocalPartySaveData, TestParticipants)
	for i := 0; i < TestParticipants; i++ {
		go func(para *Paramer,
			in <-chan tss.Message,
			out chan<- tss.Message) {
			sv, err := Decentralized(para, in, out)
			if err != nil {
				return
			}
			allsvs <- *sv
		}(paras[i], allsds[i], allrecvs)
	}

	for {
		select {
		case msg := <-allrecvs:
			dest := msg.GetTo()
			if dest == nil { // broadcast!
				for i := 0; i < TestParticipants; i++ {
					if i == msg.GetFrom().Index {
						continue
					}
					allsds[i] <- msg
				}
			} else { // point-to-point!
				if dest[0].Index == msg.GetFrom().Index {
					t.Fatalf("party %d tried to send a message to itself (%d)", dest[0].Index, msg.GetFrom().Index)
					return
				}
				allsds[dest[0].Index] <- msg
			}

		case sv := <-allsvs:
			saves.Mtx.Lock()
			saves.Svs = append(saves.Svs, &sv)
			if len(saves.Svs) == TestThreshold+1 {
				//try cases where length of saves.Svs equals to TestThreshold
				chk, _ := CheckSaves(saves.Svs, TestThreshold)
				assert.True(t, chk)
				// for i, s := range saves.Svs {
				// 	if err := StorageSavedata(s, fmt.Sprintf("saves%d", i)); err != nil {
				// 		fmt.Println(err)
				// 	}
				// }
				return
			}
			saves.Mtx.Unlock()
		}
	}

}
