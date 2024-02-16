package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/collection"
)

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
