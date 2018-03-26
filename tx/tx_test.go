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

	script, err := nep5.MintToken(scriptHash)

	tx := NewInvocationTx(script, 0)

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

func TestTransfer(t *testing.T) {
	client := rpc.NewClient(conf.GetString("neotest", "xxxxx"))

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	key2, err := keystore.KeyFromWIF(conf.GetString("wallet2", "xxxxx"))

	assert.NoError(t, err)

	from := ToInvocationAddress(key.Address)

	to := ToInvocationAddress(key2.Address)

	result, err := client.Nep5Transfer("849d095d07950b9e56d0c895ec48ec5100cfdff1", from, to, 100000000)

	assert.NoError(t, err)

	scriptHash, _ := hex.DecodeString("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	scriptHash = reverseBytes(scriptHash)

	println(result.Script, result.GasConsumed, hex.EncodeToString(scriptHash))

	bytesOfFrom, _ := hex.DecodeString(from)

	bytesOfFrom = reverseBytes(bytesOfFrom)

	bytesOfTo, _ := hex.DecodeString(to)

	bytesOfTo = reverseBytes(bytesOfTo)

	script, err := nep5.Transfer(scriptHash, bytesOfFrom, bytesOfTo, big.NewInt(100000000))

	// assert.Equal(t, result.Script, hex.EncodeToString(script))

	client2 := rpc.NewClient(conf.GetString("neotest", "xxxxx") + "/extend")

	utxos, err := client2.GetBalance(key.Address, GasAssert)

	assert.NoError(t, err)

	printResult(utxos)

	tx := NewInvocationTx(script, 0)

	// err = tx.CalcInputs(nil, utxos)

	// assert.NoError(t, err)

	tx.CheckFromWitness(bytesOfFrom)

	rawtx, _, err := tx.Tx().Sign(key.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

	// rawtx, _ := hex.DecodeString("80000001413c7fa8473898e38f3df6a28592859b715762a5ee9cc7e3f935c97538f0f71d0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e1f505000000001b37e82ba2054873cd6b3e048018bc5ba545eb9b014140b0bd5aab7cb3b04b989afe97d8f334bc2c8ecf25afa91446a580d6bca259955b21ee923f0272127194a0d00fb400e0fc0e181ef248108673bec9bcdd13dae4452321028c72ef5482e037f4795421df9c7a63fcc0e059e9314d9249e0cbf16570701bc1ac")

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestUnmarshalTx(t *testing.T) {
	tx := NewInvocationTx(nil, 0)

	data, err := hex.DecodeString("d1014b51144263d1f1b124778d66d847801fe7cb73dd4bef50144263d1f1b124778d66d847801fe7cb73dd4bef5053c1087472616e7366657267f1dfcf0051ec48ec95c8d0569e0b95075d099d8400000000000000000001b986577bcb2769bc328237d62830015a747a931e7030ac1b8d1f77bc5df8d443010001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60e3c82f21010000004263d1f1b124778d66d847801fe7cb73dd4bef50014140d8197ddab4fe5b46c473b98ff0a925b9e8c77e5e45ca280ac86f56515120a223d6c151be4b93855d6bdfb525b257923c9395d26c4373f7c818c9d8b74c7ef1d023210398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99aac")

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
