package tx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/dynamicgo/config"
	"github.com/inwecrypto/neogo"
	"github.com/inwecrypto/neogo/keystore"
	"github.com/inwecrypto/neogo/nep5"
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

func TestNep5RPC(t *testing.T) {

	client := neogo.NewClient(conf.GetString("neotest", "xxxxx"))

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	from := ToInvocationAddress(key.Address)

	tokenBalance, err := client.Nep5BalanceOf("849d095d07950b9e56d0c895ec48ec5100cfdff1", from)

	assert.NoError(t, err)

	println(tokenBalance)

	decimals, err := client.Nep5Decimals("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	assert.NoError(t, err)

	println(decimals)

	symbol, err := client.Nep5Symbol("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	assert.NoError(t, err)

	println("symbol: ", symbol)

	result, err := client.Nep5Transfer("849d095d07950b9e56d0c895ec48ec5100cfdff1", from, from, 1)

	assert.NoError(t, err)

	println(result.Script, result.GasConsumed)
}

func getAsset(address string, asset string) ([]*neogo.UTXO, error) {
	client := neogo.NewClient(conf.GetString("neotest", "xxxxx") + "/extend")

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

	client := neogo.NewClient(conf.GetString("neotest", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestTransfer(t *testing.T) {
	client := neogo.NewClient(conf.GetString("neotest", "xxxxx"))

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	from := ToInvocationAddress(key.Address)

	result, err := client.Nep5Transfer("849d095d07950b9e56d0c895ec48ec5100cfdff1", from, from, 1)

	assert.NoError(t, err)

	scriptHash, _ := hex.DecodeString("849d095d07950b9e56d0c895ec48ec5100cfdff1")

	scriptHash = reverseBytes(scriptHash)

	println(result.Script, result.GasConsumed, hex.EncodeToString(scriptHash))

	bytesOfFrom, _ := hex.DecodeString(from)

	bytesOfFrom = reverseBytes(bytesOfFrom)

	script, err := nep5.Transfer(scriptHash, bytesOfFrom, bytesOfFrom, big.NewInt(1))

	assert.Equal(t, result.Script, hex.EncodeToString(script))

	client2 := neogo.NewClient(conf.GetString("neotest", "xxxxx") + "/extend")

	utxos, err := client2.GetBalance(key.Address, GasAssert)

	assert.NoError(t, err)

	printResult(utxos)

	tx := NewInvocationTx(script, 0)

	err = tx.CalcInputs(nil, utxos)

	assert.NoError(t, err)

	rawtx, _, err := tx.Tx().Sign(key.PrivateKey)

	assert.NoError(t, err)

	println(tx.Tx().String())

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

	client := neogo.NewClient(conf.GetString("neotest", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestGetClaim(t *testing.T) {

	client := neogo.NewClient(conf.GetString("neotest", "xxxxx") + "/extend")

	key, err := keystore.KeyFromWIF(conf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

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

	client = neogo.NewClient(conf.GetString("neotest", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)

}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}
