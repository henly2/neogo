package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dynamicgo/config"
	"github.com/stretchr/testify/assert"
)

var cnf *config.Config

func init() {
	cnf, _ = config.NewFromFile("../../conf/test.json")
}

func TestRPCAccountSate(t *testing.T) {
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	accoutState, err := client.GetAccountState("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	assert.NoError(t, err)

	printResult(accoutState)
}

func TestGetBalance(t *testing.T) {
	client := NewClient(cnf.GetString("neo", "xxxxx"))

	balance, err := client.GetBalance("AerEgjukZpi6kkx7ER9bg77YuNHaUM6u2q", "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(balance)
}

func TestConvertHex(t *testing.T) {
	valBytes, err := hex.DecodeString("00e40b5402")

	valBytes = reverseBytes(valBytes)

	require.NoError(t, err)

	fmt.Printf("%d\n", new(big.Int).SetBytes(valBytes))
}

func TestGetClaim(t *testing.T) {
	client := NewClient(cnf.GetString("neoexttest", "xxxxx"))

	balance, err := client.GetClaim("Ab8vffxvjaA3JKm3weBg6ChmZMSvorMoBM")

	assert.NoError(t, err)

	printResult(balance)
}

func TestConnectionCount(t *testing.T) {
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	count, err := client.GetConnectionCount()

	assert.NoError(t, err)

	fmt.Printf("connection count :%d\n", count)
}

func TestBestBlockHash(t *testing.T) {
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	hash, err := client.GetBestBlockHash()

	assert.NoError(t, err)

	block, err := client.GetBlock(hash)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("the best block :\n\t%s\n", string(blockjson))
}

func TestBlockCount(t *testing.T) {
	client := NewClient(cnf.GetString("neotest", "xxxxx"))

	count, err := client.GetBlockCount()

	assert.NoError(t, err)

	fmt.Printf("the block count :%d\n", count)
}

func TestBlockByIndex(t *testing.T) {
	client := NewClient(cnf.GetString("neotest2", "xxxxx"))

	block, err := client.GetBlockByIndex(2003)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("the best block :\n\t%s\n", string(blockjson))
}

func TestGetRawTransaction(t *testing.T) {
	client := NewClient(cnf.GetString("neotest2", "xxxxx"))

	block, err := client.GetRawTransaction("0x9a00f762ec56299df07bf38039e06c4722cb66331acc76240e53c2afdf998972")

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("trans:\n\t%s\n", string(blockjson))
}

func TestGetTxOut(t *testing.T) {
	client := NewClient(cnf.GetString("neo", "xxxxx"))

	block, err := client.GetTxOut("15e7c13851d28b4a049082dedba368f8772d6d829c77b9948019b3232a7c356d", 0)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("trans:\n\t%s\n", string(blockjson))
}

func TestGetPeers(t *testing.T) {
	client := NewClient(cnf.GetString("mainnode", "xxxxx"))

	block, err := client.GetPeers()

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("peers:\n\t%s\n", string(blockjson))
}

func TestSendRawTransaction(t *testing.T) {
	rawtx, err := hex.DecodeString("8000000164c73796d6ad5b73842a15ecd95e2899a174d2b28bd52013ee53952892bb7c9e0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e1f505000000004263d1f1b124778d66d847801fe7cb73dd4bef5001414054a3ac89b5770f9d6430d65cc4e3fa14de9c6636c4ccb6931d5c1e322d19229c431d85e5faecdbefabf4f713f32a356dbed55851178280d75e0361f00fa1acb723210398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99aac")

	assert.NoError(t, err)

	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	_, err = client.SendRawTransaction(rawtx)

	assert.NoError(t, err)
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestApplicationLog(t *testing.T) {
	client := NewClient(cnf.GetString("neotest2", "xxxxx"))

	result, err := client.ApplicationLog("0x55f766d1e508e44be3544de2676f74b3655250749cb81461f02c594779f4d596")

	assert.NoError(t, err)

	printResult(result)
}

func TestConvert(t *testing.T) {
	data, _ := hex.DecodeString("00ca9a3b00")

	data = reverseBytes(data)

	println(fmt.Sprintf("%d", new(big.Int).SetBytes(data)))
}
