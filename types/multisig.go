package types

import "fmt"

type PendingTx struct {
	ID       uint64       `json:"id" yaml:"id "`
	Tx       Tx           `json:"tx" yaml:"tx"`
	Sender   AccAddress   `json:"sender" yaml:"sender"`
	SignedBy []AccAddress `json:"signers" yaml:"signers"`
}

type MultiSig struct {
	Owner      AccAddress   `json:"owner" yaml:"owner"`
	Threshold  int          `json:"threshold" yaml:"threshold"`
	Counter    uint64       `json:"counter" yaml:"counter"`
	Signers    []AccAddress `json:"signers" yaml:"signers"`
	PendingTxs []PendingTx  `json:"pendingTxs" yaml:"pendingTxs"`
}

// Multisig functions
func (mg MultiSig) ContainTx(txID uint64) (int, bool) {
	for i, e := range mg.PendingTxs {
		if e.ID == txID {
			return i, true
		}
	}
	return -1, false
}

func (mg MultiSig) GetSigners() []AccAddress {
	return mg.Signers
}

func (mg MultiSig) IsSigner(signer AccAddress) bool {
	for _, sig := range mg.Signers {
		if signer.Equals(sig) {
			return true
		}
	}
	return false
}

func (mg MultiSig) GetThreshold() int {
	return mg.Threshold
}

func (mg MultiSig) GetOwner() AccAddress {
	return mg.Owner
}

func (mg *MultiSig) AddTx(tx PendingTx) error {
	for _, v := range mg.PendingTxs {
		if v.ID == tx.ID {
			return fmt.Errorf("Tx ID is not unique.")
		}
	}
	mg.PendingTxs = append(mg.PendingTxs, tx)
	return nil
}

func (mg *MultiSig) RemoveTx(txID uint64) bool {
	index, exists := mg.ContainTx(txID)
	if exists {
		s := mg.PendingTxs
		mg.PendingTxs = append(s[:index], s[index+1:]...)
		return true
	}
	return false
}

func (mg *MultiSig) SignTx(signer AccAddress, txID uint64) {
	for k, v := range mg.PendingTxs {
		if v.ID == txID {
			mg.PendingTxs[k].SignedBy = append(mg.PendingTxs[k].SignedBy, signer)
		}
	}
}

func (mg MultiSig) GetPendingTx(txID uint64) (bool, PendingTx) {
	for _, v := range mg.PendingTxs {
		if v.ID == txID {
			return true, v
		}
	}
	return false, PendingTx{}
}

func (mg MultiSig) GetTx(txID uint64) Tx {
	for _, v := range mg.PendingTxs {
		if v.ID == txID {
			return v.Tx
		}
	}
	return nil
}

func (mg MultiSig) IsMetric(txID uint64) bool {
	for _, v := range mg.PendingTxs {
		if len(v.SignedBy) == mg.Threshold {
			return true
		}
	}
	return false
}

func (mg MultiSig) GetNewTxID() uint64 {
	return mg.Counter + 1
}

func (mg *MultiSig) IncCounter() {
	mg.Counter++
}

// PendingTx
func NewPendingTx(id uint64, tx Tx, sender AccAddress, signedBy []AccAddress) PendingTx {
	return PendingTx{
		ID:       id,
		Tx:       tx,
		Sender:   sender,
		SignedBy: signedBy,
	}
}
