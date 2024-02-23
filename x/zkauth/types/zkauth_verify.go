package types

import (
	"encoding/base64"
	"encoding/json"

	snarktypes "github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/pkg/errors"

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

func (v *ZKAuthVerifier) Verify(proof snarktypes.ZKProof) error {
	return verifier.VerifyGroth16(proof, v.VerifyKey)
}

func VerifyZKAuthSignature(ctx sdk.Context, zkv ZKAuthVerifier, jks *JWKs, ephPubKey []byte, msg *MsgExecution) error {
	// check max block height
	if msg.ZkAuthSignature.MaxBlockHeight < ctx.BlockHeader().Height {
		return sdkerrors.Wrap(ErrInvalidZKAuthSignature, "The permitted block height was exceeded.")
	}

	var proofData snarktypes.ProofData
	err := json.Unmarshal(msg.ZkAuthSignature.ZkAuthInputs.ProofPoints, &proofData)
	if err != nil {
		return err
	}

	// get OAuth publicKey
	jwtHeaderBytes, err := base64.RawURLEncoding.DecodeString(msg.ZkAuthSignature.ZkAuthInputs.HeaderBase64)
	if err != nil {
		return err
	}
	var jwtHeader JWTHeader
	if err = json.Unmarshal(jwtHeaderBytes, &jwtHeader); err != nil {
		return err
	}
	jwk := jks.GetJWK(jwtHeader.Kid)
	if jwk == nil {
		return errors.Errorf("no jwk of kid:%s", jwtHeader.Kid)
	}
	modulus, err := jwk.NBytes()
	if err != nil {
		return err
	}

	// calculate all input hash
	allInputHash, err := msg.ZkAuthSignature.ZkAuthInputs.CalculateAllInputsHash(ephPubKey, modulus, msg.ZkAuthSignature.MaxBlockHeight)
	if err != nil {
		return err
	}

	// verify
	proof := snarktypes.ZKProof{
		Proof:      &proofData,
		PubSignals: []string{allInputHash.String()},
	}

	return zkv.Verify(proof)
}
