package types

import (
	"errors"
)

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
