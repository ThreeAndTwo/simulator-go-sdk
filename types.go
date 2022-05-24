package simulator_go_sdk

import "github.com/ethereum/go-ethereum/core/types"

type HardhatMethod string

const (
	callHardhatImpersonateAccount       HardhatMethod = "hardhat_impersonateAccount"
	callHardhatStopImpersonatingAccount HardhatMethod = "hardhat_stopImpersonatingAccount"
	callHardhatMine                     HardhatMethod = "hardhat_mine"
	callHardhatSetBalance               HardhatMethod = "hardhat_setBalance"
	debugTraceTransaction               HardhatMethod = "debug_traceTransaction"
)

type TxSimulate struct {
	ChainId          uint64 `json:"chain_id"`
	BlockNumber      uint64 `json:"block_number"`
	TransactionIndex int    `json:"transaction_index"`
	From             string `json:"from"`
	Nonce            uint64
	Input            string           `json:"input"`
	To               string           `json:"to"`
	GasLimit         uint64           `json:"gas"`
	GasPrice         string           `json:"gas_price"`
	GasTips          string           `json:"gas_tips"`
	Value            string           `json:"value"`
	AccessList       types.AccessList `json:"access_list"`
	Overrides        txOverrides      `json:"overrides"`
}

type txOverrides struct {
	BlockNum uint64 `json:"block_num"`
}

type NodeRes struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}
