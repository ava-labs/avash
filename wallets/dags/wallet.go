package dagwallet

import (
	"fmt"
	"math"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/utils/timer"
	"github.com/ava-labs/gecko/vms/spdagvm"
)

// Wallet is a holder for keys and UTXOs for the Ava DAG.
type Wallet struct {
	clock    timer.Clock
	keyChain *spdagvm.KeyChain // Mapping from public address to the SigningKeys
	utxoSet  *UtxoSet          // Mapping from utxoIDs to Utxos
	balance  uint64
}

// NewWallet ...
func NewWallet() *Wallet {
	return &Wallet{
		keyChain: &spdagvm.KeyChain{},
		utxoSet:  &UtxoSet{},
	}
}

// CreateAddress returns a brand new address! Ready to receive funds!
func (w *Wallet) CreateAddress() ids.ShortID {
	kc, _ := w.keyChain.New()
	return kc.PublicKey().Address()
}

// ImportKey imports a private key into this wallet
func (w *Wallet) ImportKey(sk *crypto.PrivateKeySECP256K1R) { w.keyChain.Add(sk) }

// Addresses gets all the addresses in the systme
func (w *Wallet) Addresses() []string {
	addrs := w.keyChain.Addresses().List()
	results := []string{}
	for _, a := range addrs {
		results = append(results, a.String())
	}
	return results
}

// AddUtxo adds a new utxo to this wallet, if this wallet can spend it.
func (w *Wallet) AddUtxo(utxo *spdagvm.UTXO) {
	outP, ok := utxo.Out().(*spdagvm.OutputPayment)
	if !ok {
		return
	}

	if _, _, err := w.keyChain.Spend(utxo, math.MaxUint64); err == nil {
		w.utxoSet.Put(utxo)
		w.balance += outP.Amount()
	}
}

// RemoveUtxo adds a new utxo to this wallet, if this wallet can spend it.
func (w *Wallet) RemoveUtxo(utxoID ids.ID) {
	utxo := w.utxoSet.Get(utxoID)
	if utxo != nil {
		outP := utxo.Out()
		w.balance -= outP.(*spdagvm.OutputPayment).Amount()
		w.utxoSet.Remove(utxoID)
	}
}

// GetUtxos returns a copy of the UTXO set
func (w *Wallet) GetUtxos() UtxoSet { return *w.utxoSet }

// Balance returns the amount of the assets in this wallet
func (w *Wallet) Balance() uint64 { return w.balance }

// CreateTx sends some amount to the destination addresses
func (w *Wallet) CreateTx(amount uint64, locktime uint64, threshold uint32, dests []ids.ShortID) *spdagvm.Tx {
	ins, outs, signers, _ := w.txPrepare(amount, locktime, threshold, dests)
	builder := spdagvm.Builder{}
	tx, _ := builder.NewTx(ins, outs, signers)

	return tx
}

// SpendTx takes a tx, removes its utxos, and adds the inputs
func (w *Wallet) SpendTx(tx *spdagvm.Tx) {
	for _, in := range tx.Ins() {
		utxoID := in.InputID()
		w.RemoveUtxo(utxoID)
	}

	for _, out := range tx.UTXOs() {
		w.AddUtxo(out)
	}
}

func (w *Wallet) txPrepare(amount uint64, locktime uint64, threshold uint32, dests []ids.ShortID) ([]spdagvm.Input, []spdagvm.Output, []*spdagvm.InputSigner, uint64) {
	ins := []spdagvm.Input{}
	signers := []*spdagvm.InputSigner{}

	utxoIDs := []ids.ID{}
	spent := uint64(0)
	time := w.clock.Unix()
	for i := 0; i < len(w.utxoSet.Utxos) && amount > spent; i++ {
		utxo := w.utxoSet.Utxos[i]
		if in, signer, err := w.keyChain.Spend(utxo, time); err == nil {
			ins = append(ins, in)
			signers = append(signers, signer)
			utxoIDs = append(utxoIDs, utxo.ID())

			amount := in.(*spdagvm.InputPayment).Amount()
			spent += amount
		}
	}

	if spent < amount {
		return nil, nil, nil, 0 // Insufficient funds
	}

	builder := spdagvm.Builder{}

	outs := []spdagvm.Output{
		builder.NewOutputPayment(amount, locktime, threshold, dests),
	}

	if spent > amount {
		outs = append(outs,
			builder.NewOutputPayment(spent-amount, 0, 1, []ids.ShortID{w.CreateAddress()}),
		)
	}

	return ins, outs, signers, spent
}

func (w Wallet) String() string {
	return fmt.Sprintf(
		"KeyChain:\n"+
			"%s\n"+
			"UtxoSet:\n"+
			"%s",
		w.keyChain.PrefixedString("    "),
		w.utxoSet.string("    "))
}

// Wallets ...
var Wallets map[string]*Wallet

func init() {
	Wallets = map[string]*Wallet{}
}
