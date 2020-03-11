package exported

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error

	// Calculates the amount of coins that can be sent to other accounts given
	// the current time.
	SpendableCoins(blockTime time.Time) sdk.Coins

	// Ensure that account implements stringer
	String() string

	GetMultiSig() MultiSig
	SetMultiSig(signers MultiSig) error
	IsMultiSig() bool
	IsSigner(signer sdk.AccAddress) bool
}

type MultiSig interface {
	AddPendingTx(tx sdk.Tx, sender sdk.AccAddress) (PendingTx, error)
	GetPendingTx(txID uint64) PendingTx
	NumOfPendingTxa() int
	ContainTx(txID uint64) bool
	IsOwner(owner sdk.AccAddress) bool
	GetSigners() []sdk.AccAddress
	GetThreshold() int
	GetOwner() sdk.AccAddress
	RemoveTx(txID uint64) bool
	GetCounter() uint64
	CheckTx(tx sdk.Tx) (uint64, bool)
	UpdateSigners(signers []sdk.AccAddress, threshold int) error
	UpdateOwner(owner sdk.AccAddress) error
	IsMetric(txID uint64) bool
	AddSignature(txID uint64, signer sdk.AccAddress, sig []byte) (PendingTx, error)
}

type PendingTx interface {
	GetID() uint64
	GetTx() sdk.Tx
	GetSender() sdk.AccAddress
	IsSignedBy(signer sdk.AccAddress) bool
}

// GenesisAccounts defines a slice of GenesisAccount objects
type GenesisAccounts []GenesisAccount

// Contains returns true if the given address exists in a slice of GenesisAccount
// objects.
func (ga GenesisAccounts) Contains(addr sdk.Address) bool {
	for _, acc := range ga {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
}

// GenesisAccount defines a genesis account that embeds an Account with validation capabilities.
type GenesisAccount interface {
	Account
	Validate() error
}
