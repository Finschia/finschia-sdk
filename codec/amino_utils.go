package codec

import (
	"errors"
	"fmt"
	"github.com/tendermint/go-amino"
)

func decodeFieldNumberAndTyp3(bz []byte) (num uint32, typ amino.Typ3, n int, err error) {
	// Read uvarint value.
	var value64 = uint64(0)
	value64, n, err = amino.DecodeUvarint(bz)
	if err != nil {
		return
	}

	// Decode first typ3 byte.
	typ = amino.Typ3(value64 & 0x07)

	// Decode num.
	var num64 uint64
	num64 = value64 >> 3
	if num64 > (1<<29 - 1) {
		err = fmt.Errorf("invalid field num %v", num64)
		return
	}
	num = uint32(num64)
	return
}

func CheckFieldNumberAndTyp3(bz []byte, expectedFnum uint32, expectedTyp amino.Typ3) (int, error) {
	fnum, typ, n, err := decodeFieldNumberAndTyp3(bz)
	if err != nil {
		return 0, err
	}
	if fnum != expectedFnum {
		return 0, nil
	}
	if typ != expectedTyp {
		return 0, errors.New(fmt.Sprintf("expected field type %v, got %v", expectedTyp, typ))
	}
	return n, nil
}

// CONTRACT: by the time this is called, len(bz) >= _n
// Returns true so you can write one-liners.
func Slide(bz *[]byte, n *int, _n int) bool {
	if _n < 0 || _n > len(*bz) {
		panic(fmt.Sprintf("impossible slide: len:%v _n:%v", len(*bz), _n))
	}
	*bz = (*bz)[_n:]
	if n != nil {
		*n += _n
	}
	return true
}
