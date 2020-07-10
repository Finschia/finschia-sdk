package tests

const (
	TestChainID           = "test-chain"
	TestNet               = false
	TestMsgsPerTxPrepare  = 500
	TestMsgsPerTxLoadTest = 4
	TestMaxGasPrepare     = uint64(240000000)
	TestMaxGasLoadTest    = uint64(1200000)
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

	ExpectedNumTargets      = TestTPS * TestDuration
	ExpectedAttackCount     = (TestDuration-TestRampUpTime/2)*TestTPS + TestDuration
	TestMaxAttackDifference = 10
)

var TestNumPrepareRequest = GetNumPrepareTx(TestTPS*TestDuration, TestMsgsPerTxPrepare)

func GetNumPrepareTx(numMsgs, msgsPerTxPrepare int) int {
	return (numMsgs + msgsPerTxPrepare - 1) / msgsPerTxPrepare
}
