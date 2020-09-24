package msgHandle

import (
	"errors"
	"tss-hs/signing"
)

type (
	TssHandler struct {
		parties []signing.Party
		route   Router
	}

	Router interface {
		SetChan(inchan chan *signing.Message, outchan chan *signing.Message)
		HotRoute(mode int)
	}
)

const (
	ValidatorMode = iota
	ProposerMode
)

var (
	Handler TssHandler
	path    = ""
)

// func init() {
// 	Handler.parties = signing.Init(path)
// }

func validMode(PartyMod int) error {
	if !(PartyMod == ValidatorMode || PartyMod == ProposerMode) {
		return newHandleError(PartyMod, task, errors.New("Only ValidatorMode & ProposerMode can be used"))
	}
	return nil
}

func (t *TssHandler) SetRouter(r Router) {
	t.route = r
}

func (t *TssHandler) Start(PartyMod int, msg []byte) (*[32]byte, error) {
	if t.route == nil {
		return nil, newHandleError(PartyMod, task, errors.New("Should set a non-nil Router"))
	}
	if err := validMode(PartyMod); err != nil {
		return nil, err
	}
	party := t.parties[PartyMod]
	inChan, outChan := party.GetChannel()
	t.route.SetChan(inChan, outChan)
	go t.route.HotRoute(PartyMod)
	sig, err := party.Negotiation(inChan, outChan, msg)
	if err != nil {
		return nil, newHandleError(PartyMod, task, err)
	}
	return sig, nil
}

func (t *TssHandler) Verify(PartyMod int, sumSig *[32]byte) (bool, error) {
	if err := validMode(PartyMod); err != nil {
		return false, err
	}
	return t.parties[PartyMod].Verify(sumSig), nil
}

func SumSis(sis []*[32]byte) (sumS *[32]byte) {
	return signing.SumSis(sis)
}
