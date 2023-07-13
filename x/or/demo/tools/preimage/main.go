package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
)

func main() {
	preimagePath := os.Args[1]

	channel := preimage.ClientPreimageChannel()
	srv := preimage.NewOracleServer(channel)
	preimages, err := preimages(preimagePath)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := srv.NextPreimageRequest(func(key [32]byte) ([]byte, error) {
				dat, ok := preimages[key]
				if !ok {
					return nil, fmt.Errorf("cannot find %s", key)
				}
				return dat, nil
			})
			if err != nil {
				panic(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-c
		if sig.String() == "interrupt" {
			return
		}
	}
}

func preimages(path string) (map[[32]byte][]byte, error) {
	if path == "" {
		return nil, errors.New("no path specified")
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer f.Close()

	var jsonMap map[string]string
	if err := json.NewDecoder(f).Decode(&jsonMap); err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", path, err)
	}

	preimages := map[[32]byte][]byte{}

	for k := range jsonMap {
		kb, err := hex.DecodeString(k)
		if err != nil {
			return nil, err
		}
		vb, err := hex.DecodeString(jsonMap[k])
		if err != nil {
			return nil, err
		}
		preimages[[32]byte(kb)] = vb
	}

	return preimages, nil
}

type sha256Key [32]byte

const sha256KeyType = 100

func (s sha256Key) PreimageKey() (out [32]byte) {
	out = s                      // copy the keccak hash
	out[0] = byte(sha256KeyType) // apply prefix
	return
}

func (s sha256Key) String() string {
	return "0x" + hex.EncodeToString(s[:])
}

func (s sha256Key) TerminalString() string {
	return "0x" + hex.EncodeToString(s[:])
}
