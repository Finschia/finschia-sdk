package ibctesting_test

import (
	"testing"

	octypes "github.com/line/ostracon/types"
	"github.com/stretchr/testify/require"

	ibctesting "github.com/line/lbm-sdk/x/ibc/testing"
	"github.com/line/lbm-sdk/x/ibc/testing/mock"
)

func TestCreateSortedSignerArray(t *testing.T) {
	privVal1 := mock.NewPV()
	pubKey1, err := privVal1.GetPubKey()
	require.NoError(t, err)

	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	require.NoError(t, err)

	validator1 := ibctesting.NewTestValidator(pubKey1, 1)
	validator2 := ibctesting.NewTestValidator(pubKey2, 2)

	expected := []octypes.PrivValidator{privVal2, privVal1}

	actual := ibctesting.CreateSortedSignerArray(privVal1, privVal2, validator1, validator2)
	require.Equal(t, expected, actual)

	// swap order
	actual = ibctesting.CreateSortedSignerArray(privVal2, privVal1, validator2, validator1)
	require.Equal(t, expected, actual)

	// smaller address
	validator1.Address = []byte{1}
	validator2.Address = []byte{2}
	validator2.VotingPower = 1

	expected = []octypes.PrivValidator{privVal1, privVal2}

	actual = ibctesting.CreateSortedSignerArray(privVal1, privVal2, validator1, validator2)
	require.Equal(t, expected, actual)

	// swap order
	actual = ibctesting.CreateSortedSignerArray(privVal2, privVal1, validator2, validator1)
	require.Equal(t, expected, actual)
}
