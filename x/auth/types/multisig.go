package types

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

type PendingTx struct {
	ID       uint64           `json:"id" yaml:"id "`
	StdTx    StdTx            `json:"tx" yaml:"tx"`
	Sender   sdk.AccAddress   `json:"sender" yaml:"sender"`
	SignedBy []sdk.AccAddress `json:"signers" yaml:"signers"`
}

func (ptx PendingTx) GetID() uint64 {
	return ptx.ID
}
func (ptx PendingTx) GetTx() sdk.Tx {
	return ptx.StdTx
}
func (ptx PendingTx) GetSender() sdk.AccAddress {
	return ptx.Sender
}
func (ptx PendingTx) IsSignedBy(signer sdk.AccAddress) bool {
	for _, s := range ptx.SignedBy {
		if s.Equals(signer) {
			return true
		}
	}
	return false
}

type MultiSig struct {
	Owner      sdk.AccAddress   `json:"owner" yaml:"owner"`
	Threshold  int              `json:"threshold" yaml:"threshold"`
	Counter    uint64           `json:"counter" yaml:"counter"`
	Signers    []sdk.AccAddress `json:"signers" yaml:"signers"`
	PendingTxs []PendingTx      `json:"pendingTxs" yaml:"pendingTxs"`
}

func NewMultiSig(owner sdk.AccAddress, threshold int, signers []sdk.AccAddress) *MultiSig {
	return &MultiSig{
		Owner:     owner,
		Threshold: threshold,
		Counter:   0,
		Signers:   signers,
	}
}

// PendingTx
func (mg *MultiSig) AddPendingTx(tx sdk.Tx, sender sdk.AccAddress) (exported.PendingTx, error) {
	stdTx, ok := tx.(StdTx)
	if !ok {
		panic("Invalid StdTx for multisig")
	}
	var signedBy []sdk.AccAddress
	for i, sig := range stdTx.Signatures {
		addr := sdk.AccAddress(sig.PubKey.Address())
		signedBy = append(signedBy, addr)
		stdTx.Signatures[i].PubKey = nil
	}
	ptx := PendingTx{
		ID:       mg.Counter,
		StdTx:    stdTx,
		Sender:   sender,
		SignedBy: signedBy,
	}
	_, exist := mg.CheckTx(stdTx)
	if exist {
		return nil, fmt.Errorf("This transaction already exists")
	}

	mg.PendingTxs = append(mg.PendingTxs, ptx)
	mg.Counter++
	return ptx, nil
}

// Multisig functions
func (mg MultiSig) ContainTx(txID uint64) bool {
	for _, e := range mg.PendingTxs {
		if e.ID == txID {
			return true
		}
	}
	return false
}

func (mg MultiSig) GetSigners() []sdk.AccAddress {
	return mg.Signers
}

func (mg MultiSig) GetThreshold() int {
	return mg.Threshold
}

func (mg MultiSig) GetOwner() sdk.AccAddress {
	return mg.Owner
}

func (mg *MultiSig) RemoveTx(txID uint64) bool {
	for i, ptx := range mg.PendingTxs {
		if ptx.ID == txID {
			mg.PendingTxs = append(mg.PendingTxs[:i], mg.PendingTxs[i+1:]...)
			return true
		}
	}
	return false
}

func (mg MultiSig) GetPendingTx(txID uint64) exported.PendingTx {
	for _, ptx := range mg.PendingTxs {
		if ptx.ID == txID {
			return &ptx
		}
	}
	return nil
}

func (mg MultiSig) IsMetric(txID uint64) bool {
	for _, ptx := range mg.PendingTxs {
		if ptx.GetID() == txID {
			if len(ptx.StdTx.Signatures) == mg.Threshold {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (mg *MultiSig) GetCounter() uint64 {
	return mg.Counter
}

func (mg MultiSig) CheckTx(tx sdk.Tx) (uint64, bool) {
	for _, v := range mg.PendingTxs {
		stdTx, ok := tx.(StdTx)
		if !ok {
			panic("Invalid StdTx for multisig")
		}
		if reflect.DeepEqual(v.StdTx, stdTx) {
			return v.ID, true
		}
	}
	return 0, false
}

func (mg MultiSig) NumOfPendingTxa() int {
	return len(mg.PendingTxs)
}

func (mg *MultiSig) UpdateSigners(signers []sdk.AccAddress, threshold int) error {
	if len(mg.PendingTxs) > 0 {
		return fmt.Errorf("Unable to update multisig account because it has Pending Txs.")
	}
	if len(signers) < threshold {
		return fmt.Errorf("Unable to update multisig account because signers are less than threshold.")
	}
	mg.Threshold = threshold
	mg.Signers = signers
	return nil
}

func (mg *MultiSig) UpdateOwner(owner sdk.AccAddress) error {
	mg.Owner = owner
	return nil
}

func (mg *MultiSig) IsOwner(owner sdk.AccAddress) bool {
	return mg.Owner.Equals(owner)
}

func (mg *MultiSig) AddSignature(txID uint64, signer sdk.AccAddress, sig []byte) (exported.PendingTx, error) {
	stdSig := StdSignature{Signature: sig}
	for i, ptx := range mg.PendingTxs {
		if ptx.ID == txID {
			if ptx.IsSignedBy(signer) {
				return nil, fmt.Errorf("Signer has signed this transaction before.")
			}
			mg.PendingTxs[i].StdTx.Signatures = append(ptx.StdTx.Signatures, stdSig)
			mg.PendingTxs[i].SignedBy = append(ptx.SignedBy, signer)
			return mg.PendingTxs[i], nil
		}
	}

	return nil, fmt.Errorf("There is no pending transaction with this id:%v", txID)
}
