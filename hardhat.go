package simulator_go_sdk

import (
	"fmt"
	"github.com/deng00/req"
)

type hardhat struct {
	nodeUrl string
	api     *net
	header  map[string]string
}

func initHeader() map[string]string {
	header := make(map[string]string)
	header["content-type"] = "application/json"
	return header
}

func newHardhat(nodeUrl string, header map[string]string) *hardhat {
	_api := newNet(nodeUrl, header, nil)
	return &hardhat{nodeUrl: nodeUrl, api: _api, header: header}
}

func (h *hardhat) rpcCall(method HardhatMethod, params interface{}) (string, error) {
	h.api.Params = req.Param{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	return h.api.Request(PostTy)
}

func (h *hardhat) HardhatImpersonateAccount(account []string) (string, error) {
	return h.rpcCall(callHardhatImpersonateAccount, account)
}

func (h *hardhat) HardhatStopImpersonatingAccount(account []string) (string, error) {
	return h.rpcCall(callHardhatStopImpersonatingAccount, account)
}

func (h *hardhat) HardhatDebugTransaction(txHash []string) (string, error) {
	return h.rpcCall(debugTraceTransaction, txHash)
}

func (h *hardhat) HardhatMine(count uint64) (string, error) {
	if count == 0 {
		return "", fmt.Errorf("mine count eq 0")
	}
	return h.rpcCall(callHardhatMine, []string{fmt.Sprintf("0x%x", count)})
}

// HardhatSetBalance
// fixedValue, set 10_000 eth
func (h *hardhat) HardhatSetBalance(account string) (string, error) {
	params := []string{account, "0x21e19e0c9bab2400000"}
	return h.rpcCall(callHardhatSetBalance, params)
}
