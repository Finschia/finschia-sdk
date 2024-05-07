package types

import (
	"errors"
)

var QueryParamToRole = map[string]Role{
	"unspecified": 0,
	"guardian":    1,
	"operator":    2,
	"judge":       3,
}

func IsValidRole(role Role) error {
	switch role {
	case RoleGuardian, RoleOperator, RoleJudge:
		return nil
	}

	return errors.New("unsupported role")
}

func IsValidVoteOption(option VoteOption) error {
	switch option {
	case OptionYes, OptionNo:
		return nil
	}

	return errors.New("unsupported vote option")
}

func IsValidBridgeStatus(status BridgeStatus) error {
	switch status {
	case StatusActive, StatusInactive:
		return nil
	}

	return errors.New("unsupported bridge status")
}
