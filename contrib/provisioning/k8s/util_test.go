package k8s

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const checkMark = "\u2713"
const ballotX = "\u2717"
const failMsg = "The code did not panic" + ballotX
const strParam = "test"
const intParam = 1
const boolParam = true

var strArrParam = []string{strParam}

func TestErrCheckedParam(t *testing.T) {
	t.Log("Panic test", checkMark)
	{
		expectedError := fmt.Errorf("flag accessed but not defined: %s",
			"param")

		Panics(t, func() {
			ErrCheckedStrParam(strParam, expectedError)
		}, failMsg)

		Panics(t, func() {
			ErrCheckedIntParam(intParam, expectedError)
		}, failMsg)

		Panics(t, func() {
			ErrCheckedBoolParam(boolParam, expectedError)
		}, failMsg)
		Panics(t, func() {
			ErrCheckedStrArrParam(strArrParam, expectedError)
		}, failMsg)
	}
	t.Log("Not panic test", checkMark)
	{
		require.Equal(t, ErrCheckedStrParam(strParam, nil), strParam, failMsg)
		require.Equal(t, ErrCheckedIntParam(intParam, nil), intParam, failMsg)
		require.Equal(t, ErrCheckedBoolParam(boolParam, nil), boolParam, failMsg)
		require.Equal(t, ErrCheckedStrArrParam(strArrParam, nil), strArrParam, failMsg)
	}
}

func TestDefParam(t *testing.T) {
	t.Log("Test newVal", checkMark)
	{
		expectedStr := "expectedStr"
		var confField string
		require.Equal(t, DefIfEmpty(&confField, strParam, expectedStr), expectedStr, failMsg)
		require.Equal(t, confField, expectedStr, failMsg)
		expectedInt := intParam + 1
		expectedConfField := fmt.Sprintf(listenLoopbackIngressPortTemplate, expectedInt)
		require.Equal(t, DefFormatSetIfLTEZero(&confField, listenLoopbackIngressPortTemplate, intParam, expectedInt),
			expectedInt, failMsg)
		require.Equal(t, confField, expectedConfField, failMsg)
	}

	t.Log("Test defVal", checkMark)
	{
		var confField string
		require.Equal(t, DefIfEmpty(&confField, strParam, ""), strParam, failMsg)
		require.Equal(t, confField, strParam, failMsg)
		expectedConfField := fmt.Sprintf(listenLoopbackIngressPortTemplate, intParam)
		require.Equal(t, DefFormatSetIfLTEZero(&confField, listenLoopbackIngressPortTemplate, intParam, -1),
			intParam, failMsg)
		require.Equal(t, confField, expectedConfField, failMsg)
	}
}

func TestRandomHash(t *testing.T) {
	hash, e := RandomHash()
	require.Equal(t, e, nil)
	require.Equal(t, hash.Size(), 32)
}
