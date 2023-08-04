package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum-optimism/optimism/cannon/cmd"
	"github.com/ethereum-optimism/optimism/cannon/mipsevm"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/profile"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/urfave/cli/v2"
)

type Proof struct {
	Step uint64 `json:"step"`

	Pre  tmbytes.HexBytes `json:"pre"`
	Post tmbytes.HexBytes `json:"post"`

	Witness Witness `json:"witness"`
}

type Witness struct {
	State          tmbytes.HexBytes `json:"state"`
	MemProof       tmbytes.HexBytes `json:"mem_proof"`
	PreimageKey    tmbytes.HexBytes `json:"preimage_key"`
	PreimageValue  tmbytes.HexBytes `json:"preimage_value"`
	PreimageOffset uint32           `json:"preimage_offset"`
}

func WitnessFrom(w mipsevm.StepWitness) Witness {
	return Witness{
		State:          w.State,
		MemProof:       w.MemProof,
		PreimageKey:    w.PreimageKey[:],
		PreimageValue:  w.PreimageValue,
		PreimageOffset: w.PreimageOffset,
	}
}

type rawHint string

func (rh rawHint) Hint() string {
	return string(rh)
}

type rawKey [32]byte

func (rk rawKey) PreimageKey() [32]byte {
	return rk
}

type ProcessPreimageOracle struct {
	pCl *preimage.OracleClient
	hCl *preimage.HintWriter
	cmd *exec.Cmd
}

func (p *ProcessPreimageOracle) Hint(v []byte) {
	if p.hCl == nil { // no hint processor
		return
	}
	p.hCl.Hint(rawHint(v))
}

func (p *ProcessPreimageOracle) GetPreimage(k [32]byte) []byte {
	if p.pCl == nil {
		panic("no pre-image retriever available")
	}
	return p.pCl.Get(rawKey(k))
}

func (p *ProcessPreimageOracle) Start() error {
	if p.cmd == nil {
		return nil
	}
	return p.cmd.Start()
}

func (p *ProcessPreimageOracle) Close() error {
	if p.cmd == nil {
		return nil
	}
	_ = p.cmd.Process.Signal(os.Interrupt)
	// Go 1.20 feature, to introduce later
	//p.cmd.WaitDelay = time.Second * 10
	err := p.cmd.Wait()
	if err, ok := err.(*exec.ExitError); ok {
		if err.Success() {
			return nil
		}
	}
	return err
}

func NewProcessPreimageOracle(name string, args []string) (*ProcessPreimageOracle, error) {
	if name == "" {
		return &ProcessPreimageOracle{}, nil
	}

	pClientRW, pOracleRW, err := preimage.CreateBidirectionalChannel()
	if err != nil {
		return nil, err
	}
	hClientRW, hOracleRW, err := preimage.CreateBidirectionalChannel()
	if err != nil {
		return nil, err
	}

	ecmd := exec.Command(name, args...) // nosemgrep
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	ecmd.ExtraFiles = []*os.File{
		hOracleRW.Reader(),
		hOracleRW.Writer(),
		pOracleRW.Reader(),
		pOracleRW.Writer(),
	}
	out := &ProcessPreimageOracle{
		pCl: preimage.NewOracleClient(pClientRW),
		hCl: preimage.NewHintWriter(hClientRW),
		cmd: ecmd,
	}
	return out, nil
}

func Run(ctx *cli.Context) error {
	if ctx.Bool(cmd.RunPProfCPU.Name) {
		defer profile.Start(profile.NoShutdownHook, profile.ProfilePath("."), profile.CPUProfile).Stop()
	}

	state, err := loadJSON[mipsevm.State](ctx.Path(cmd.RunInputFlag.Name))
	if err != nil {
		return err
	}

	l := cmd.Logger(os.Stderr, log.LvlInfo)
	outLog := &mipsevm.LoggingWriter{Name: "program std-out", Log: l}
	errLog := &mipsevm.LoggingWriter{Name: "program std-err", Log: l}

	// split CLI args after first '--'
	args := ctx.Args().Slice()
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) == 0 {
		args = []string{""}
	}

	po, err := NewProcessPreimageOracle(args[0], args[1:])
	if err != nil {
		return fmt.Errorf("failed to create pre-image oracle process: %w", err)
	}
	if err := po.Start(); err != nil {
		return fmt.Errorf("failed to start pre-image oracle server: %w", err)
	}
	defer func() {
		if err := po.Close(); err != nil {
			l.Error("failed to close pre-image server", "err", err)
		}
	}()

	stopAt := ctx.Generic(cmd.RunStopAtFlag.Name).(*cmd.StepMatcherFlag).Matcher()
	proofAt := ctx.Generic(cmd.RunProofAtFlag.Name).(*cmd.StepMatcherFlag).Matcher()
	snapshotAt := ctx.Generic(cmd.RunSnapshotAtFlag.Name).(*cmd.StepMatcherFlag).Matcher()
	infoAt := ctx.Generic(cmd.RunInfoAtFlag.Name).(*cmd.StepMatcherFlag).Matcher()

	var meta *mipsevm.Metadata
	if metaPath := ctx.Path(cmd.RunMetaFlag.Name); metaPath == "" {
		l.Info("no metadata file specified, defaulting to empty metadata")
		meta = &mipsevm.Metadata{Symbols: nil} // provide empty metadata by default
	} else {
		if m, err := loadJSON[mipsevm.Metadata](metaPath); err != nil {
			return fmt.Errorf("failed to load metadata: %w", err)
		} else {
			meta = m
		}
	}

	us := mipsevm.NewInstrumentedState(state, po, outLog, errLog)
	proofFmt := ctx.String(cmd.RunProofFmtFlag.Name)
	snapshotFmt := ctx.String(cmd.RunSnapshotFmtFlag.Name)

	stepFn := us.Step
	if po.cmd != nil {
		stepFn = cmd.Guard(po.cmd.ProcessState, stepFn)
	}

	start := time.Now()
	startStep := state.Step

	// avoid symbol lookups every instruction by preparing a matcher func
	sleepCheck := meta.SymbolMatcher("runtime.notesleep")

	for !state.Exited {
		if state.Step%100 == 0 { // don't do the ctx err check (includes lock) too often
			if err := ctx.Context.Err(); err != nil {
				return err
			}
		}

		step := state.Step

		if infoAt(state) {
			delta := time.Since(start)
			l.Info("processing",
				"step", step,
				"pc", mipsevm.HexU32(state.PC),
				"insn", mipsevm.HexU32(state.Memory.GetMemory(state.PC)),
				"ips", float64(step-startStep)/(float64(delta)/float64(time.Second)),
				"pages", state.Memory.PageCount(),
				"mem", state.Memory.Usage(),
				"name", meta.LookupSymbol(state.PC),
			)
		}

		if sleepCheck(state.PC) { // don't loop forever when we get stuck because of an unexpected bad program
			return fmt.Errorf("got stuck in Go sleep at step %d", step)
		}

		if stopAt(state) {
			break
		}

		if snapshotAt(state) {
			if err := writeJSON(fmt.Sprintf(snapshotFmt, step), state, false); err != nil {
				return fmt.Errorf("failed to write state snapshot: %w", err)
			}
		}

		if proofAt(state) {
			preStateHash := sha256.Sum256(state.EncodeWitness())
			witness, err := stepFn(true)
			if err != nil {
				return fmt.Errorf("failed at proof-gen step %d (PC: %08x): %w", step, state.PC, err)
			}
			postStateHash := sha256.Sum256(state.EncodeWitness())
			proof := &Proof{
				Step:    step,
				Pre:     preStateHash[:],
				Post:    postStateHash[:],
				Witness: WitnessFrom(*witness),
			}
			if err := writeJSON(fmt.Sprintf(proofFmt, step), proof, true); err != nil {
				return fmt.Errorf("failed to write proof data: %w", err)
			}
		} else {
			_, err = stepFn(false)
			if err != nil {
				return fmt.Errorf("failed at step %d (PC: %08x): %w", step, state.PC, err)
			}
		}
	}

	if err := writeJSON(ctx.Path(cmd.RunOutputFlag.Name), state, true); err != nil {
		return fmt.Errorf("failed to write state output: %w", err)
	}
	return nil
}

var RunCommand = &cli.Command{
	Name:        "run",
	Usage:       "Run VM step(s) and generate proof data to replicate onchain.",
	Description: "Run VM step(s) and generate proof data to replicate onchain. See flags to match when to output a proof, a snapshot, or to stop early.",
	Action:      Run,
	Flags: []cli.Flag{
		cmd.RunInputFlag,
		cmd.RunOutputFlag,
		cmd.RunProofAtFlag,
		cmd.RunProofFmtFlag,
		cmd.RunSnapshotAtFlag,
		cmd.RunSnapshotFmtFlag,
		cmd.RunStopAtFlag,
		cmd.RunMetaFlag,
		cmd.RunInfoAtFlag,
		cmd.RunPProfCPU,
	},
}
