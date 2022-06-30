package keeper_test

import (
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestAttach() {
	testCases := map[string]struct {
		contractID string
		subject    string
		target     string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain+1),
			target:     collection.NewNFTID(s.nftClassID, 1),
			valid:      true,
		},
		"not owner of subject": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain*2+1),
			target:     collection.NewNFTID(s.nftClassID, 1),
		},
		"target not found": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain+1),
			target:     collection.NewNFTID(s.nftClassID, s.lenChain*6+1),
		},
		"target's depth exceeds the limit": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain+1),
			target:     collection.NewNFTID(s.nftClassID, s.lenChain),
		},
		"not owner of target": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain+1),
			target:     collection.NewNFTID(s.nftClassID, s.lenChain*2+1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.Attach(s.ctx, tc.contractID, s.customer, tc.subject, tc.target)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			children := s.keeper.GetChildren(s.ctx, tc.contractID, tc.target)
			s.Require().Equal(2, len(children))
		})
	}
}

func (s *KeeperTestSuite) TestDetach() {
	testCases := map[string]struct {
		contractID string
		subject    string
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, 2),
			valid:      true,
		},
		"subject not found": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain*6+1),
		},
		"subject has no parent": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, 1),
		},
		"parent's depth exceeds the limit": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain),
		},
		"not owner of subject": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.lenChain*2+2),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.Detach(s.ctx, tc.contractID, s.customer, tc.subject)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			parent, err := s.keeper.GetParent(s.ctx, tc.contractID, tc.subject)
			s.Require().Error(err)
			s.Require().Nil(parent)
		})
	}
}
