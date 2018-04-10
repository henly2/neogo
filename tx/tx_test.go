package tx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dynamicgo/config"
	"github.com/inwecrypto/neogo/keystore"
	"github.com/inwecrypto/neogo/nep5"
	"github.com/inwecrypto/neogo/rpc"
	"github.com/stretchr/testify/assert"
)

var conf *config.Config
var scriptHash []byte
var scriptAddress string

func init() {
	conf, _ = config.NewFromFile("../../conf/test.json")

	scriptHash, _ = hex.DecodeString("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	scriptHash = reverseBytes(scriptHash)

	scriptAddress = encodeAddress(scriptHash)

	println(scriptAddress)
}

func TestFixed8ReadWrite(t *testing.T) {
	val := float64(0.00013874)

	fixed8 := MakeFixed8(val)

	println(fixed8.String())

	assert.Equal(t, fixed8.Float64(), val)

	var buff bytes.Buffer

	assert.NoError(t, fixed8.Write(&buff))

	var other Fixed8

	assert.NoError(t, other.Read(&buff))

	assert.Equal(t, fixed8.Float64(), other.Float64())
}

func TestEncodeDecodeAddress(t *testing.T) {
	address, err := decodeAddress("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")
	assert.NoError(t, err)

	assert.Equal(t, encodeAddress(address), "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")
}

func TestVarint(t *testing.T) {
	varint := Varint(253)

	var buff bytes.Buffer

	assert.NoError(t, varint.Write(&buff))

	println(hex.EncodeToString(buff.Bytes()))

	assert.Equal(t, len(buff.Bytes()), 3)

	var other Varint

	assert.NoError(t, other.Read(&buff))

	assert.Equal(t, other, varint)

	buff.Reset()

	varint = Varint(252)

	assert.NoError(t, varint.Write(&buff))

	println(hex.EncodeToString(buff.Bytes()))

	assert.Equal(t, len(buff.Bytes()), 1)
}

func TestPrintTx(t *testing.T) {
	tx := &Transaction{}

	tx.Scripts = []*Scripts{
		&Scripts{
			StackScript:  make([]byte, 0),
			RedeemScript: make([]byte, 0),
		},
	}

	// var buff bytes.Buffer

	// writer := neogo.NewScriptWriter(&buff)

	// writer.EmitPushBytes()

	// inv := NewInvocationTx()

	println(tx.String())
}

func TestA(t *testing.T) {
	address, _ := hex.DecodeString("8cec4a755be0fac1613df2b549798ca25ea0e37e")

	address = reverseBytes(address)

	println(encodeAddress(address))
}

func TestNep5RPC(t *testing.T) {

	client := rpc.NewClient(conf.GetString("neo", "xxxxx"))

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	println(key.Address)

	from := ToInvocationAddress("AYYiDtPGaxt7rVtEEp9tiw4wgtg8jVEnSP")

	// // from := "8cec4a755be0fac1613df2b549798ca25ea0e37e"

	// tokenBalance, err := client.Nep5BalanceOf("9beb45a55bbc1880043e6bcd734805a22be8371b", from)

	// assert.NoError(t, err)

	// println(tokenBalance)

	// decimals, err := client.Nep5Decimals("9beb45a55bbc1880043e6bcd734805a22be8371b")

	// assert.NoError(t, err)

	// println(decimals)

	// symbol, err := client.Nep5Symbol("9beb45a55bbc1880043e6bcd734805a22be8371b")

	// assert.NoError(t, err)

	// println("symbol: ", symbol)

	result, err := client.Nep5Transfer("ab719b8baa2310f232ee0d277c061704541cfb", from, from, 1)

	assert.NoError(t, err)

	println(result.Script, result.GasConsumed)
}

func getAsset(address string, asset string) ([]*rpc.UTXO, error) {
	client := rpc.NewClient(conf.GetString("neo", "xxxxx") + "/extend")

	return client.GetBalance(address, asset)
}

func TestMintToken(t *testing.T) {
	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	gasAsset, err := getAsset(key.Address, GasAssert)

	assert.NoError(t, err)

	neoAsset, err := getAsset(key.Address, NEOAssert)

	assert.NoError(t, err)

	asset := append(gasAsset, neoAsset...)

	printResult(asset)

	from := ToInvocationAddress(key.Address)

	bytesOfFrom, _ := hex.DecodeString(from)

	bytesOfFrom = reverseBytes(bytesOfFrom)

	script, err := nep5.MintToken(scriptHash)

	nonce, _ := time.Now().MarshalBinary()

	tx := NewInvocationTx(script, 0, bytesOfFrom, nonce)

	vout := []*Vout{
		&Vout{
			Asset:   NEOAssert,
			Value:   MakeFixed8(1),
			Address: scriptAddress,
		},
	}

	err = tx.CalcInputs(vout, asset)

	assert.NoError(t, err)

	rawtx, _, err := tx.Tx().Sign(key.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

	client := rpc.NewClient(conf.GetString("neotest", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestTimeNow(t *testing.T) {
	println(time.Now().String())
}

type Test struct {
	Data string `json:"data"`
}

func TestTransfer(t *testing.T) {
	client := rpc.NewClient(conf.GetString("neotest2", "xxxxx"))

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	key2, err := keystore.KeyFromWIF(conf.GetString("wallet2", "xxxxx"))

	assert.NoError(t, err)

	from := ToInvocationAddress(key.Address)

	to := ToInvocationAddress(key2.Address)

	// result, err := client.Nep5Transfer("849d095d07950b9e56d0c895ec48ec5100cfdff1", from, to, 100000000)

	// assert.NoError(t, err)

	scriptHash, _ := hex.DecodeString("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	scriptHash = reverseBytes(scriptHash)

	// println(result.Script, result.GasConsumed, hex.EncodeToString(scriptHash))

	bytesOfFrom, _ := hex.DecodeString(from)

	bytesOfFrom = reverseBytes(bytesOfFrom)

	bytesOfTo, _ := hex.DecodeString(to)

	bytesOfTo = reverseBytes(bytesOfTo)

	script, err := nep5.Transfer(scriptHash, bytesOfFrom, bytesOfTo, big.NewInt(100000000))

	// assert.Equal(t, result.Script, hex.EncodeToString(script))

	// client2 := rpc.NewClient(conf.GetString("neotest", "xxxxx") + "/extend")

	// utxos, err := client2.GetBalance(key.Address, GasAssert)

	assert.NoError(t, err)

	// data, _ := json.Marshal(utxos)

	// test := &Test{
	// 	Data: string(data),
	// }

	// printResult(test)

	nonce, _ := time.Now().MarshalBinary()

	tx := NewInvocationTx(script, 0, bytesOfFrom, nonce)

	// err = tx.CalcInputs(nil, utxos)

	assert.NoError(t, err)

	// tx.CheckFromWitness(bytesOfFrom)

	rawtx, _, err := tx.Tx().Sign(key.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

	// rawtx, _ := hex.DecodeString("d101500500e1f50500140debf40cabd7c745bb8baa85bdf579ad380bc37e144263d1f1b124778d66d847801fe7cb73dd4bef5053c1087472616e7366657267f1dfcf0051ec48ec95c8d0569e0b95075d099d84000000000000000000011962ffadd11147311d0a85f6489ba2981737b00d95ee72901d3a006d713943d9000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6011c05548170000004263d1f1b124778d66d847801fe7cb73dd4bef50014140bdc48547b96aaa302cd13a0c968d9395b2c9ae3e0e2d3a01c8c1c23dcca0e01f6e4f23418350d102a53f7baf4f7a55cf96695afe9749e14ad925abeeca25bee323210398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99aac")

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestUnmarshalTx(t *testing.T) {
	tx := NewInvocationTx(nil, 0, []byte{}, []byte{})

	data, err := hex.DecodeString("d101500500e1f50500147ee3a05ea28c7949b5f23d61c1fae05b754aec8c144263d1f1b124778d66d847801fe7cb73dd4bef5053c1087472616e7366657267f1dfcf0051ec48ec95c8d0569e0b95075d099d84000000000000000002204263d1f1b124778d66d847801fe7cb73dd4bef50ff0f010000000ed24eb8d32402572601e00000014140d70b2dc42059ed47ee674130533267556ee272e33e003bcdd84783236e4fb3111535352fd9c0a1f586930a43181ea1508c02a08e036a180bdd1125d9e9c5d72e23210398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99aac")

	assert.NoError(t, err)

	assert.NoError(t, tx.Tx().Read(bytes.NewBuffer(data)))

	println(tx.Tx().String())
}

func TestTransferNEO(t *testing.T) {
	to, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	from, err := keystore.KeyFromWIF(conf.GetString("wallet2", "xxxxx"))

	assert.NoError(t, err)

	tx := NewContractTx()

	vout := []*Vout{
		&Vout{
			Asset:   NEOAssert,
			Value:   MakeFixed8(1),
			Address: to.Address,
		},
	}

	asset, err := getAsset(from.Address, NEOAssert)

	assert.NoError(t, err)

	err = tx.CalcInputs(vout, asset)

	assert.NoError(t, err)

	rawtx, _, err := tx.Tx().Sign(from.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

	client := rpc.NewClient(conf.GetString("neotest", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestUTXO(t *testing.T) {
	utxos, err := getAsset("AXXYB4NBu1uDChFmE4vkbCgNrq2tvkxQuK", NEOAssert)
	require.NoError(t, err)

	printResult(utxos)
}

func TestGetClaim(t *testing.T) {
	client := rpc.NewClient(conf.GetString("neo", "xxxxx") + "/extend")

	claims, err := client.GetClaim("AJFnsA8y2UFwnqcris5KrnAijK2qCMtu7R")

	assert.NoError(t, err)

	printResult(claims)
}

var key *keystore.Key

func init() {
	rawdata, err := ioutil.ReadFile("../../conf/lala2.json")

	if err != nil {
		panic(err)
	}

	key, err = keystore.ReadKeyStore(rawdata, "Lalala123")

	if err != nil {
		panic(err)
	}
}

func TestDoClaim(t *testing.T) {

	client := rpc.NewClient(conf.GetString("neo", "xxxxx") + "/extend")

	// key, err := keystore.KeyFromWIF(conf.GetString("wallet2", "xxxxx"))

	// assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	printResult(claims)

	val, err := strconv.ParseFloat(claims.Available, 64)

	assert.NoError(t, err)

	tx := NewClaimTx()

	err = tx.Claim(val, key.Address, claims.Claims)

	assert.NoError(t, err)

	rawtx, _, err := tx.Tx().Sign(key.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

	client = rpc.NewClient(conf.GetString("neo", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)

}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestFixed8(t *testing.T) {

	valueBytes, _ := hex.DecodeString("ffd31fd7e800")

	valueBytes = reverseBytes(valueBytes)

	// fixed := Fixed8(new(big.Int).SetBytes(valueBytes).Int64())

	println(hex.EncodeToString(valueBytes))

	fixed := MakeFixed8(10000.41600000)

	println(hex.EncodeToString(fixed.Int().Bytes()))

}
