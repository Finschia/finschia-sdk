package types

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Approval struct {
	TokenId         string         `json:"token_id"`
	ApprovedAddress sdk.AccAddress `json:"approved_address"`
}

func NewApproval(tokenId string, approvedAddress sdk.AccAddress) Approval {
	return Approval{
		TokenId:         tokenId,
		ApprovedAddress: approvedAddress,
	}
}

func (tokenApproval Approval) String() string {
	return fmt.Sprintf(`TokenId 			%s
	ApprovedAddress 			%s`,
		tokenApproval.TokenId, tokenApproval.ApprovedAddress)
}

type OperatorApprovals struct {
	Denom        string         `json:"denom"`
	OwnerAddress sdk.AccAddress `json:"owner_address"`
	Operators    Operators      `json:"operators"`
}

func NewOperatorApprovals(denom string, ownerAddress sdk.AccAddress, operators Operators) OperatorApprovals {
	return OperatorApprovals{
		Denom:        denom,
		OwnerAddress: ownerAddress,
		Operators:    operators,
	}
}

type Operators []sdk.AccAddress

func (operators Operators) DeleteOperator(operatorAddress sdk.AccAddress) (Operators, sdk.Error) {
	index := operators.find(operatorAddress)
	if index == -1 {
		return operators, ErrNotExistOperator(DefaultCodespace)
	}

	operators = append(operators[:index], operators[index+1:]...)

	return operators, nil
}

type QueryOperatorApproveParams struct {
	Denom           string
	OwnerAddress    sdk.AccAddress
	OperatorAddress sdk.AccAddress
}

func NewQueryOperatorApproveParams(denom string, ownerAddress, operatorAddress sdk.AccAddress) QueryOperatorApproveParams {
	return QueryOperatorApproveParams{
		Denom:           denom,
		OwnerAddress:    ownerAddress,
		OperatorAddress: operatorAddress,
	}
}

type TokenBalance struct {
	Address sdk.AccAddress
	Denom   string
	Amount  int
}

func NewTokenBalance(address sdk.AccAddress, denom string, amount int) TokenBalance {
	return TokenBalance{
		Address: address,
		Denom:   denom,
		Amount:  amount,
	}
}

// String follows stringer interface
func (tokenBalance TokenBalance) String() string {
	return fmt.Sprintf(`
	address: 				%s
	address: 				%s
	amount: 				%d`,
		tokenBalance.Address.String(),
		tokenBalance.Denom,
		tokenBalance.Amount,
	)
}

type TokenOwner struct {
	Denom   string
	TokenId string
	Address sdk.AccAddress
}

func NewTokenOwner(denom string, tokenId string, address sdk.AccAddress) TokenOwner {
	return TokenOwner{
		Denom:   denom,
		TokenId: tokenId,
		Address: address,
	}
}

// String follows stringer interface
func (tokenOwner TokenOwner) String() string {
	return fmt.Sprintf(`
	Denom: 			    %s,
	TokenId: 			%s,
	Owner: 				%s`,
		tokenOwner.Denom,
		tokenOwner.TokenId,
		tokenOwner.Address,
	)
}

func (operators Operators) Find(operatorAddress sdk.AccAddress) bool {
	index := operators.find(operatorAddress)
	if index == -1 {
		return false
	}
	return true
}

func (operatorApprovals OperatorApprovals) String() string {
	return fmt.Sprintf(`Denom: 			%s
	OwnerAddress: 			%s
	operators:        	%s`,
		operatorApprovals.Denom,
		operatorApprovals.OwnerAddress,
		operatorApprovals.Operators,
	)
}

func (operators Operators) find(el sdk.AccAddress) (idx int) {
	if len(operators) == 0 {
		return -1
	}

	operators.Sort()

	midIdx := len(operators) / 2
	stringArrayEl := operators[midIdx]

	switch {
	case el.String() < stringArrayEl.String():
		return operators[:midIdx].find(el)
	case stringArrayEl.Equals(el):
		return midIdx
	default:
		return operators[midIdx+1:].find(el)
	}
}

func (operators Operators) Len() int { return len(operators) }
func (operators Operators) Less(i, j int) bool {
	return strings.Compare(operators[i].String(), operators[j].String()) == -1
}
func (operators Operators) Swap(i, j int) {
	operators[i], operators[j] = operators[j], operators[i]
}
func (operators Operators) Sort() Operators {
	sort.Sort(operators)
	return operators
}

type ApprovedForAllResult struct {
	Found bool
}

func NewApprovedForAllResult(found bool) ApprovedForAllResult {
	return ApprovedForAllResult{
		Found: found,
	}
}

// String follows stringer interface
func (approvedForAllResult ApprovedForAllResult) String() string {
	return fmt.Sprintf(`
	found: 				%t`,
		approvedForAllResult.Found,
	)
}
