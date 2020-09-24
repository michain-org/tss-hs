package msgHandle

import (
	"fmt"
	"sync"
	"testing"

	"github.com/michain-org/tss-hs/keygen"
	sig "github.com/michain-org/tss-hs/signing"
	"github.com/stretchr/testify/assert"
)

type (
	safeInchan struct {
		mux    *sync.Mutex
		inchan chan *sig.Message
	}

	Vroute struct {
		fromPros, toPros, outchan, inchan chan *sig.Message
	}
	Proute struct {
		*safeInchan
		outchan          chan *sig.Message
		fromValt, toValt []chan *sig.Message
	}
)

var (
	PathPrefix   = "test/Sd_"
	participant  = 10
	threshold    = 6
	ConsensusMsg = []byte("hello wolid")

	globalProInchan  chan *sig.Message
	globalProOutchan chan *sig.Message
	globalValInchan  []chan *sig.Message
	globalValOutchan []chan *sig.Message
)

func start_1(parties []sig.Party, PartyMod int) {
	if !(PartyMod == ValidatorMode || PartyMod == ProposerMode) {
		return
	}
	party := parties[PartyMod]
	inChan, outChan := party.GetChannel()
	if PartyMod == ValidatorMode {
		globalValInchan = append(globalValInchan, inChan)
		globalValOutchan = append(globalValOutchan, outChan)
	} else {
		globalProInchan = inChan
		globalProOutchan = outChan
	}
}

func start_2(i int, parties []sig.Party, PartyMod int, msg []byte, router Router) *[32]byte {
	if !(PartyMod == ValidatorMode || PartyMod == ProposerMode) {
		return nil
	}
	party := parties[PartyMod]
	var inChan, outChan chan *sig.Message
	if PartyMod == ProposerMode {
		inChan = globalProInchan
		outChan = globalProOutchan
	} else {
		inChan = globalValInchan[i-1]
		outChan = globalValOutchan[i-1]
	}
	router.SetChan(inChan, outChan)
	go router.HotRoute(PartyMod)
	sig, _ := party.Negotiation(inChan, outChan, msg)
	return sig
}

func (p *Proute) SetChan(inchan chan *sig.Message, outchan chan *sig.Message) {
	p.safeInchan = &safeInchan{
		mux:    new(sync.Mutex),
		inchan: inchan,
	}
	p.outchan = outchan
}

func (v *Vroute) SetChan(inchan chan *sig.Message, outchan chan *sig.Message) {
	v.inchan = inchan
	v.outchan = outchan
}

//Testing router tests consensus only one times ,do not need to worry about kill HotRoute gorutine

func (p *Proute) HotRoute(mode int) {
	for i := range p.fromValt {
		go func(i int, p *Proute) {
			msg := <-p.fromValt[i]
			// fmt.Println("4.Proute receive a message from Vroute")
			p.mux.Lock()
			p.inchan <- msg
			// fmt.Println("5.Proute send a message to Proposer's inChannel")
			p.mux.Unlock()
		}(i, p)
	}
	msg := <-p.outchan

	// fmt.Println("8.Proute receive a message from  proposer's outchannel")
	for _, v := range p.toValt {
		go func(msg *sig.Message, c chan *sig.Message) {
			c <- msg
		}(msg, v)
	}
}

func (v *Vroute) HotRoute(mode int) {

	go func(v *Vroute) {
		msg := <-v.outchan
		// fmt.Println("2.Vroute receive a message from validator's outChannel!!!")
		v.toPros <- msg
		// fmt.Println("3.Vroute send a R message to Proute's fromValt Channel ")
	}(v)

}

func TestMsgHandle(t *testing.T) {
	//keygen locally and store savedata into files
	svs, err := keygen.Local(threshold, participant)
	if err != nil {
		t.Fatal(err, "keygen locally failed ")
	}
	fs := make([]string, participant)
	for i := range svs {
		f := fmt.Sprintf("%s%d.json", PathPrefix, i)
		fs[i] = f
		if err := keygen.StorageSavedata(svs[i], f); err != nil {
			t.Fatal(err)
		}

	}

	nodes := make([][]sig.Party, participant)
	sigChan := make(chan *[32]byte, participant)
	var wg sync.WaitGroup
	for i := range nodes {
		wg.Add(1)
		mode := ValidatorMode
		if i == 0 {
			mode = ProposerMode
		}
		path := fs[i]
		nodes[i] = sig.Init4test(path)
		start_1(nodes[i], mode)

	}

	//set router to listening

	for i := range nodes {
		var route Router
		mode := ValidatorMode
		if i == 0 {
			mode = ProposerMode
		}
		if mode == ProposerMode {
			route = &Proute{
				fromValt: globalValOutchan[:],
				toValt:   globalValInchan[:],
			}
		} else {
			route = &Vroute{
				fromPros: globalProOutchan,
				toPros:   globalProInchan,
			}
		}
		go func(i int, node []sig.Party, mode int, msg []byte, route Router, siga chan *[32]byte) {
			siga <- start_2(i, node, mode, msg, route)
			defer wg.Done()
		}(i, nodes[i], mode, ConsensusMsg, route, sigChan)
	}
	wg.Wait()
	close(sigChan)
	var sigs []*[32]byte
	for s := range sigChan {
		if s != nil {
			sigs = append(sigs, s)
		}
	}
	for i := range sigs {
		t.Log(sigs[i])
	}
	sumSig := SumSis(sigs)
	t.Log(sumSig)
	assert.True(t, nodes[0][1].Verify(sumSig))
}
