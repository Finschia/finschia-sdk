package v1

import (

	types3 "github.com/line/lbm-sdk/x/gov/types"
	yaml "gopkg.in/yaml.v2"

	"github.com/line/lbm-sdk/codec/types"
)

func ProposalToProposalV1(org types3.Proposal) ProposalV1 {
	p := ProposalV1{
		ProposalId:       org.ProposalId,
		Status:           org.Status,
		FinalTallyResult: org.FinalTallyResult,
		TotalDeposit:     org.TotalDeposit,
		SubmitTime:       org.SubmitTime,
		DepositEndTime:   org.DepositEndTime,
		VotingStartTime:  org.VotingStartTime,
		VotingEndTime:    org.VotingEndTime,
		Content:          org.Content,
	}

	return p
}

func ProposalV1ToProposal(v1 ProposalV1) types3.Proposal {
	p := types3.Proposal{
		ProposalId:       v1.ProposalId,
		Content:          v1.Content,
		Status:           v1.Status,
		FinalTallyResult: v1.FinalTallyResult,
		SubmitTime:       v1.SubmitTime,
		DepositEndTime:   v1.DepositEndTime,
		TotalDeposit:     v1.TotalDeposit,
		VotingStartTime:  v1.VotingStartTime,
		VotingEndTime:    v1.VotingEndTime,
		NewFieldSample:   0,
	}

	return p
}

// String implements stringer interface
func (p ProposalV1) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p ProposalV1) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var content types3.Content
	return unpacker.UnpackAny(p.Content, &content)
}
