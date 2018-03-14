package script

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestScriptHash(t *testing.T) {
	data, _ := hex.DecodeString("0480969800146063795d3b9b3cd55aef026eae992b91063db0db14a0e3bf726d9790a51aa42d9d3e006b1c32b1e1ae53c1087472616e7366657267f91d6b7085db7c5aaf09f19eeec1ca3c0db2c6ecf166c19707e2f89c3a93")

	println(hex.EncodeToString(Hash(data)))
}

func TestSignBigInt(t *testing.T) {
	val := big.NewInt(-1)

	val2 := new(big.Int).SetBytes(val.Bytes())

	println(val.Int64(), val2.Int64())
}
