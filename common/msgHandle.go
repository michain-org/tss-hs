package common

import (
	"github.com/binance-chain/tss-lib/tss"
)

//MsginHandle handles msg from others
func MsginHandle(party tss.Party, msg tss.Message, errCh chan<- *tss.Error) {
	if party.PartyID() == msg.GetFrom() {
		return
	}
	bz, _, err := msg.WireBytes()
	if err != nil {
		errCh <- party.WrapError(err)
		return
	}
	pMsg, err := tss.ParseWireMessage(bz, msg.GetFrom(), msg.IsBroadcast())
	if err != nil {
		errCh <- party.WrapError(err)
		return
	}
	// fmt.Println("Update start")
	if _, err := party.Update(pMsg); err != nil {
		// fmt.Println("handle the in msgs err:", err)
		errCh <- err
	}
}
