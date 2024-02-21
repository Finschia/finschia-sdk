package types

import (
	"encoding/json"

	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"

	types2 "github.com/Finschia/finschia-sdk/crypto/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

type ZKAuthVerifier struct {
	VerifyKey []byte
}

func NewZKAuthVerifier(vk []byte) ZKAuthVerifier {
	return ZKAuthVerifier{
		VerifyKey: vk,
	}
}

func (v *ZKAuthVerifier) Verify(proof types.ZKProof) error {
	return verifier.VerifyGroth16(proof, v.VerifyKey)
}

func VerifyZKAuthSignature(ctx sdk.Context, zkv ZKAuthVerifier, ephPubKey types2.PubKey, msg *MsgExecution) error {
	// check max block height
	if msg.ZkAuthSignature.MaxBlockHeight < ctx.BlockHeader().Height {
		return sdkerrors.Wrap(ErrInvalidZKAuthSignature, "The permitted block height was exceeded.")
	}

	var proofData types.ProofData
	err := json.Unmarshal(msg.ZkAuthSignature.ZkAuthInputs.ProofPoints, &proofData)
	if err != nil {
		return nil
	}

	// get OAuth pubKey

	// calculate all input hash
	allInputHash, err := msg.ZkAuthSignature.ZkAuthInputs.CalculateAllInputsHash(ephPubKey.Bytes(), []byte{}, msg.ZkAuthSignature.MaxBlockHeight)
	if err != nil {
		return err
	}

	// verify
	proof := types.ZKProof{
		Proof:      &proofData,
		PubSignals: []string{allInputHash.String()},
	}

	return zkv.Verify(proof)
}
