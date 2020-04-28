package tests

const (
	TestChainID           = "test-chain"
	TestNet               = false
	TestMsgsPerTxPrepare  = 100
	TestMsgsPerTxLoadTest = 4
	TestTPS               = 100
	TestDuration          = 4
	TestRampUpTime        = 2
	TestMaxWorkers        = 100
	TestTargetURL         = "http://testurl.com:1234"
	TestLoadGeneratorURL  = "http://test-lg-url.com:1234"
	TestCustomURL         = "/custom/url"
	TestCoinName          = "link"
	TestMasterMnemonic    = "embrace catch hover lab birth gap gorilla price boost chapter vicious crowd draft announce skin swift harvest stage gas fragile artwork bid solar village"
	TestMnemonic          = "fever tell fancy ridge fly glow reflect decline voice coil reflect ski empty forum frost rebuild slide nut invite chase swarm flag dizzy diet"
	TestMnemonic2         = "tribe slot pioneer either light fossil broken scissors okay special update place want trash soldier rural portion lock couple venue cushion bind enact one"
	InvalidMnemonic       = "invalid mnemonic"

	ExpectedNumTargets  = TestTPS * TestDuration
	ExpectedAttackCount = (TestDuration-TestRampUpTime/2)*TestTPS + TestDuration
)

var TestNumPrepareRequest = TestTPS*TestDuration/TestMsgsPerTxPrepare + isPositive(TestTPS*TestDuration%TestMsgsPerTxPrepare)

func isPositive(x int) int {
	if x > 0 {
		return 1
	}
	return 0
}
