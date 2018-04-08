package keystore

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dynamicgo/config"
	"github.com/inwecrypto/bip39"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conf *config.Config

func init() {
	conf, _ = config.NewFromFile("../../conf/test.json")
}

func TestAddress(t *testing.T) {
	fromBytes, _ := hex.DecodeString("2b41aea9d405fef2e809e3c8085221ce944527a7")
	println(base58.CheckEncode(fromBytes, 0x17))
}

func TestWIF(t *testing.T) {
	privateKey, err := DecodeWIF(conf.GetString("wallet", "xxx"))

	require.NoError(t, err)

	wif, err := EncodeWIF(privateKey)

	require.NoError(t, err)

	require.Equal(t, wif, conf.GetString("wallet", "xxx"))

	address, err := PrivateToAddress(privateKey)

	require.NoError(t, err)

	require.Equal(t, address, conf.GetString("walletaddr", "xx"))

}

func publicKeyToBytes(pub *ecdsa.PublicKey) (b []byte) {
	/* See Certicom SEC1 2.3.3, pg. 10 */

	x := pub.X.Bytes()

	/* Pad X to 32-bytes */
	paddedx := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedx...)
	}

	return append([]byte{0x03}, paddedx...)
}

func TestMnemonic(t *testing.T) {
	key, _ := KeyFromWIF("L4sSGSGh15dtocMMSYS115fhZEVN9UuETWDjgGKu2JDu59yncyVf")

	privateKeyBytes := key.ToBytes()

	dic, _ := bip39.GetDict("en_US")

	data, _ := bip39.NewMnemonic(privateKeyBytes, dic)

	println(len(privateKeyBytes), len(strings.Split(data, " ")))

	println(string(data))

	data2, err := bip39.MnemonicToByteArray(data, dic)

	data2 = data2[1 : len(data2)-1]

	assert.NoError(t, err)

	assert.Equal(t, privateKeyBytes, data2)

	key2, err := KeyFromPrivateKey(data2)

	assert.NoError(t, err)

	assert.Equal(t, key.Address, key2.Address)
}

func TestMnemonic2(t *testing.T) {
	dic, _ := bip39.GetDict("en_US")

	data2, err := bip39.MnemonicToByteArray("increase local decade among nerve brisk eyebrow palm law humble drama wreck mean endorse yard slight fiber entry harsh senior fuel fetch cekery panda", dic)

	require.NoError(t, err)

	// data2 = data2[1 : len(data2)-1]

	data2 = data2[1 : len(data2)-1]

	key2, _ := KeyFromPrivateKey(data2)

	printResult(key2)
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestNEOAddress(t *testing.T) {
	key, err := KeyFromWIF("L4Ns4Uh4WegsHxgDG49hohAYxuhj41hhxG6owjjTWg95GSrRRbLL")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t,
		hex.EncodeToString(toBytes(key.PrivateKey)),
		"d59208b9228bff23009a666262a800f20f9dad38b0d9291f445215a0d4542beb")

	assert.Equal(t, hex.EncodeToString(publicKeyToBytes(&key.PrivateKey.PublicKey)), "0398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99a")
	assert.Equal(t, key.Address, "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	ks, err := WriteLightScryptKeyStore(key, "test")

	assert.NoError(t, err)

	key2, err := ReadKeyStore(ks, "test")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t,
		hex.EncodeToString(toBytes(key2.PrivateKey)),
		"d59208b9228bff23009a666262a800f20f9dad38b0d9291f445215a0d4542beb")

	assert.Equal(t, hex.EncodeToString(publicKeyToBytes(&key2.PrivateKey.PublicKey)), "0398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99a")
	assert.Equal(t, key2.Address, "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

}
