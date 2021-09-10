package keys

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"

	cryptokeyring "github.com/line/lbm-sdk/crypto/keyring"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"
)

type bechKeyOutFn func(keyInfo cryptokeyring.Info) (cryptokeyring.KeyOutput, error)

func printKeyInfo(w io.Writer, keyInfo cryptokeyring.Info, bechKeyOut bechKeyOutFn, output string) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}

	switch output {
	case OutputFormatText:
		printTextInfos(w, []cryptokeyring.KeyOutput{ko})

	case OutputFormatJSON:
		out, err := KeysCdc.MarshalJSON(ko)
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, string(out))
	}
}

func printInfos(w io.Writer, infos []cryptokeyring.Info, output string) {
	kos, err := cryptokeyring.Bech32KeysOutput(infos)
	if err != nil {
		panic(err)
	}

	switch output {
	case OutputFormatText:
		printTextInfos(w, kos)

	case OutputFormatJSON:
		out, err := KeysCdc.MarshalJSON(kos)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "%s", out)
	}
}

func printTextInfos(w io.Writer, kos []cryptokeyring.KeyOutput) {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(out))
}

func printKeyAddress(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.Address)
}

func printPubKey(w io.Writer, info cryptokeyring.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, ko.PubKey)
}
