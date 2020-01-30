package verifier

import (
	"github.com/line/link/cmd/contract_tests/unmarshaler"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestElementType(t *testing.T) {
	var elementTests = []struct {
		json1    string
		json2    string
		expected bool
	}{
		{
			"30",
			"50",
			true,
		},
		{
			`"30"`,
			`"50"`,
			true,
		},
		{
			`"null"`,
			`"null"`,
			true,
		},
		{
			`"null"`,
			`"50"`,
			true,
		},
		{
			`"null"`,
			`30`,
			true,
		},
	}

	for _, tt := range elementTests {
		t.Logf("compare %s and %s", tt.json1, tt.json2)
		{
			data1 := unmarshaler.UnmarshalJSON(&tt.json1)
			data2 := unmarshaler.UnmarshalJSON(&tt.json2)
			require.Equal(t, tt.expected, CompareJSONFormat(data1.Body, data2.Body))
		}
	}
}

func TestArray(t *testing.T) {
	var arrayTests = []struct {
		json1    string
		json2    string
		expected bool
	}{
		{
			`[]`,
			`[]`,
			true,
		},
		{
			`"30"`,
			`"50"`,
			true,
		},
		{
			`"abc"`,
			`[]`,
			false,
		},
		{
			`""`,
			`[]`,
			true,
		},
		{
			`[123]`,
			`[123123, "abc", "ad"]`,
			true,
		},
		{
			`["a", "b", 1, 43]`,
			`["a", "b", 1, 43]`,
			true,
		},
		{
			`[{"a":1}, {"b": 1}]`,
			`[{"a":2}, {"b": 2}]`,
			true,
		},
		{
			`[{"a":1}]`,
			`[{"a":2}, {"b": 2}]`,
			true,
		},
		{
			`[{"b":1}]`,
			`[{"a":2}, {"b": 2}]`,
			false,
		},
	}

	for _, tt := range arrayTests {
		t.Logf("compare %s and %s", tt.json1, tt.json2)
		{
			data1 := unmarshaler.UnmarshalJSON(&tt.json1)
			data2 := unmarshaler.UnmarshalJSON(&tt.json2)
			require.Equal(t, tt.expected, CompareJSONFormat(data1.Body, data2.Body))
		}
	}

}

func TestObject(t *testing.T) {
	var objectTests = []struct {
		json1    string
		json2    string
		expected bool
	}{
		{
			`{"key1":"value1"}`,
			`{"key1":"value1"}`,
			true,
		},
		{
			`{"key1":"value1"}`,
			`{"key2":"value1"}`,
			false,
		},
		{
			`{"key1":"value1"}`,
			`{"key1":"value2"}`,
			true,
		},
		{
			`{"key1":"value1", "key2":"value2"}`,
			`{"key1":"value2"}`,
			false,
		},
		{
			`{"key1":"value1", "key2":113}`,
			`{"key1":119, "key2":"value2"}`,
			true,
		},
		{
			`{"key1":"value1", "key2":113, "key3":113, "key4":113, "key5":113}`,
			`{"key1":119, "key2":"value2", "key3":"value2", "key4":"value2", "key5":"value2"}`,
			true,
		},
		{
			`{"key1":"value1", "key2":113, "key3":113, "key4":113, "key5":113}`,
			`{"key0":119, "key2":"value2", "key3":"value2", "key4":"value2", "key5":"value2"}`,
			false,
		},
	}

	for _, tt := range objectTests {
		t.Logf("compare %s and %s", tt.json1, tt.json2)
		{
			data1 := unmarshaler.UnmarshalJSON(&tt.json1)
			data2 := unmarshaler.UnmarshalJSON(&tt.json2)
			require.Equal(t, tt.expected, CompareJSONFormat(data1.Body, data2.Body))
		}
	}
}

func TestNestedObject(t *testing.T) {
	var nestedObjectTests = []struct {
		json1    string
		json2    string
		expected bool
	}{
		{
			`{"key1":"value1", "sub1": {"key2":113, "sub2":{"key3":113}}, "sub3":{"key4":113, "key5":113}}`,
			`{"key1":119, "sub1":{"key2":"value2", "sub2":{"key3":"value2"}}, "sub3":{"key4":"value2", "key5":"value2"}}`,
			true,
		},
		{
			`{"key1":"value1", "sub1": {"key2":113, "sub2":{"key":113}}, "sub3":{"key4":113, "key5":113}}`,
			`{"key1":119, "sub1":{"key2":"value2", "sub2":{"key3":"value2"}}, "sub3":{"key4":"value2", "key5":"value2"}}`,
			false,
		},
		{
			`{"key1":"value1", "sub1": [{"key2":113, "sub2":{"key":[113]}}], "sub3":{"key4":113, "key5":113}}`,
			`{"key1":119, "sub1":[{"key2":"value2", "sub2":{"key3":["value2"]}}], "sub3":{"key4":"value2", "key5":"value2"}}`,
			false,
		},
		{
			`{"key1":"value1", "sub1": [1, 2], "sub3":{"key4":113, "key5":113}}`,
			`{"key1":119, "sub1":[{"key2":"value2", "sub2":{"key3":["value2"]}}], "sub3":{"key4":"value2", "key5":"value2"}}`,
			false,
		},
		{
			`{"key1":"value1", "sub1": [{"key0":113, "sub4":{"key5":[113]}}], "sub3":{"key4":113, "key5":113}}`,
			`{"key1":119, "sub1":[{"key2":"value2", "sub2":{"key3":["value2"]}}], "sub3":{"key4":"value2", "key5":"value2"}}`,
			false,
		},
	}

	for _, tt := range nestedObjectTests {
		t.Logf("compare %s and %s", tt.json1, tt.json2)
		{
			data1 := unmarshaler.UnmarshalJSON(&tt.json1)
			data2 := unmarshaler.UnmarshalJSON(&tt.json2)
			require.Equal(t, tt.expected, CompareJSONFormat(data1.Body, data2.Body))
		}
	}
}

func TestTxResult(t *testing.T) {
	var txResultsTests = []struct {
		json1    string
		json2    string
		expected bool
	}{
		{
			exampleJSON,
			exampleJSONWithDiffrentValue,
			true,
		},
		{
			exampleJSON,
			exampleJSONWithoutTimeStamp,
			false,
		},
	}

	for _, tt := range txResultsTests {
		t.Logf("compare %s and %s", tt.json1, tt.json2)
		{
			data1 := unmarshaler.UnmarshalJSON(&tt.json1)
			data2 := unmarshaler.UnmarshalJSON(&tt.json2)
			require.Equal(t, tt.expected, CompareJSONFormat(data1.Body, data2.Body))
		}
	}
}

const (
	exampleJSON = `{
	  "total_count": "1",
	  "count": "1",
	  "page_number": "1",
	  "page_total": "1",
	  "limit": "30",
	  "txs": [
		{
		  "height": "368",
		  "txhash": "777467B3B9FF7918B15D374EC0510AADE8026D3B91795AFA25B50888199DFD20",
		  "index": 1,
		  "code": 0,
		  "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
		  "logs": [
			{
			  "msg_index": 0,
			  "success": false,
			  "log": ""
			}
		  ],
		  "gas_wanted": "200000",
		  "gas_used": "26354",
		  "events": [
			{
			  "type": "",
			  "attributes": [
				{
				  "key": "",
				  "value": ""
				}
			  ]
			}
		  ],
		  "tx": {
			"type": "cosmos-sdk/StdTx",
			"value": {
			  "msg": [
				{
				  "type": "cosmos-sdk/MsgSend",
				  "value": {
					"from_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"to_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"amount": [
					  {
						"denom": "link",
						"amount": "1"
					  }
					]
				  }
				}
			  ],
			  "fee": {
				"amount": [
				  {
					"denom": "link",
					"amount": "50"
				  }
				],
				"gas": "30000"
			  },
			  "signatures": [
				{
				  "pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "AmtHHdw+9xASQejvAzOwj/BT8AumZGJHWKKb4Gg5hssC"
				  },
				  "signature": "7uTC74QlknqYWEwg7Vn6M8Om7FuZ0EO4bjvuj6rwH1mTUJrRuMMZvAAqT9VjNgP0RA/TDp6u/92AqrZfXJSpBQ=="
				}
			  ],
			  "memo": ""
			}
		  },
          "timestamp": "2020-01-03T12:26:11Z"
		}
	  ]
	}
	`
	exampleJSONWithDiffrentValue = `
	{
	  "total_count": "4",
	  "count": "4",
	  "page_number": "1",
	  "page_total": "1",
	  "limit": "50",
	  "txs": [
		{
		  "height": "2",
		  "txhash": "9DCF862A32F2718A915D8277F827CA74DD967379C6349677924EBABE451853B5",
		  "index": 0,
		  "code": 0,
		  "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
		  "logs": [
			{
			  "msg_index": 0,
			  "success": true,
			  "log": ""
			}
		  ],
		  "gas_wanted": "200000",
		  "gas_used": "23291",
		  "events": [
			{
			  "type": "message",
			  "attributes": [
				{
				  "key": "action",
				  "value": "send"
				},
				{
				  "key": "sender",
				  "value": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80"
				},
				{
				  "key": "module",
				  "value": "bank"
				}
			  ]
			},
			{
			  "type": "transfer",
			  "attributes": [
				{
				  "key": "recipient",
				  "value": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80"
				},
				{
				  "key": "amount",
				  "value": "1link"
				}
			  ]
			}
		  ],
		  "tx": {
			"type": "cosmos-sdk/StdTx",
			"value": {
			  "msg": [
				{
				  "type": "cosmos-sdk/MsgSend",
				  "value": {
					"from_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"to_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"amount": [
					  {
						"denom": "link",
						"amount": "1"
					  }
					]
				  }
				}
			  ],
			  "fee": {
				"amount": [],
				"gas": "200000"
			  },
			  "signatures": [
				{
				  "pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "A1Z2vdw7e1SEGxzvLLyADzI1JKvyZD3V81/yUdTvBnUL"
				  },
				  "signature": "Z3gKHw9qEu8JCz1R9L8D2e8AAZJVU1doksSUDmRfqJRq4iQ+PMNC/TAJa5C0kT58Kilot2BZoXs+mMbzFqksWw=="
				}
			  ],
			  "memo": ""
			}
		  },
		  "timestamp": "2020-01-03T12:26:11Z"
		}
	  ]
	}
	`
	exampleJSONWithoutTimeStamp = `
	{
      "total_count": "1",
	  "count": "1",
	  "page_number": "1",
	  "page_total": "1",
	  "limit": "30",
	  "txs": [
		{
		  "height": "368",
		  "txhash": "777467B3B9FF7918B15D374EC0510AADE8026D3B91795AFA25B50888199DFD20",
		  "index": 1,
		  "code": 0,
		  "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
		  "logs": [
			{
			  "msg_index": 0,
			  "success": false,
			  "log": ""
			}
		  ],
		  "gas_wanted": "200000",
		  "gas_used": "26354",
		  "events": [
			{
			  "type": "",
			  "attributes": [
				{
				  "key": "",
				  "value": ""
				}
			  ]
			}
		  ],
		  "tx": {
			"type": "cosmos-sdk/StdTx",
			"value": {
			  "msg": [
				{
				  "type": "cosmos-sdk/MsgSend",
				  "value": {
					"from_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"to_address": "link16k7fxx6thgpltttzqczr5r40fx7c8sqt3tpc80",
					"amount": [
					  {
						"denom": "link",
						"amount": "1"
					  }
					]
				  }
				}
			  ],
			  "fee": {
				"amount": [
				  {
					"denom": "link",
					"amount": "50"
				  }
				],
				"gas": "30000"
			  },
			  "signatures": [
				{
				  "pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "AmtHHdw+9xASQejvAzOwj/BT8AumZGJHWKKb4Gg5hssC"
				  },
				  "signature": "7uTC74QlknqYWEwg7Vn6M8Om7FuZ0EO4bjvuj6rwH1mTUJrRuMMZvAAqT9VjNgP0RA/TDp6u/92AqrZfXJSpBQ=="
				}
			  ],
			  "memo": ""
			}
		  }
		}
	  ]
	}
	`
)
