package testutil

import (
	"fmt"
	"github.com/line/lbm-sdk/client/flags"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/wasm/client/cli"
	"github.com/line/lbm-sdk/x/wasm/keeper"
	"io/ioutil"
)

func (s *IntegrationTestSuite) TestInstantiateContractCmd() {
	val := s.network.Validators[0]
	owner := val.Address.String()

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid instantiateContract": {
			[]string{
				s.codeId,
				fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", owner, keeper.RandomAccountAddress(s.T())),
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			true,
		},
		"wrong args count": {
			[]string{"0"},
			false,
		},
		"no label error": {
			[]string{
				s.codeId,
				fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", owner, keeper.RandomAccountAddress(s.T())),
			},
			false,
		},
		"no admin error": {
			[]string{
				s.codeId,
				fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", owner, keeper.RandomAccountAddress(s.T())),
				fmt.Sprintf("--label=%v", "TestContract"),
			},
			false,
		},
		"no sender error": {
			[]string{
				s.codeId,
				fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", owner, keeper.RandomAccountAddress(s.T())),
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
			},
			false,
		},
		"no instantiate params error": {
			[]string{
				s.codeId,
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			false,
		},
		"no exist codeId error": {
			[]string{
				"0",
				fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", owner, keeper.RandomAccountAddress(s.T())),
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.InstantiateContractCmd()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestExecuteContractCmd() {
	val := s.network.Validators[0]

	params := fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", s.network.Validators[0].Address.String(), keeper.RandomAccountAddress(s.T()))
	contractAddr := s.instantiate(s.codeId, params)

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid executeContract": {
			[]string{
				contractAddr,
				fmt.Sprintf("{\"release\":{}}"),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			true,
		},
		"wrong amount": {
			[]string{
				contractAddr,
				fmt.Sprintf("{\"release\":{}}"),
				fmt.Sprintf("--%s=%v", "amount", "100"),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			false,
		},
		"wrong param": {
			[]string{
				contractAddr,
				fmt.Sprintf("{release:{}}"),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			false,
		},
		"no contract address": {
			[]string{
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
			},
			false,
		},
		"no sender": {
			[]string{
				contractAddr,
				fmt.Sprintf("{\"release\":{}}"),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.ExecuteContractCmd()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}

func (s *IntegrationTestSuite) TestStoreCodeAndInstantiateContractCmd() {
	val := s.network.Validators[0]
	owner := val.Address.String()

	wasmPath := "../../keeper/testdata/hackatom.wasm"
	_, err := ioutil.ReadFile(wasmPath)
	s.Require().NoError(err)

	params := fmt.Sprintf("{\"verifier\": \"%s\", \"beneficiary\": \"%s\"}", s.network.Validators[0].Address.String(), keeper.RandomAccountAddress(s.T()))

	testCases := map[string]struct {
		args  []string
		valid bool
	}{
		"valid storeCodeAndInstantiateContract": {
			[]string{
				wasmPath,
				params,
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
				fmt.Sprintf("--%s=%v", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=%v", flags.FlagGas, 1600000),
			},
			true,
		},
		"wrong args count": {
			[]string{"0"},
			false,
		},
		"no label error": {
			[]string{
				wasmPath,
				params,
			},
			false,
		},
		"no sender error": {
			[]string{
				wasmPath,
				params,
				fmt.Sprintf("--label=%v", "TestContract"),
				fmt.Sprintf("--admin=%v", owner),
			},
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.StoreCodeAndInstantiateContractCmd()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, append(tc.args, commonArgs...))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			var res sdk.TxResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
			s.Require().EqualValues(0, res.Code, out.String())
		})
	}
}
