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
	networkID uint32
	subnetID  ids.ID
	clock     timer.Clock
	keyChain  *spdagvm.KeyChain // Mapping from public address to the SigningKeys
	utxoSet   *UtxoSet          // Mapping from utxoIDs to Utxos
	balance   uint64
	txFee     uint64
}

// NewWallet returns a new Wallet
func NewWallet(networkID uint32, subnetID ids.ID, txFee uint64) *Wallet {
	return &Wallet{
		networkID: networkID,
		subnetID:  subnetID,
		keyChain:  &spdagvm.KeyChain{},
		utxoSet:   &UtxoSet{},
		txFee:     txFee,
	}
}

// GetAddress returns one of the addresses this wallet manages. If no address
// exists, one will be created.
func (w *Wallet) GetAddress() ids.ShortID {
	if w.keyChain.Addrs.Len() == 0 {
		return w.CreateAddress()
	}
	return w.keyChain.Addrs.CappedList(1)[0]
}

// CreateAddress returns a new address.
// It also saves the address and the private key that controls it
// so the address can be used later
func (w *Wallet) CreateAddress() ids.ShortID {
	privKey, _ := w.keyChain.New()
	return privKey.PublicKey().Address()
}

// ImportKey imports a private key into this wallet
func (w *Wallet) ImportKey(sk *crypto.PrivateKeySECP256K1R) { w.keyChain.Add(sk) }

// AddUtxo adds a new UTXO to this wallet if this wallet may spend it
// The UTXO's output must be an OutputPayment
func (w *Wallet) AddUtxo(utxo *spdagvm.UTXO) {
	out, ok := utxo.Out().(*spdagvm.OutputPayment)
	if !ok {
		return
	}

	if _, _, err := w.keyChain.Spend(utxo, math.MaxUint64); err == nil {
		w.utxoSet.Put(utxo)
		w.balance += out.Amount()
	}
}

// Balance returns the amount of the assets in this wallet
func (w *Wallet) Balance() uint64 { return w.balance }

// CreateTx sends some amount to the destination addresses
func (w *Wallet) CreateTx(amount uint64, locktime uint64, threshold uint32, dests []ids.ShortID) *spdagvm.Tx {
	//ins, outs, signers, _ := w.txPrepare(amount, locktime, threshold, dests)
	builder := spdagvm.Builder{
		NetworkID: w.networkID,
		SubnetID:  w.subnetID,
	}
	currentTime := w.clock.Unix()

	// Send any change to an address this wallet controls
	changeAddr := ids.ShortID{}
	if w.keyChain.Addrs.Len() < 1000 {
		changeAddr = w.CreateAddress()
	} else {
		changeAddr = w.GetAddress()
	}

	// Build the transaction
	tx, err := builder.NewTxFromUTXOs(w.keyChain, w.utxoSet.Utxos, amount, w.txFee, locktime, dests, changeAddr, currentTime)
	if err != nil {
		panic(err)
	}

	return tx
}

// SpendTx takes a tx, removes its utxos, and adds the inputs
func (w *Wallet) SpendTx(tx *spdagvm.Tx) {
	for _, in := range tx.Ins() {
		if in, ok := in.(*spdagvm.InputPayment); ok {
			utxoID := in.InputID()
			w.RemoveUtxo(utxoID)
			w.balance -= in.Amount() // Deduct from [w.balance] the amount sent
		}
	}

	for _, out := range tx.UTXOs() {
		w.AddUtxo(out)
	}
}

// GetNetworkID returns the networkID for the wallet
func (w *Wallet) GetNetworkID() uint32 {
	return w.networkID
}

// GetSubnetID returns the blockchainID for the wallet
func (w *Wallet) GetSubnetID() ids.ID {
	return w.subnetID
}

/*
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
}*/

func (w Wallet) String() string {
	return fmt.Sprintf(
		"KeyChain:\n"+
			"%s\n"+
			"UtxoSet:\n"+
			"%s",
		w.keyChain.PrefixedString("    "),
		w.utxoSet.string("    "))
}

// ADDED for Avash

// Addresses gets all the addresses in the systme
func (w *Wallet) Addresses() []string {
	addrs := w.keyChain.Addresses().List()
	results := []string{}
	for _, a := range addrs {
		results = append(results, a.String())
	}
	return results
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

// Wallets ...
var Wallets map[string]*Wallet

func init() {
	Wallets = map[string]*Wallet{}
}
