package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccAddressFromAddressSeed(t *testing.T) {
	addressSeeds := map[string]string{
		"6837746941479125904804058040336379861713403246941015384197322102329773944159":  "link15uxm4w27fyuy0dcczvnu7edq7dnuf3r6t97t82j3fuc2pcfux4zsrxq8mu",
		"15631290963325153295531082175644052605061447806381313417788319283841674180543": "link1etuxuak3q5tgjs8ew09h5q4yct5vrcc7thl2llcjpk72fuq8r50q958h26",
		"11967859756179757815812269077480243698674971831305032302706373939074904776780": "link1muet5uz59utufwvux67tn5lgez28dasg4renvzhuku4a0jz459cs7m6tuw",
		"10248927216838193947032332743622432828052626024625324374470479888091711985426": "link14x0mfl95kxg739fqpqdyvrykhd9ghsqk4uew3tr4aqcprfmtpaesh858kj",
		"16423356555234621455999265971523270820664611063213539441365583262907894845196": "link1a4m70fx99udpr6c2excsgxa9x0vmra9gvzs4w6cz78dx3djrhuaqlppfp2",
	}

	iss := "accounts.google.com"
	for seed, exp := range addressSeeds {
		acc, err := AccAddressFromAddressSeed(seed, iss)
		require.NoError(t, err)
		require.Equal(t, exp, acc.String())
	}
}
