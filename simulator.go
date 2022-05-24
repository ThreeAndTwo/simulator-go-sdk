package simulator_go_sdk

func NewSimulator(platformUrl, rpc, pkStr string) (ISimulator, error) {
	return newEvmSimulator(platformUrl, rpc, pkStr)
}
