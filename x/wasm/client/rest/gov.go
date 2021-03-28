package rest

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/line/lbm-sdk/v2/x/wasm/internal/types"
)

type StoreCodeProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`

	RunAs sdk.AccAddress `json:"run_as" yaml:"run_as"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte `json:"wasm_byte_code" yaml:"wasm_byte_code"`
	// Source is a valid absolute HTTPS URI to the contract's source code, optional
	Source string `json:"source" yaml:"source"`
	// Builder is a valid docker image name with tag, optional
	Builder string `json:"builder" yaml:"builder"`
	// InstantiatePermission to apply on contract creation, optional
	InstantiatePermission *types.AccessConfig `json:"instantiate_permission" yaml:"instantiate_permission"`
}

func (s StoreCodeProposalJSONReq) Content() gov.Content {
	return types.StoreCodeProposal{
		WasmProposal: types.WasmProposal{
			Title:       s.Title,
			Description: s.Description,
		},
		RunAs:                 s.RunAs,
		WASMByteCode:          s.WASMByteCode,
		Source:                s.Source,
		Builder:               s.Builder,
		InstantiatePermission: s.InstantiatePermission,
	}
}
func (s StoreCodeProposalJSONReq) GetProposer() sdk.AccAddress {
	return s.Proposer
}
func (s StoreCodeProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s StoreCodeProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func StoreCodeProposalHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_store_code",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req StoreCodeProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type InstantiateProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins      `json:"deposit" yaml:"deposit"`

	RunAs sdk.AccAddress `json:"run_as" yaml:"run_as"`
	// Admin is an optional address that can execute migrations
	Admin     sdk.AccAddress  `json:"admin,omitempty" yaml:"admin"`
	Code      uint64          `json:"code_id" yaml:"code_id"`
	Label     string          `json:"label" yaml:"label"`
	InitMsg   json.RawMessage `json:"init_msg" yaml:"init_msg"`
	InitFunds sdk.Coins       `json:"init_funds" yaml:"init_funds"`
}

func (s InstantiateProposalJSONReq) Content() gov.Content {
	return types.InstantiateContractProposal{
		WasmProposal: types.WasmProposal{Title: s.Title, Description: s.Description},
		RunAs:        s.RunAs,
		Admin:        s.Admin,
		CodeID:       s.Code,
		Label:        s.Label,
		InitMsg:      s.InitMsg,
		InitFunds:    s.InitFunds,
	}
}
func (s InstantiateProposalJSONReq) GetProposer() sdk.AccAddress {
	return s.Proposer
}
func (s InstantiateProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s InstantiateProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func InstantiateProposalHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_instantiate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req InstantiateProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type MigrateProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins      `json:"deposit" yaml:"deposit"`

	Contract   sdk.AccAddress  `json:"contract" yaml:"contract"`
	Code       uint64          `json:"code_id" yaml:"code_id"`
	MigrateMsg json.RawMessage `json:"msg" yaml:"msg"`
	// RunAs is the role that is passed to the contract's environment
	RunAs sdk.AccAddress `json:"run_as" yaml:"run_as"`
}

func (s MigrateProposalJSONReq) Content() gov.Content {
	return types.MigrateContractProposal{
		WasmProposal: types.WasmProposal{Title: s.Title, Description: s.Description},
		Contract:     s.Contract,
		CodeID:       s.Code,
		MigrateMsg:   s.MigrateMsg,
		RunAs:        s.RunAs,
	}
}
func (s MigrateProposalJSONReq) GetProposer() sdk.AccAddress {
	return s.Proposer
}
func (s MigrateProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s MigrateProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func MigrateProposalHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_migrate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req MigrateProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type UpdateAdminJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins      `json:"deposit" yaml:"deposit"`

	NewAdmin sdk.AccAddress `json:"new_admin" yaml:"new_admin"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
}

func (s UpdateAdminJSONReq) Content() gov.Content {
	return types.UpdateAdminProposal{
		WasmProposal: types.WasmProposal{Title: s.Title, Description: s.Description},
		Contract:     s.Contract,
		NewAdmin:     s.NewAdmin,
	}
}
func (s UpdateAdminJSONReq) GetProposer() sdk.AccAddress {
	return s.Proposer
}
func (s UpdateAdminJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s UpdateAdminJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func UpdateContractAdminProposalHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_update_admin",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req UpdateAdminJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type ClearAdminJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins      `json:"deposit" yaml:"deposit"`

	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
}

func (s ClearAdminJSONReq) Content() gov.Content {
	return types.ClearAdminProposal{
		WasmProposal: types.WasmProposal{Title: s.Title, Description: s.Description},
		Contract:     s.Contract,
	}
}
func (s ClearAdminJSONReq) GetProposer() sdk.AccAddress {
	return s.Proposer
}
func (s ClearAdminJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s ClearAdminJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func ClearContractAdminProposalHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wasm_clear_admin",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req ClearAdminJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type wasmProposalData interface {
	Content() gov.Content
	GetProposer() sdk.AccAddress
	GetDeposit() sdk.Coins
	GetBaseReq() rest.BaseReq
}

func toStdTxResponse(cliCtx context.CLIContext, w http.ResponseWriter, data wasmProposalData) {
	msg := gov.NewMsgSubmitProposal(data.Content(), data.GetDeposit(), data.GetProposer())
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	baseReq := data.GetBaseReq().Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}
	utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
}
