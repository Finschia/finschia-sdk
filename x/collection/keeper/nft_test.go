package keeper_test

import (
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestAttach() {
	testCases := map[string]struct {
		contractID string
		subject    string
		target     string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			target:     collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit),
		},
		"not owner of subject": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			target:     collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrTokenNotOwnedBy,
		},
		"target not found": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			target:     collection.NewNFTID(s.nftClassID, s.numNFTs*3+1),
			err:        collection.ErrTokenNotExist,
		},
		"result exceeds the limit": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+2),
			target:     collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit),
			err:        collection.ErrCompositionTooDeep,
		},
		"not owner of target": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, collection.DefaultDepthLimit+1),
			target:     collection.NewNFTID(s.nftClassID, s.numNFTs+1),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Attach(ctx, tc.contractID, s.customer, tc.subject, tc.target)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			parent, err := s.keeper.GetParent(ctx, tc.contractID, tc.subject)
			s.Require().NoError(err)
			s.Require().Equal(*parent, tc.target)
		})
	}
}

func (s *KeeperTestSuite) TestDetach() {
	testCases := map[string]struct {
		contractID string
		subject    string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, 2),
		},
		"subject not found": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.numNFTs*3+1),
			err:        collection.ErrTokenNotExist,
		},
		"subject has no parent": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, 1),
			err:        collection.ErrTokenNotAChild,
		},
		"not owner of subject": {
			contractID: s.contractID,
			subject:    collection.NewNFTID(s.nftClassID, s.numNFTs+2),
			err:        collection.ErrTokenNotOwnedBy,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Detach(ctx, tc.contractID, s.customer, tc.subject)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			parent, err := s.keeper.GetParent(ctx, tc.contractID, tc.subject)
			s.Require().Error(err)
			s.Require().Nil(parent)
		})
	}
}

func (s *KeeperTestSuite) TestGetNFT() {
	testCases := map[string]struct {
		tokenID string
		err     error
	}{
		"valid request": {
			tokenID: collection.NewNFTID(s.nftClassID, 1),
		},
		"not found (not existing nft id)": {
			tokenID: collection.NewNFTID("deadbeef", 1),
			err:     collection.ErrTokenNotExist,
		},
		"not found (not existing ft id)": {
			tokenID: collection.NewFTID("00bab10c"),
			err:     collection.ErrTokenNotExist,
		},
		"not found (existing nft class id)": {
			tokenID: collection.NewNFTID(s.nftClassID, 0),
			err:     collection.ErrTokenNotExist,
		},
		"not found (existing ft class id)": {
			tokenID: collection.NewNFTID(s.ftClassID, 0),
			err:     collection.ErrTokenNotNFT,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			token, err := s.keeper.GetNFT(s.ctx, s.contractID, tc.tokenID)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(token)
		})
	}
}
