package client

import (
	"github.com/spf13/pflag"

	"github.com/line/lbm-sdk/v2/client/flags"
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
	"github.com/line/lbm-sdk/v2/types/query"
)

// Paginate returns the correct starting and ending index for a paginated query,
// given that client provides a desired page and limit of objects and the handler
// provides the total number of objects. The start page is assumed to be 1-indexed.
// If the start page is invalid, non-positive values are returned signaling the
// request is invalid; it returns non-positive values if limit is non-positive and
// defLimit is negative.
func Paginate(numObjs, page, limit, defLimit int) (start, end int) {
	if page <= 0 {
		// invalid start page
		return -1, -1
	}

	// fallback to default limit if supplied limit is invalid
	if limit <= 0 {
		if defLimit < 0 {
			// invalid default limit
			return -1, -1
		}
		limit = defLimit
	}

	start = (page - 1) * limit
	end = limit + start

	if end >= numObjs {
		end = numObjs
	}

	if start >= numObjs {
		// page is out of bounds
		return -1, -1
	}

	return start, end
}

// ReadPageRequest reads and builds the necessary page request flags for pagination.
func ReadPageRequest(flagSet *pflag.FlagSet) (*query.PageRequest, error) {
	pageKey, _ := flagSet.GetString(flags.FlagPageKey)
	offset, _ := flagSet.GetUint64(flags.FlagOffset)
	limit, _ := flagSet.GetUint64(flags.FlagLimit)
	countTotal, _ := flagSet.GetBool(flags.FlagCountTotal)
	page, _ := flagSet.GetUint64(flags.FlagPage)

	return NewPageRequest(pageKey, offset, limit, page, countTotal)
}

func NewPageRequest(pageKey string, offset, limit, page uint64, countTotal bool) (*query.PageRequest, error) {
	if page > 1 && offset > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "page and offset cannot be used together")
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return &query.PageRequest{
		Key:        []byte(pageKey),
		Offset:     offset,
		Limit:      limit,
		CountTotal: countTotal,
	}, nil
}
