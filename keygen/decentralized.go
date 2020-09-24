package keygen

import (
	"tss-hs/common"

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
) (*kg.LocalPartySaveData, error) {

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
			return &save, nil
		}

	}

}
