package types

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	snarktypes "github.com/iden3/go-rapidsnark/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
)

func TestZKAuthVerifierVerify(t *testing.T) {
	const proofStr = "{\n  \"pi_a\": [\n    \"19522371839270073620652913303981953856602066371307169133095252437574898167087\",\n    \"4833523328406165158329834156323299381989099864962000556295485103939841522677\",\n    \"1\"\n  ],\n  \"pi_b\": [\n    [\n      \"13590751830803055634813187832453152641731583060610323629642183986767393594680\",\n      \"14562876116960780180699310377103898768025569460014751784685583318505733065584\"\n    ],\n    [\n      \"9109019813807833053703748676868056381436829864508782138974623759022676848785\",\n      \"15372401783891458941369342924584882274874259166617081630036949233929956089219\"\n    ],\n    [\n      \"1\",\n      \"0\"\n    ]\n  ],\n  \"pi_c\": [\n    \"19650458257976957471150466891156721008030985770755836553951832770722262072792\",\n    \"18859262322622715946926814109590642931253610417578764508946716428980514815460\",\n    \"1\"\n  ],\n  \"protocol\": \"groth16\",\n  \"curve\": \"bn128\"\n}"
	const publicDataStr = "2897363891707776560374456764255972429123418378332336636175790055619926988108"

	verificationKey, err := os.ReadFile("../testutil/testdata/verification_key.json")
	require.NoError(t, err)

	zkAuthVerifier := NewZKAuthVerifier(verificationKey)

	var proofData snarktypes.ProofData
	err = json.Unmarshal([]byte(proofStr), &proofData)
	require.NoError(t, err)

	proof := snarktypes.ZKProof{
		Proof:      &proofData,
		PubSignals: []string{publicDataStr},
	}

	err = zkAuthVerifier.Verify(proof)
	require.NoError(t, err)
}

type zkKeeperMock struct {
	jwks     *JWKsMap
	verifier *ZKAuthVerifier
}

var _ ZKAuthKeeper = &zkKeeperMock{}

func (z *zkKeeperMock) GetJWK(kid string) *JWK {
	return z.jwks.GetJWK(kid)
}

func (z *zkKeeperMock) GetVerifier() *ZKAuthVerifier {
	return z.verifier
}

func TestVerifyZKAuthSignature(t *testing.T) {
	verificationKey, err := os.ReadFile("../testutil/testdata/verification_key.json")
	require.NoError(t, err)

	zkAuthVerifier := NewZKAuthVerifier(verificationKey)

	ephPubKey, ok := new(big.Int).SetString("18948426102457371978524559226152399917062673825697601263047735920285791872240", 10)
	require.True(t, ok)

	const proofStr = "{\n \"pi_a\": [\n  \"7575287679446209007446416020137456670042570578978230730578011103770415897062\",\n  \"20469978368515629364541212704109752583692706286549284712208570249653184893207\",\n  \"1\"\n ],\n \"pi_b\": [\n  [\n   \"4001119070037193619600086014535210556571209449080681376392853276923728808564\",\n   \"18475391841797083641468254159150812922259839776046448499150732610021959794558\"\n  ],\n  [\n   \"19781252109528278034156073207688818205850783935629584279449144780221040670063\",\n   \"5873714313814830719712095806732872482213125567325442209795797618441438990229\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"pi_c\": [\n  \"18920522434978516095250248740518039198650690968720755259416280639852277665022\",\n  \"1945774583580804632084048753815901730674007769630810705050114062476636502591\",\n  \"1\"\n ],\n \"protocol\": \"groth16\",\n \"curve\": \"bn128\"\n}"
	zkAuthMsg := MsgExecution{
		Msgs: nil,
		ZkAuthSignature: ZKAuthSignature{
			ZkAuthInputs: &ZKAuthInputs{
				ProofPoints:  []byte(proofStr),
				IssBase64:    "aHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29t",
				HeaderBase64: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU1YzE4OGE4MzU0NmZjMTg4ZTUxNTc2YmE3MjgzNmUwNjAwZThiNzMiLCJ0eXAiOiJKV1QifQ",
				AddressSeed:  "15035161560159971633800983619931498696152633426768016966057770643262022096073",
			},
			MaxBlockHeight: 32754,
		},
	}

	jks := NewJWKs()
	jks.AddJWK(&JWK{
		Kty: "RSA",
		E:   "AQAB",
		N:   "q0CrF3x3aYsjr0YOLMOAhEGMvyFp6o4RqyEdUrnTDYkhZbcud-fJEQafCTnjS9QHN1IjpuK6gpx5i3-Z63vRjs5EQX7lP1jG8Qg-CnBdTTLw4uJi7RmmlKPsYaO1DbNkFO2uEN62sOOzmJCh1od3CZXI1UYH5cvZ_sLJaN2A4TwvUTU3aXlXbUNJz_Hy3l0q1Jjta75NrJtJ7Pfj9tVXs8qXp15tZXrnbaM-AI0puswt35VsQbmLwUovFFGeToo5q2c_c1xYnV5uQYMadANekGPRFPM9JZpSSIvH0Lv_f15V2zRqmIgX7a3RcmTnr3-w3QNQTogdy-MogxPUdRbxow",
		Alg: "RS256",
		Kid: "55c188a83546fc188e51576ba72836e0600e8b73",
	})

	require.NoError(t, err)

	var zkKeeper = zkKeeperMock{
		jwks:     jks,
		verifier: &zkAuthVerifier,
	}

	ctx := sdk.Context{}
	err = VerifyZKAuthSignature(ctx, &zkKeeper, ephPubKey.Bytes(), &zkAuthMsg)
	require.NoError(t, err)
}
