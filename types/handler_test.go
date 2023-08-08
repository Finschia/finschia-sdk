package types_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-rdk/tests/mocks"
	sdk "github.com/Finschia/finschia-rdk/types"
)

type handlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (s *handlerTestSuite) SetupSuite() {
	s.T().Parallel()
}

func (s *handlerTestSuite) TestChainAnteDecorators() {
	// test panic
	s.Require().Nil(sdk.ChainAnteDecorators([]sdk.AnteDecorator{}...))

	ctx, tx := sdk.Context{}, sdk.Tx(nil)
	mockCtrl := gomock.NewController(s.T())
	mockAnteDecorator1 := mocks.NewMockAnteDecorator(mockCtrl)
	mockAnteDecorator1.EXPECT().AnteHandle(gomock.Eq(ctx), gomock.Eq(tx), true, gomock.Any()).Times(1)
	_, err := sdk.ChainAnteDecorators(mockAnteDecorator1)(ctx, tx, true)
	s.Require().NoError(err)

	called := []bool{false, false}
	testAnteDecorator1 := TestAnteDecorator{suite: s, ctx: ctx, tx: tx, simulate: true, called: &called[0]}
	testAnteDecorator2 := TestAnteDecorator{suite: s, ctx: ctx, tx: tx, simulate: true, called: &called[1]}

	_, err = sdk.ChainAnteDecorators(
		testAnteDecorator1,
		testAnteDecorator2)(ctx, tx, true)
	s.Require().NoError(err)
	s.Require().True(called[0])
	s.Require().True(called[1])
}

type TestAnteDecorator struct {
	suite    *handlerTestSuite
	ctx      sdk.Context
	tx       sdk.Tx
	simulate bool
	called   *bool
}

func (t TestAnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	t.suite.Require().Equal(t.ctx, ctx)
	t.suite.Require().Equal(t.tx, tx)
	t.suite.Require().Equal(t.simulate, simulate)
	*t.called = true
	return next(ctx, tx, simulate)
}
