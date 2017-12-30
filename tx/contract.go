package tx

// ContractTx contract transaction
type ContractTx Transaction

// NewContractTx create new contract transaction
func NewContractTx() *ContractTx {
	tx := &ContractTx{}

	tx.Type = ContractTransaction

	return tx
}
