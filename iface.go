package simulator_go_sdk

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type (
	ISimulator interface {
		SimulateTx(tx *TxSimulate) (string, error)
		DebuggerTx(txHash []string) (string, error)
		Mine(count uint64) error
		ResetBlockNumber(uint64) error
		SetBalance(common.Address) error
		BlockNumber() (uint64, error)
		BalanceAt(common.Address) (float64, error)
		NonceAt(account common.Address, blockNum *big.Int) (uint64, error)

		SuggestGasPrice() (*big.Int, error)
		PendingBalanceAt(account common.Address) (*big.Int, error)
		PendingStorageAt(account common.Address, key common.Hash) ([]byte, error)
		PendingCodeAt(account common.Address) ([]byte, error)
		PendingNonceAt(account common.Address) (uint64, error)
		PendingTransactionCount() (uint, error)
	}

	IHardhat interface {
		HardhatImpersonateAccount(account []string) (string, error)
		HardhatStopImpersonatingAccount(account []string) (string, error)
		HardhatMine(count uint64) (string, error)
		HardhatSetBalance(account string) (string, error)
		HardhatDebugTransaction(txHash []string) (string, error)
	}
)
