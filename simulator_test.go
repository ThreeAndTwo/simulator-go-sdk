package simulator_go_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"testing"
)

const PlatformURL = "http://127.0.0.1:49160"
const RPC = "http://127.0.0.1:8545"

func TestEvmSimulator_SimulateTx(t *testing.T) {
	tests := []struct {
		name   string
		params string
		pkStr  string
	}{
		{
			name:   "transfer",
			params: "{\"chain_id\":1,\"block_number\":14829628,\"transaction_index\":0,\"from\":\"0x8Fd30ec7FF8B74bcbc3daB47601c3DE4Afb34A5E\",\"input\":\"0x\",\"to\":\"0xA51Fc19f0430614F22B9Caf10491298E5D571313\",\"gas\":21000,\"gas_price\":\"100\",\"gas_tips\":\"2\",\"value\":\"3\",\"access_list\":[]}\n",
			pkStr:  os.Getenv("PK"),
		},
		{
			name:   "mint ldr",
			params: "{\"chain_id\":1,\"block_number\":14829628,\"transaction_index\":0,\"from\":\"0x562A24fF5Abf5e6003b2a41C94B5B1470b0c0e3A\",\"input\":\"0xa3cba09a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000011b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000413e34d2baf0b8949f748ad66b44228a3f6365d2acf8053fd005fb6b9b717e52e544ea994f4443cf8d9597c247eb320185cb66442d1c5ab25570bb97acdf76f0a11c00000000000000000000000000000000000000000000000000000000000000\",\"to\":\"0xfd43d1da000558473822302e1d44d81da2e4cc0d\",\"gas\":140000,\"gas_price\":\"32000000000\",\"gas_tips\":\"2\",\"value\":\"0\",\"access_list\":[]}\n",
			pkStr:  os.Getenv("PK"),
		},
		{
			name:   "",
			params: "{\"chain_id\":1,\"block_number\":14829628,\"transaction_index\":0,\"from\":\"0x562A24fF5Abf5e6003b2a41C94B5B1470b0c0e3A\",\"input\":\"0xa3cba09a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000011b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000413e34d2baf0b8949f748ad66b44228a3f6365d2acf8053fd005fb6b9b717e52e544ea994f4443cf8d9597c247eb320185cb66442d1c5ab25570bb97acdf76f0a11c00000000000000000000000000000000000000000000000000000000000000\",\"to\":\"0xfd43d1da000558473822302e1d44d81da2e4cc0d\",\"gas\":140000,\"gas_price\":\"32000000000\",\"gas_tips\":\"2\",\"value\":\"0\",\"access_list\":[],\"overrides\":{\"block_num\":14834856}}\n",
			pkStr:  os.Getenv("PK"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulator, err := NewSimulator(PlatformURL, RPC, tt.pkStr)
			if err != nil {
				t.Fatalf("new simulator ins error: %s", err)
			}

			tx := &TxSimulate{}
			err = json.Unmarshal([]byte(tt.params), tx)
			if err != nil {
				t.Fatalf("fmt json error: %s", err)
			}

			blockNum, _ := simulator.BlockNumber()
			fmt.Printf("blockNum: %d \n", blockNum)

			simulateTx, err := simulator.SimulateTx(tx)
			if err != nil {
				t.Fatalf("simulator tx error: %s", err)
			}
			fmt.Printf("simulateTx: %s \n", simulateTx)
		})
	}
}

func TestEvmSimulator_DebuggerTx(t *testing.T) {
	tests := []struct {
		name  string
		pkStr string
		hash  []string
	}{
		{
			name:  "test trace",
			pkStr: os.Getenv("PK"),
			hash:  []string{"0xaf9dcfa3906f24e184f0fab41dc7285e7c584ae34070c6ca33da4a2e80afc064"},
		},
		{
			name:  "test mint ldr",
			pkStr: os.Getenv("PK"),
			hash:  []string{"0x8df3529bec2f2a3040ad471e83255a04958a01785b134ee14c25922407251ac0"},
		},
		{
			name:  "test error tx trace",
			pkStr: os.Getenv("PK"),
			hash:  []string{"0x1f0149efb88913dd76460d829eaf23174369fee390fdfdca638922ac409a0fab"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulator, err := NewSimulator(PlatformURL, RPC, tt.pkStr)
			if err != nil {
				t.Fatalf("new simulator ins error: %s", err)
			}

			traceLog, err := simulator.DebuggerTx(tt.hash)
			t.Logf("log: %s", traceLog)
			//t.Logf("err: %s", err)
		})
	}
}

func TestEvmSimulator_GetBalance(t *testing.T) {
	tests := []struct {
		name    string
		pkStr   string
		account common.Address
	}{
		{
			name:    "ledger",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress("0x632Ea37aAc7A086f74E2E3AD9062EAA29da5F7F8"),
		},
		{
			name:    "billing",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress("0xA51Fc19f0430614F22B9Caf10491298E5D571313"),
		},
		{
			name:    "pubKey",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress("0x8Fd30ec7FF8B74bcbc3daB47601c3DE4Afb34A5E"),
		},
		{
			name:    "blockhole null",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress(""),
		},
		{
			name:    "account is null",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress("0x00"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulator, err := NewSimulator(PlatformURL, RPC, tt.pkStr)
			if err != nil {
				t.Fatalf("new simulator ins error: %s", err)
			}

			balance, _ := simulator.BalanceAt(tt.account, nil)
			fmt.Printf("balance: %f \n", balance)
		})
	}
}

func TestEvmSimulator_SetBalance(t *testing.T) {
	tests := []struct {
		name    string
		pkStr   string
		account common.Address
	}{
		{
			name:    "test setBalance",
			pkStr:   os.Getenv("PK"),
			account: common.HexToAddress("0x8Fd30ec7FF8B74bcbc3daB47601c3DE4Afb34A5E"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulator, err := NewSimulator(PlatformURL, RPC, tt.pkStr)
			if err != nil {
				t.Fatalf("new simulator ins error: %s", err)
			}

			err = simulator.SetBalance(tt.account)
			t.Logf("account balance: %s", err)
		})
	}
}
