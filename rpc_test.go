package neogo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dynamicgo/config"
	"github.com/stretchr/testify/assert"
)

var cnf *config.Config

func init() {
	cnf, _ = config.NewFromFile("./test.json")
}

func TestRPCAccountSate(t *testing.T) {
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	accoutState, err := client.GetAccountState("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	assert.NoError(t, err)

	printResult(accoutState)
}

func TestGetBalance(t *testing.T) {
	client := NewClient(cnf.GetString("testnodeext", "xxxxx"))

	balance, err := client.GetBalance("AJnNUn6HynVcco1p8LER72s4zXtNFYDnys", "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(balance)
}

func TestGetClaim(t *testing.T) {
	client := NewClient(cnf.GetString("testnodeext", "xxxxx"))

	balance, err := client.GetClaim("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

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
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	count, err := client.GetBlockCount()

	assert.NoError(t, err)

	fmt.Printf("the block count :%d\n", count)
}

func TestBlockByIndex(t *testing.T) {
	client := NewClient(cnf.GetString("testnode", "xxxxx"))

	block, err := client.GetBlockByIndex(4356)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("the best block :\n\t%s\n", string(blockjson))
}

func TestGetRawTransaction(t *testing.T) {
	client := NewClient(cnf.GetString("mainnode", "xxxxx"))

	block, err := client.GetRawTransaction("fafd00170b22480f46ee4506c0ddc97c02eeb23729c611d27329a78aa250a1e7")

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("trans:\n\t%s\n", string(blockjson))
}

func TestGetTxOut(t *testing.T) {
	client := NewClient(cnf.GetString("mainnode", "xxxxx"))

	block, err := client.GetTxOut("0x0ae13c1ba01d30a8238a0ec89019171fcf9eee61802dd468cc797a02ac48798d", 0)

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

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}
