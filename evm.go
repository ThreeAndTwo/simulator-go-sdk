package simulator_go_sdk

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"strconv"
)

var (
	ErrPk       = errors.New("privateKey error")
	ErrTxParams = errors.New("tx params error")
	ErrJsonRpc  = errors.New("request jsonrpc error")
	ErrAddress  = errors.New("address len ne 42")
)

type evmSimulator struct {
	client     *ethclient.Client
	api        *hApi
	hardhatApi *hardhat
	account    *account
	header     map[string]string
}

type account struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

func (es *evmSimulator) SuggestGasPrice() (*big.Int, error) {
	return es.client.SuggestGasPrice(context.Background())
}

func (es *evmSimulator) PendingBalanceAt(account common.Address) (*big.Int, error) {
	return es.client.PendingBalanceAt(context.Background(), account)
}

func (es *evmSimulator) PendingStorageAt(account common.Address, key common.Hash) ([]byte, error) {
	return es.client.PendingStorageAt(context.Background(), account, key)
}

func (es *evmSimulator) PendingCodeAt(account common.Address) ([]byte, error) {
	return es.client.PendingCodeAt(context.Background(), account)
}

func (es *evmSimulator) PendingNonceAt(account common.Address) (uint64, error) {
	return es.client.PendingNonceAt(context.Background(), account)
}

func (es *evmSimulator) PendingTransactionCount() (uint, error) {
	return es.client.PendingTransactionCount(context.Background())
}

func newEvmSimulator(platformUrl, rpc, pkStr string) (*evmSimulator, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	header := initHeader()

	api := newHApi(platformUrl, header)
	hardhatApi := newHardhat(rpc, header)

	priKey := getAccountFromPKStr(pkStr)
	if priKey == nil {
		return nil, ErrPk
	}

	_account := &account{
		address:    crypto.PubkeyToAddress(priKey.PublicKey),
		privateKey: priKey,
	}
	return &evmSimulator{client: client, api: api, hardhatApi: hardhatApi, account: _account, header: header}, nil
}

func getAccountFromPKStr(pkStr string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(pkStr)
	if err != nil {
		return nil
	}
	return privateKey
}

func (es *evmSimulator) checkTx(tx *TxSimulate) bool {
	return tx.ChainId >= 1 && tx.BlockNumber >= 1 && tx.TransactionIndex >= 0 &&
		len(tx.From) == 42 && tx.Nonce >= 0 && len(tx.To) == 42 && tx.GasLimit >= 21_000 && tx.GasPrice >= "1"
}

func (es *evmSimulator) SimulateTx(tx *TxSimulate) (string, error) {
	// check params
	if !es.checkTx(tx) {
		return "", ErrTxParams
	}

	// reset to hardhat default chainId
	tx.ChainId = 31337

	// callHardhatImpersonateAccount
	accountRes, err := es.hardhatApi.HardhatImpersonateAccount([]string{tx.From})
	if err != nil {
		return "", err
	}

	_, err = decodeJsonRpc(accountRes)
	if err != nil {
		return "", err
	}

	if tx.Overrides.BlockNum != 0 {
		err = es.api.Reset(tx.Overrides.BlockNum)
		if err != nil {
			return "", err
		}
		tx.BlockNumber = tx.Overrides.BlockNum
	}

	if err = es.autoMine(tx); err != nil {
		return "", err
	}

	// set nonce
	_nonce, err := es.NonceAt(es.account.address, nil)
	if err != nil {
		return "", err
	}
	tx.Nonce = _nonce
	fmt.Printf("nonce: %d \n", tx.Nonce)

	// wrapperTx
	signerTx, err := es.wrapperTx(tx)
	if err != nil {
		return "", err
	}

	// sendTransaction
	err = es.client.SendTransaction(context.Background(), signerTx)
	if err != nil {
		return "", err
	}

	// stop HardhatImpersonateAccount
	_, err = es.hardhatApi.HardhatStopImpersonatingAccount([]string{tx.From})
	if err != nil {
		return "", err
	}

	return signerTx.Hash().Hex(), nil
}

func fmtTxVal(tx *TxSimulate) (*big.Int, error) {
	if tx.Value == "" {
		tx.Value = "0"
	}

	fVal, err := strconv.ParseFloat(tx.Value, 10)
	if err != nil {
		return nil, err
	}
	return parseAmountFloat(fVal, 18), nil
}

func parseAmountFloat(amount float64, decimals int) *big.Int {
	value := decimal.NewFromFloat(math.Pow10(decimals)).Mul(decimal.NewFromFloat(amount))
	return value.BigInt()
}

func (es *evmSimulator) wrapperTx(tx *TxSimulate) (*types.Transaction, error) {
	to := common.HexToAddress(tx.To)

	val, err := fmtTxVal(tx)
	if err != nil {
		return nil, err
	}

	data, _ := hexutil.Decode(tx.Input)
	// convert gasFee
	iFeePrice, err := strconv.ParseInt(tx.GasPrice, 10, 64)
	if err != nil {
		return nil, err
	}
	feePrice := big.NewInt(iFeePrice)

	// convert gasTips
	iTipPrice, err := strconv.ParseInt(tx.GasTips, 10, 64)
	if err != nil {
		return nil, err
	}
	tipPrice := big.NewInt(iTipPrice)
	signer := types.LatestSignerForChainID(big.NewInt(int64(tx.ChainId)))
	signerTx, err := types.SignNewTx(es.account.privateKey, signer, &types.DynamicFeeTx{
		ChainID:    big.NewInt(int64(tx.ChainId)),
		Nonce:      tx.Nonce,
		GasTipCap:  tipPrice,
		GasFeeCap:  feePrice,
		Gas:        tx.GasLimit,
		To:         &to,
		Value:      val,
		Data:       data,
		AccessList: tx.AccessList,
	})
	return signerTx, nil
}

// autoMine
// if tx.BlockNumber > current(blockNum) auto mine block to transaction block number
// else ignore the step
func (es *evmSimulator) autoMine(tx *TxSimulate) error {
	blockNum, err := es.BlockNumber()
	if err != nil {
		return err
	}

	if tx.BlockNumber > blockNum {
		_, err = es.hardhatApi.HardhatMine(tx.BlockNumber - blockNum)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func decodeJsonRpc(data string) (interface{}, error) {
	res := &NodeRes{}
	if err := json.Unmarshal([]byte(data), res); err != nil {
		return "", err
	}

	if res.JsonRpc == "2.0" && res.Id == 1 {
		return res.Result, nil
	}
	return "", ErrJsonRpc
}

func (es *evmSimulator) DebuggerTx(txHash []string) (string, error) {
	traces, err := es.hardhatApi.HardhatDebugTransaction(txHash)
	if err != nil {
		return "", err
	}
	return traces, nil
}

func (es *evmSimulator) Mine(count uint64) (err error) {
	_, err = es.hardhatApi.HardhatMine(count)
	return
}

func (es *evmSimulator) ResetBlockNumber(blockNum uint64) error {
	return es.api.Reset(blockNum)
}

func (es *evmSimulator) SetBalance(account common.Address) error {
	_, err := es.hardhatApi.HardhatSetBalance(account.Hex())
	return err
}

func (es *evmSimulator) BlockNumber() (uint64, error) {
	return es.client.BlockNumber(context.Background())
}

func (es *evmSimulator) NonceAt(account common.Address, blockNum *big.Int) (uint64, error) {
	return es.client.NonceAt(context.Background(), account, blockNum)
}

func (es *evmSimulator) BalanceAt(account common.Address, blockNum *big.Int) (float64, error) {
	if !es.checkAddress(account) {
		return 0, ErrAddress
	}

	balanceAt, err := es.client.BalanceAt(context.Background(), account, blockNum)
	if err != nil {
		return 0, err
	}

	balance := formatAmount(balanceAt, 18)
	return balance, nil
}

func formatAmount(amount *big.Int, decimal int) float64 {
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	value, _ := new(big.Float).Quo(new(big.Float).SetInt(amount), tenDecimal).Float64()
	return value
}

func (es *evmSimulator) checkAddress(account common.Address) bool {
	return len(account.Hex()) == 42
}
