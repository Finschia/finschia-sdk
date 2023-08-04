package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/urfave/cli/v2"
)

var (
	ChallengePathFlag = &cli.PathFlag{
		Name:      "path",
		Usage:     "Path to challenge json file",
		TakesFile: true,
		Required:  true,
	}
)

func GetID(ctx *cli.Context) error {
	path := ctx.Path(ChallengePathFlag.Name)
	challenge, err := getChallenge(path)
	if err != nil {
		return err
	}
	fmt.Println(challenge.ID())
	return nil
}

var IDCommand = &cli.Command{
	Name:   "id",
	Usage:  "Get challenge id",
	Action: GetID,
	Flags: []cli.Flag{
		ChallengePathFlag,
	},
}

func GetIDFrom(ctx *cli.Context) error {
	rollupName := ctx.Args().Get(0)
	blockHeight, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}
	challenger := ctx.Args().Get(2)
	defender := ctx.Args().Get(3)

	challenge := types.Challenge{
		RollupName:  rollupName,
		BlockHeight: int64(blockHeight),
		Challenger:  challenger,
		Defender:    defender,
	}
	fmt.Println(challenge.ID())
	return nil
}

var IDFromCommand = &cli.Command{
	Name:   "id-from",
	Usage:  "Get challenge id from input. example, id-from  [rollup_name] [block_height] [challenger_address] [defender_address]",
	Action: GetIDFrom,
}

func IsSearching(ctx *cli.Context) error {
	path := ctx.Path(ChallengePathFlag.Name)
	challenge, err := getChallenge(path)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", challenge.IsSearching())
	return nil
}

var IsSearchingCommand = &cli.Command{
	Name:   "is-searching",
	Usage:  "Get flag if challenge is searching",
	Action: IsSearching,
	Flags: []cli.Flag{
		ChallengePathFlag,
	},
}

func GetSteps(ctx *cli.Context) error {
	path := ctx.Path(ChallengePathFlag.Name)
	challenge, err := getChallenge(path)
	if err != nil {
		return err
	}
	steps := challenge.GetSteps()
	for i := range steps {
		fmt.Printf("%d", steps[i])
		if len(steps) == i+1 {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" ")
		}
	}
	return nil
}

var StepsCommand = &cli.Command{
	Name:   "steps",
	Usage:  "Get steps to send state hashes",
	Action: GetSteps,
	Flags: []cli.Flag{
		ChallengePathFlag,
	},
}

func getChallenge(path string) (*types.Challenge, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	challengeStrMap := challengeStrMap{}
	if err := tmjson.Unmarshal(b, &challengeStrMap); err != nil {
		panic(err)
	}
	return &types.Challenge{
		RollupName:          challengeStrMap.RollupName,
		BlockHeight:         challengeStrMap.BlockHeight,
		Challenger:          challengeStrMap.Challenger,
		Defender:            challengeStrMap.Defender,
		L:                   challengeStrMap.L,
		R:                   challengeStrMap.R,
		AssertedStateHashes: challengeStrMap.getAssertedStateHashes(),
		DefendedStateHashes: challengeStrMap.getDefendedStateHashes(),
	}, nil
}

type challengeStrMap struct {
	RollupName          string            `json:"rollup_name,omitempty"`
	BlockHeight         int64             `json:"block_height,omitempty"`
	Challenger          string            `json:"challenger,omitempty"`
	Defender            string            `json:"defender,omitempty"`
	L                   uint64            `json:"l,omitempty"`
	R                   uint64            `json:"r,omitempty"`
	AssertedStateHashes map[string][]byte `json:"asserted_state_hashes,omitempty"`
	DefendedStateHashes map[string][]byte `json:"defended_state_hashes,omitempty"`
}

func (c challengeStrMap) getAssertedStateHashes() map[uint64][]byte {
	hashes := map[uint64][]byte{}
	for k := range c.AssertedStateHashes {
		key, err := strconv.Atoi(k)
		if err != nil {
			panic("something wrong")
		}
		hashes[uint64(key)] = c.AssertedStateHashes[k]
	}
	return hashes
}

func (c challengeStrMap) getDefendedStateHashes() map[uint64][]byte {
	hashes := map[uint64][]byte{}
	for k := range c.DefendedStateHashes {
		key, err := strconv.Atoi(k)
		if err != nil {
			panic("something wrong")
		}
		hashes[uint64(key)] = c.DefendedStateHashes[k]
	}
	return hashes
}
