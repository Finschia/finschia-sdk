package client_test

import (
	"strconv"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/client/flags"
)

func TestPaginate(t *testing.T) {
	testCases := []struct {
		name                           string
		numObjs, page, limit, defLimit int
		expectedStart, expectedEnd     int
	}{
		{
			"all objects in a single page",
			100, 1, 100, 100,
			0, 100,
		},
		{
			"page one of three",
			75, 1, 25, 100,
			0, 25,
		},
		{
			"page two of three",
			75, 2, 25, 100,
			25, 50,
		},
		{
			"page three of three",
			75, 3, 25, 100,
			50, 75,
		},
		{
			"end is greater than total number of objects",
			75, 2, 50, 100,
			50, 75,
		},
		{
			"fallback to default limit",
			75, 5, 0, 10,
			40, 50,
		},
		{
			"invalid start page",
			75, 4, 25, 100,
			-1, -1,
		},
		{
			"invalid zero start page",
			75, 0, 25, 100,
			-1, -1,
		},
		{
			"invalid negative start page",
			75, -1, 25, 100,
			-1, -1,
		},
		{
			"invalid default limit",
			75, 2, 0, -10,
			-1, -1,
		},
	}

	for i, tc := range testCases {
		i, tc := i, tc
		t.Run(tc.name, func(t *testing.T) {
			start, end := client.Paginate(tc.numObjs, tc.page, tc.limit, tc.defLimit)
			require.Equal(t, tc.expectedStart, start, "invalid result; test case #%d", i)
			require.Equal(t, tc.expectedEnd, end, "invalid result; test case #%d", i)
		})
	}
}

func TestReadPageRequest(t *testing.T) {

	testCases := []struct {
		name                string
		pageKey             string
		offset, limit, page int
		countTotal          bool
		ok                  bool
	}{
		{
			"use page ok",
			"page key",
			0, 100, 10,
			true,
			true,
		},
		{
			"use offset ok",
			"page key",
			10, 100, 0,
			true,
			true,
		},
		{
			"page and offset cannot be used together",
			"page key",
			100, 100, 10,
			true,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := pflag.NewFlagSet("test flag set", pflag.ContinueOnError)
			flagSet.String(flags.FlagPageKey, "default page key", "page key")
			flagSet.Uint64(flags.FlagOffset, 0, "offset")
			flagSet.Uint64(flags.FlagLimit, 0, "limit")
			flagSet.Uint64(flags.FlagPage, 0, "page")
			flagSet.Bool(flags.FlagCountTotal, false, "count total")

			err := flagSet.Set(flags.FlagPageKey, tc.pageKey)
			err = flagSet.Set(flags.FlagOffset, strconv.Itoa(tc.offset))
			err = flagSet.Set(flags.FlagLimit, strconv.Itoa(tc.limit))
			err = flagSet.Set(flags.FlagPage, strconv.Itoa(tc.page))
			err = flagSet.Set(flags.FlagCountTotal, strconv.FormatBool(tc.countTotal))

			pr, err := client.ReadPageRequest(flagSet)
			if tc.ok {
				require.NoError(t, err)
				require.NotNil(t, pr)
			} else {
				require.Error(t, err)
			}
		})
	}
}
