package dagwallet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/gecko/utils/formatting"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/modules/dags/ava"
)

// UtxoSetResult result for a UtxoSet
type UtxoSetResult struct {
	// This can be used to iterate over. However, it should not be modified externally.
	UtxoMap map[string]int
	Utxos   []UTXO
}

// Put ...
func (us *UtxoSetResult) Put(utxo *ava.UTXO) {
	if us.UtxoMap == nil {
		us.UtxoMap = make(map[string]int)
	}
	if _, ok := us.UtxoMap[utxo.ID().String()]; !ok {
		us.UtxoMap[utxo.ID().String()] = len(us.Utxos)
		newUTXO := UTXO{}
		newUTXO.copyFromAvaUTXO(utxo)
		us.Utxos = append(us.Utxos, newUTXO)
	}
}

// UTXO just used to produce a set difference
type UTXO struct {
	SourceID    string `json:"sourceID"`
	SourceIndex uint32 `json:"sourceIndex"`
	ID          string `json:"id"`

	Out OutputPayment `json:"out"`

	Bytes string `json:"bytes"`
}

// OutputPayment is a structure which allows serialization of UTXO outputs
type OutputPayment struct {
	Amount   uint64 `json:"amount"`
	Locktime uint64 `json:"locktime"`

	Threshold uint32   `json:"threshold"`
	Addresses []string `json:"addresses"`
}

func (u *UTXO) copyFromAvaUTXO(utxo *ava.UTXO) {
	sid, sidx := utxo.Source()
	u.SourceID = sid.String()
	u.SourceIndex = sidx
	id := utxo.ID().String()
	u.ID = id
	b := utxo.Bytes()
	fb := formatting.FormatBytes{}
	fb.Bytes = b
	bstr := fb.String()
	u.Bytes = bstr

	out := utxo.Out()
	addrs := out.(*ava.OutputPayment).Addresses()
	newAddrs := []string{}

	for _, a := range addrs {
		newAddrs = append(newAddrs, a.String())
	}

	newOut := OutputPayment{
		Amount:    out.(*ava.OutputPayment).Amount(),
		Locktime:  out.(*ava.OutputPayment).Locktime(),
		Threshold: out.(*ava.OutputPayment).Threshold(),
		Addresses: newAddrs,
	}
	u.Out = newOut
}

// UtxoSet ...
type UtxoSet struct {
	// This can be used to iterate over. However, it should not be modified externally.
	utxoMap map[[32]byte]int
	Utxos   []*ava.UTXO
}

// Put ...
func (us *UtxoSet) Put(utxo *ava.UTXO) {
	if us.utxoMap == nil {
		us.utxoMap = map[[32]byte]int{}
	}
	if _, ok := us.utxoMap[utxo.ID().Key()]; !ok {
		us.utxoMap[utxo.ID().Key()] = len(us.Utxos)
		us.Utxos = append(us.Utxos, utxo)
	}
}

func (us *UtxoSet) JSON() ([]byte, error) {
	result := UtxoSetResult{
		UtxoMap: map[string]int{},
		Utxos:   []UTXO{},
	}
	for _, v := range us.utxoMap {
		result.Put(us.Utxos[v])
	}
	resultJSON, err := json.MarshalIndent(result, "", "    ")
	return resultJSON, err
}

// SetDiff takes two UtxoSets and returns a set difference result
func (us *UtxoSet) SetDiff(us2 UtxoSet) UtxoSetResult {

	unionSet := UtxoSet{
		utxoMap: map[[32]byte]int{},
		Utxos:   []*ava.UTXO{},
	}

	intersectSet := UtxoSet{
		utxoMap: map[[32]byte]int{},
		Utxos:   []*ava.UTXO{},
	}

	resultSet := UtxoSetResult{
		UtxoMap: map[string]int{},
		Utxos:   []UTXO{},
	}
	for k, v := range us2.utxoMap {
		unionSet.Put(us2.Utxos[v])
		if _, ok := us.utxoMap[k]; ok && v < len(us2.Utxos) {
			intersectSet.Put(us2.Utxos[v])
		}
	}
	for _, v := range us.utxoMap {
		if v < len(us2.Utxos) {
			unionSet.Put(us2.Utxos[v])
		}
	}

	for k, v := range unionSet.utxoMap {
		if _, ok := intersectSet.utxoMap[k]; !ok {
			resultSet.Put(unionSet.Utxos[v])
		}
	}

	return resultSet
}

// Get ...
func (us *UtxoSet) Get(id ids.ID) *ava.UTXO {
	if us.utxoMap == nil {
		return nil
	}
	if i, ok := us.utxoMap[id.Key()]; ok {
		utxo := us.Utxos[i]
		return utxo
	}
	return nil
}

// Remove ...
func (us *UtxoSet) Remove(id ids.ID) *ava.UTXO {
	i, ok := us.utxoMap[id.Key()]
	if !ok {
		return nil
	}
	utxoI := us.Utxos[i]

	j := len(us.Utxos) - 1
	utxoJ := us.Utxos[j]

	us.Utxos[i] = us.Utxos[j]
	us.Utxos = us.Utxos[:j]

	us.utxoMap[utxoJ.ID().Key()] = i
	delete(us.utxoMap, utxoI.ID().Key())

	return utxoI
}

func (us *UtxoSet) string(prefix string) string {
	s := strings.Builder{}

	for i, utxo := range us.Utxos {
		out := utxo.Out().(*ava.OutputPayment)
		sourceID, sourceIndex := utxo.Source()

		s.WriteString(fmt.Sprintf("%sUtxo[%d]:"+
			"\n%s    InputID: %s"+
			"\n%s    InputIndex: %d"+
			"\n%s    Locktime: %d"+
			"\n%s    Amount: %d\n",
			prefix, i,
			prefix, sourceID,
			prefix, sourceIndex,
			prefix, out.Locktime(),
			prefix, out.Amount()))
	}

	return strings.TrimSuffix(s.String(), "\n")
}

func (us *UtxoSet) String() string {
	return us.string("")
}
