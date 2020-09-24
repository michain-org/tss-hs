package signing

import (
	"math/big"
	"tss-hs/common"

	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/tss"
)

func (m *SignMessagePoint) Unmashal() (*crypto.ECPoint, int) {
	x := m.GetX()
	y := m.GetY()
	point, err := crypto.NewECPoint(tss.EC(), new(big.Int).SetBytes(x), new(big.Int).SetBytes(y))
	if err != nil {
		return nil, -1
	}
	index := m.GetIndex()
	return point, int(index)
}

func (m *SignMessagePoint) Validate() bool {
	if !(common.NonEmptyBytes(m.GetX()) &&
		common.NonEmptyBytes(m.GetY())) &&
		!(m.GetIndex() < 0) {
		return false
	}
	return true
}

func (m *SignMessagePointList) Unmashal() ([]*crypto.ECPoint, []int) {
	bs := m.GetLi()
	l := common.Bytes2ints(bs)
	ps := make([]*crypto.ECPoint, len(l)/2)
	var err error
	for i := 0; i < len(l); i += 2 {
		if ps[i/2], err = crypto.NewECPoint(tss.EC(), l[i], l[i+1]); err != nil {
			return nil, nil
		}
	}
	indexes := m.GetIndexes()
	//Validate
	idxs := make([]int, len(indexes))
	for i := range indexes {
		idxs[i] = int(indexes[i])
	}
	// fmt.Println(ps)
	return ps, idxs
}

func (m *SignMessagePointList) Validate() bool {
	bs := m.GetLi()
	l := common.Bytes2ints(bs)
	indexes := m.GetIndexes()
	if !(common.NonEmptyMulBytes(bs) &&
		len(l)%2 == 0 &&
		indexes != nil &&
		len(indexes) == len(l)/2) {
		return false
	}
	return true
}

func NewSignMessagePoint(p *crypto.ECPoint, index int) *Message {
	smp := &SignMessagePoint{
		X:     p.X().Bytes(),
		Y:     p.Y().Bytes(),
		Index: int32(index),
	}
	return &Message{
		Content: &Message_Smp{Smp: smp},
	}
}

func NewSignMessagePointList(pl []*crypto.ECPoint, indexes []int) *Message {
	bigs := make([]*big.Int, len(pl)*2)
	idxs := make([]int32, len(indexes))
	for i, e := range pl {
		index := i * 2
		bigs[index] = e.X()
		bigs[index+1] = e.Y()
		idxs[i] = int32(indexes[i])
	}
	smpl := &SignMessagePointList{
		Li:      common.Ints2bytes(bigs),
		Indexes: idxs,
	}
	return &Message{
		Content: &Message_Smpl{
			Smpl: smpl,
		},
	}
}
