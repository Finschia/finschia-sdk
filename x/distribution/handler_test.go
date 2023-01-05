package distribution

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	mocktypes "github.com/line/lbm-sdk/x/distribution/mocks/types"
	"github.com/line/lbm-sdk/x/distribution/types"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

type handlerTestSuite struct {
	suite.Suite
	mockSrv *mocktypes.MsgServer
	cut     func(sdk.Context, sdk.Msg) (*sdk.Result, error)
}

func (s *handlerTestSuite) SetupTest() {
	s.mockSrv = &mocktypes.MsgServer{}
	s.cut = newHandler(s.mockSrv)
}

func (s *handlerTestSuite) TearDownTest() {
	s.mockSrv.AssertExpectations(s.T())
}

func (s *handlerTestSuite) TestMsgSetWithdrawAddress() {
	// Arrange
	ctx := sdk.NewContext(nil, ocproto.Header{}, false, nil)
	withdrawAddress := &types.MsgSetWithdrawAddress{}
	s.mockSrv.On("SetWithdrawAddress", sdk.WrapSDKContext(ctx), mock.Anything).Return(nil, nil)

	// Act
	_, err := s.cut(ctx, withdrawAddress)

	// Assert
	if err != nil {
		assert.Fail(s.T(), err.Error())
	}
}

func (s *handlerTestSuite) TestMsgWithdrawDelegatorReward() {
	// Arrange
	ctx := sdk.NewContext(nil, ocproto.Header{}, false, nil)
	msgWithdrawDelegatorReward := &types.MsgWithdrawDelegatorReward{}
	s.mockSrv.On("WithdrawDelegatorReward", sdk.WrapSDKContext(ctx), mock.Anything).Return(nil, nil)

	// Act
	_, err := s.cut(ctx, msgWithdrawDelegatorReward)

	// Assert
	require.NoError(s.T(), err)
}

func (s *handlerTestSuite) TestMsgWithdrawValidatorCommission() {
	// Arrange
	ctx := sdk.NewContext(nil, ocproto.Header{}, false, nil)
	msgWithdrawValidatorCommission := &types.MsgWithdrawValidatorCommission{}
	s.mockSrv.On("WithdrawValidatorCommission", sdk.WrapSDKContext(ctx), mock.Anything).Return(nil, nil)

	// Act
	_, err := s.cut(ctx, msgWithdrawValidatorCommission)

	// Assert
	require.NoError(s.T(), err)
}

func (s *handlerTestSuite) TestMsgFundCommunityPool() {
	// Arrange
	ctx := sdk.NewContext(nil, ocproto.Header{}, false, nil)
	msgFundCommunityPool := &types.MsgFundCommunityPool{}
	s.mockSrv.On("FundCommunityPool", sdk.WrapSDKContext(ctx), mock.Anything).Return(nil, nil)

	// Act
	_, err := s.cut(ctx, msgFundCommunityPool)

	// Assert
	require.NoError(s.T(), err)
}

func (s *handlerTestSuite) TestThrowErrorForUnknownMsgType() {
	// Arrange
	ctx := sdk.NewContext(nil, ocproto.Header{}, false, nil)

	// Act
	_, err := s.cut(ctx, nil)

	// Assert
	require.ErrorIs(s.T(), err, sdkerrors.ErrUnknownRequest)
}
