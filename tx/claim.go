package tx

import (
	"encoding/json"
	"io"

	"github.com/inwecrypto/neogo"
)

// ClaimTx .
type ClaimTx Transaction

type claimTx struct {
	inputs []*Vin
}

// NewClaimTx .
func NewClaimTx() *ClaimTx {
	tx := &ClaimTx{
		Type: ClaimTransaction,
	}

	return tx
}

// JSON .
func (tx *ClaimTx) JSON() string {
	data, _ := json.Marshal(tx.Inputs)

	return string(data)
}

// Tx get basic transaction object
func (tx *ClaimTx) Tx() *Transaction {
	return (*Transaction)(tx)
}

// Claim .
func (tx *ClaimTx) Claim(outputs []*Vout, unspent []*neogo.UTXO) error {

	gasOutputs := make([]*Vout, 0)
	otherOutputs := make([]*Vout, 0)

	for _, vout := range outputs {
		if vout.Asset == GasAssert {
			gasOutputs = append(gasOutputs, vout)
		} else {
			otherOutputs = append(otherOutputs, vout)
		}
	}

	gasInputs, unselected, err := (*Transaction)(tx).CalcInputs(gasOutputs, unspent)

	if err != nil {
		return err
	}

	tx.Extend = &claimTx{
		inputs: gasInputs,
	}

	otherInputs, _, err := (*Transaction)(tx).CalcInputs(otherOutputs, unselected)

	if err != nil {
		return err
	}

	tx.Inputs = otherInputs

	return nil
}

func (tx *claimTx) Write(writer io.Writer) error {

	length := Varint(len(tx.inputs))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, vin := range tx.inputs {
		if err := (*Vin)(vin).Write(writer); err != nil {
			return err
		}
	}

	return nil
}

func (tx *claimTx) Read(reader io.Reader) error {
	var length Varint

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		var vin Vin

		if err := vin.Read(reader); err != nil {
			return err
		}

		tx.inputs = append(tx.inputs, &vin)
	}

	return nil
}
