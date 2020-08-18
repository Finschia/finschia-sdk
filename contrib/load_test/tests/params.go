package tests

const (
	TestChainID           = "test-chain"
	TestNet               = false
	TestMsgsPerTxPrepare  = 500
	TestMsgsPerTxLoadTest = 3
	TestMaxGasPrepare     = uint64(240000000)
	TestMaxGasLoadTest    = uint64(10000000)
	TestTPS               = 100
	TestDuration          = 4
	TestRampUpTime        = 2
	TestMaxWorkers        = 100
	TestTargetURL         = "http://testurl.com:1234"
	TestLoadGeneratorURL  = "http://test-lg-url.com:1234"
	TestCoinName          = "link"
	TestMasterMnemonic    = "embrace catch hover lab birth gap gorilla price boost chapter vicious crowd draft announce skin swift harvest stage gas fragile artwork bid solar village"
	TestMnemonic          = "fever tell fancy ridge fly glow reflect decline voice coil reflect ski empty forum frost rebuild slide nut invite chase swarm flag dizzy diet"
	TestMnemonic2         = "tribe slot pioneer either light fossil broken scissors okay special update place want trash soldier rural portion lock couple venue cushion bind enact one"
	InvalidMnemonic       = "invalid mnemonic"
	TokenContractID       = "9be17165"
	CollectionContractID  = "678c146a"
	FTTokenID             = "0000000100000000"
	NFTTokenType          = "10000001"
	TxHash                = "D20985E8B70B54B7C79D37B8E214EE815EB8D9818CF793A20304678FFA2A4A92"
	Address               = "link1muu5cza33kttadr5wylsqhfgxnwlcdrxls0wwn"

	ExpectedNumTargets      = TestTPS * TestDuration
	ExpectedAttackCount     = (TestDuration-TestRampUpTime/2)*TestTPS + TestDuration
	TestMaxAttackDifference = 10
)

var TestNumPrepareRequest = GetNumPrepareTx(TestTPS*TestDuration, TestMsgsPerTxPrepare)

func GetNumPrepareTx(numMsgs, msgsPerTxPrepare int) int {
	return (numMsgs + msgsPerTxPrepare - 1) / msgsPerTxPrepare
}
