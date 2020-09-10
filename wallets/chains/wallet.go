package chainwallet

import (
	"fmt"

	"github.com/ava-labs/avalanche-go/ids"
	"github.com/ava-labs/avalanche-go/modules/chains/payments"
	"github.com/ava-labs/avalanche-go/utils/crypto"
)

// Wallet is a holder for keys and UTXOs.
type Wallet struct {
	keyChain   *payments.KeyChain            // Mapping from public address to the SigningKeys
	accountSet map[[20]byte]payments.Account // Mapping from addresses to accounts
	balance    uint64
}

// NewWallet ...
func NewWallet() Wallet {
	return Wallet{
		keyChain:   &payments.KeyChain{},
		accountSet: make(map[[20]byte]payments.Account),
	}
}

// CreateAddress returns a brand new address! Ready to receive funds!
func (w *Wallet) CreateAddress() ids.ShortID { return w.keyChain.New().PublicKey().Address() }

// ImportKey imports a private key into this wallet
func (w *Wallet) ImportKey(sk *crypto.PrivateKeySECP256K1R) { w.keyChain.Add(sk) }

// AddAccount adds a new account to this wallet, if this wallet can spend it.
func (w *Wallet) AddAccount(account payments.Account) {
	if account.Balance() > 0 {
		// if _, _, err := w.keyChain.Spend(account, account.Balance(), ids.NewShortID([20]byte{})); account.Balance() > 0 && err == nil {
		w.accountSet[account.ID().Key()] = account
		w.balance += account.Balance()
	}
}

// Balance returns the amount of the assets in this wallet
func (w *Wallet) Balance() uint64 { return w.balance }

// Send sends some amount to a new address
func (w *Wallet) Send(newAccount bool) (*payments.Tx, payments.Account, payments.Account) {
	destination := payments.Account{}
	if len(w.accountSet) <= 1 || newAccount {
		destAddr := w.CreateAddress()

		builder := payments.Builder{}
		destination = builder.NewAccount(destAddr, 0, 0)
	} else {
		for accoutID, account := range w.accountSet {
			delete(w.accountSet, accoutID)
			w.balance -= account.Balance()
			destination = account
			break
		}
	}
	for accountID, account := range w.accountSet {
		if key, exists := w.keyChain.Get(account.ID()); exists {
			amount := uint64(1)
			if tx, sendAccount, err := account.CreateTx(amount, destination.ID(), key); err == nil {
				delete(w.accountSet, accountID)
				w.balance -= account.Balance()

				builder := payments.Builder{}
				return tx, sendAccount, builder.NewAccount(destination.ID(), destination.Nonce(), destination.Balance()+amount)
			}
		}
	}

	return nil, payments.Account{}, payments.Account{} // Insufficient funds
}

func (w Wallet) String() string {
	return fmt.Sprintf(
		"KeyChain:\n"+
			"%s",
		w.keyChain.PrefixedString("    "))
}
