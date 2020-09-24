package signing

import (
	sync "sync"

	"github.com/michain-org/tss-hs/keygen"

	"github.com/binance-chain/tss-lib/crypto"
	kg "github.com/michain-org/tss-hs/keygen"
	"github.com/pkg/errors"
)

type (
	Party interface {
		Negotiation(inchan <-chan *Message, outchan chan<- *Message, conMsg []byte) (*[32]byte, error)
		Verify(*[32]byte) bool
		GetChannel() (inchan chan *Message, outchan chan *Message)
		// SumSis(sis []*[32]byte) *[32]byte
	}

	Validator struct {
		*Base
	}

	waitRi struct {
		mux  *sync.Mutex
		rl   []*crypto.ECPoint
		idxs []int
	}

	Proposer struct {
		*Base
		threshold int
		wr        *waitRi
	}
)

var (
	validator *Validator
	proposer  *Proposer
)

func load(path string) (*kg.LocalPartySaveData, error) {
	svdata, err := keygen.LoadSavedata(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Load savedate failed from %s", path)
	}
	if svdata.Xi == nil ||
		len(svdata.Ks) == 0 ||
		svdata.ShareID == nil ||
		svdata.EDDSAPub == nil {
		return nil, errors.New("Bad format savedata")
	}
	return svdata, nil
}

func Init(path string) []Party {
	sv, err := load(path)
	if err != nil {
		panic(err)
	}
	parties := len(sv.Ks)
	tolerance := int((parties - 1) / 3)
	threshold := parties - tolerance - 1
	b := newBase(sv)
	validator = &Validator{
		Base: b,
	}
	proposer = &Proposer{
		Base:      b,
		threshold: threshold,
	}

	return []Party{validator, proposer}
}

func Init4test(path string) []Party {
	sv, err := load(path)
	if err != nil {
		panic(err)
	}
	parties := len(sv.Ks)
	tolerance := int((parties - 1) / 3)
	threshold := parties - tolerance - 1
	b := newBase(sv)
	validator := &Validator{
		Base: b,
	}
	proposer := &Proposer{
		Base:      b,
		threshold: threshold,
	}

	return []Party{validator, proposer}
}

func (v *Validator) must(idxs []int) bool {
	for _, val := range idxs {
		if val == v.index {
			return true
		}
	}
	return false
}

func (v *Validator) Negotiation(inchan <-chan *Message,
	outchan chan<- *Message,
	conMsg []byte) (*[32]byte, error) {
	if v == nil {
		return nil, errors.New("Init first")
	}
	v.newRound(conMsg)

	go func(outchan chan<- *Message, v *Validator) {
		outchan <- NewSignMessagePoint(v.temp.R, v.index)
		// fmt.Println("1.Validator send a R mesaage to it's outChannel!!!")
	}(outchan, v)

	var Rs []*crypto.ECPoint
	var idxs []int
	msgin := <-inchan
	if smpl := msgin.GetSmpl(); smpl != nil {
		if !smpl.Validate() {
			return nil, errors.New("Get a bad smpl message from Proposer")
		}
		Rs, idxs = smpl.Unmashal()
	}
	if !v.must(idxs) {
		v.calcSi(Rs)
		return nil, nil
	}
	if err := v.calcWi(idxs); err != nil {
		return nil, err
	}
	v.calcSi(Rs)
	var tempsi [32]byte
	copy(tempsi[:], v.temp.si[:])
	return &tempsi, nil
}

func (p *Proposer) Negotiation(inchan <-chan *Message,
	outchan chan<- *Message,
	conMsg []byte) (*[32]byte, error) {
	if p == nil {
		return nil, errors.New("Init first")
	}

	p.newRound(conMsg)
	p.wr = &waitRi{
		mux:  new(sync.Mutex),
		rl:   make([]*crypto.ECPoint, 0, p.threshold+1),
		idxs: make([]int, 0, p.threshold+1),
	}
	p.wr.rl = append(p.wr.rl, p.temp.R)
	p.wr.idxs = append(p.wr.idxs, p.index)
	signal := make(chan bool)

CollectRiS:
	for {
		select {
		case msgri := <-inchan:
			// fmt.Println("6.Proposer's inChannel receive a message from Proute")
			if smp := msgri.GetSmp(); smp != nil {
				if !smp.Validate() {
					return nil, errors.New("Get a bad smp message from Proposer")
				}
				r, index := smp.Unmashal()
				go func(p *Proposer, point *crypto.ECPoint, index int, sign chan bool) {
					p.wr.mux.Lock()
					p.wr.rl = append(p.wr.rl, point)
					p.wr.idxs = append(p.wr.idxs, index)
					if len(p.wr.rl) == p.threshold+1 {
						sign <- true
					}
					defer p.wr.mux.Unlock()
				}(p, r, index, signal)
			}

		case _ = <-signal:
			p.wr.mux.Lock()
			break CollectRiS
		}
	}
	p.wr.rl = p.wr.rl[:p.threshold+1]
	p.wr.idxs = p.wr.idxs[:p.threshold+1]
	go func(outchan chan<- *Message, sumR []*crypto.ECPoint, idxs []int) {
		outchan <- NewSignMessagePointList(sumR, idxs)
		// fmt.Println("7.Proposer send a message to it's outchannel ")
	}(outchan, p.wr.rl, p.wr.idxs)
	if err := p.calcWi(p.wr.idxs); err != nil {
		return nil, err
	}
	p.calcSi(p.wr.rl)
	var tempsi [32]byte
	copy(tempsi[:], p.temp.si[:])
	return &tempsi, nil

}

func (p *Proposer) GetChannel() (inchan chan *Message, outchan chan *Message) {
	inchan = make(chan *Message, p.threshold+1)
	outchan = make(chan *Message)
	return
}

func (v *Validator) GetChannel() (inchan chan *Message, outchan chan *Message) {
	inchan = make(chan *Message)
	outchan = make(chan *Message)
	return
}
