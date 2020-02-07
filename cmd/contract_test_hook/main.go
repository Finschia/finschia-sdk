package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/line/link/contract_test/unmarshaler"
	"github.com/line/link/contract_test/verifier"
	"github.com/snikch/goodman/hooks"
	"github.com/snikch/goodman/transaction"
)

const swaggerYAMLPath = "/tmp/contract_test/swagger.yaml"

func main() {
	// This must be compiled beforehand and given to dredd as parameter, in the meantime the server should be running
	h := hooks.NewHooks()
	server := hooks.NewServer(hooks.NewHooksRunner(h))

	// TODO: dredd only runs happy path test, so we have to skip when status code is not 200
	//  We have to add test of non-200 cases using another method
	h.BeforeEach(func(t *transaction.Transaction) {
		if t.Expected.StatusCode != "200" {
			t.Skip = true
		}
	})

	// It is difficult to reproduce the evidence, so skip to validate evidence
	h.Before("/blocks/latest > Get the latest block > 200 > application/json", func(t *transaction.Transaction) {
		makeExpectedEvidenceNull(t)
	})

	h.Before("/blocks/{height} > Get a block at a certain height > 200 > application/json", func(t *transaction.Transaction) {
		makeExpectedEvidenceNull(t)
	})

	// dredd can not validate items inside array in 12.1.0, so validate them in hook
	h.BeforeEachValidation(func(t *transaction.Transaction) {
		compareEachBody(t)
	})

	// Sometimes unconfirmed tx may not be there, so skip if not
	h.BeforeValidation("/unconfirmed_txs > Get the list of unconfirmed transactions > 200 > application/json", func(t *transaction.Transaction) {
		actual := unmarshaler.UnmarshalJSON(&t.Real.Body)
		if actual.GetProperty("txs") == nil {
			t.Skip = true
		}
	})

	// dredd can not parsing items in anyOf in 12.1.0, so add messages to expected body for verification
	h.BeforeValidation("/txs/{hash} > Get a Tx by hash > 200 > application/json", func(t *transaction.Transaction) {
		addMsgExamplesToExpected(t)
		compareEachBody(t)
	})

	server.Serve()
	defer server.Listener.Close()
	fmt.Print(h)
}

func makeExpectedEvidenceNull(t *transaction.Transaction) {
	expected := unmarshaler.UnmarshalJSON(&t.Expected.Body)
	expected.SetProperty([]string{"block", "evidence", "evidence"}, nil)
	newBody, err := json.Marshal(expected.Body)
	if err != nil {
		panic(fmt.Sprintf("fail to marshal expected body with %s", err))
	}
	t.Expected.Body = string(newBody)
}

func addMsgExamplesToExpected(t *transaction.Transaction) {
	bytes, err := ioutil.ReadFile(swaggerYAMLPath)
	if err != nil {
		panic(err)
	}
	yamlString := string(bytes)
	swaggerYAML := unmarshaler.UnmarshalYAML(&yamlString)
	value := swaggerYAML.GetProperty("components", "examples", "MsgExamples", "value")

	expected := unmarshaler.UnmarshalJSON(&t.Expected.Body)
	expected.SetProperty([]string{"tx", "value", "msg"}, value)
	newBody, err := json.Marshal(expected.Body)
	if err != nil {
		panic(fmt.Sprintf("fail to marshal expected body with %s", err))
	}
	t.Expected.Body = string(newBody)
}

func compareEachBody(t *transaction.Transaction) {
	expected := unmarshaler.UnmarshalJSON(&t.Expected.Body)
	actual := unmarshaler.UnmarshalJSON(&t.Real.Body)
	if !verifier.CompareJSONFormat(expected.Body, actual.Body) {
		t.Fail = true
	}
}
