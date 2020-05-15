package mock

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"sync"

	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
	"github.com/line/link/app"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)
var mutex sync.Mutex

type CallCounter struct {
	QueryAccountCallCount int
	QueryBlockCallCount   int
	BroadcastTxCallCount  int
	TargetLoadCallCount   int
	TargetFireCallCount   int
}

func (c *CallCounter) increment(count *int) {
	mutex.Lock()
	*count++
	mutex.Unlock()
}

func NewCallCounter() *CallCounter {
	return &CallCounter{}
}

func NewServer() *httptest.Server {
	callCounter := NewCallCounter()

	r := mux.NewRouter()
	r.HandleFunc("/call_counter", GetCallCounterHandler(callCounter)).Methods("GET")
	r.HandleFunc("/call_counter/clear", ClearCallCounterHandler(callCounter)).Methods("GET")
	r.HandleFunc("/target/load", TargetLoadHandler(callCounter)).Methods("POST")
	r.HandleFunc("/target/fire", TargetFireHandler(callCounter)).Methods("POST")
	r.HandleFunc("/auth/accounts/{address}", QueryAccountHandler(callCounter)).Methods("GET")
	r.HandleFunc("/blocks/{height}", QueryBlockHandler(callCounter)).Methods("GET")
	r.HandleFunc("/blocks/latest", QueryBlockHandler(callCounter)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxHandler(callCounter)).Methods("POST")
	return httptest.NewServer(r)
}

func GetCallCounterHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		bytes, err := json.Marshal(cc)
		if err != nil {
			log.Printf("Marshal CallCounter failed: %v", err)
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)
		logerr(w.Write(bytes))
	}
}

func ClearCallCounterHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nc := NewCallCounter()
		*cc = *nc
	}
}

func TargetLoadHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc.increment(&cc.TargetLoadCallCount)
		success(w)
	}
}

func TargetFireHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc.increment(&cc.TargetFireCallCount)
		w.WriteHeader(http.StatusOK)
		logerr(w.Write(loadResponse(TargetFireResponse)))
	}
}

func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	logerr(w.Write([]byte("success")))
}

func QueryAccountHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc.increment(&cc.QueryAccountCallCount)
		w.WriteHeader(http.StatusOK)
		logerr(w.Write(loadResponse(AccountResponse)))
	}
}

func QueryBlockHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc.increment(&cc.QueryBlockCallCount)
		w.WriteHeader(http.StatusOK)
		logerr(w.Write(loadResponse(BlockResponse)))
	}
}

func BroadcastTxHandler(cc *CallCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cc.increment(&cc.BroadcastTxCallCount)

		var req authrest.BroadcastReq
		status := http.StatusOK
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("BroadcastTx failed: %v", err)
			status = http.StatusInternalServerError
		}
		err = app.MakeCodec().UnmarshalJSON(body, &req)
		if err != nil {
			log.Printf("BroadcastTx failed: %v", err)
			status = http.StatusInternalServerError
		}

		var response string
		switch req.Mode {
		case "block":
			response = TxBlockResponse
		case "sync":
			response = TxSyncResponse
		case "async":
			response = TxAsyncResponse
		default:
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		logerr(w.Write(loadResponse(response)))
	}
}

func loadResponse(fileName string) []byte {
	data, err := ioutil.ReadFile(basepath + "/mock_response/" + fileName)
	if err != nil {
		panic(err)
	}
	return data
}

func logerr(n int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func GetCallCounter(url string) (res CallCounter) {
	resp, err := http.Get(url + "/call_counter")
	if err != nil {
		log.Printf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Request failed: %v", err)
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Printf("Request failed: %v", err)
	}

	return
}

func ClearCallCounter(url string) {
	resp, err := http.Get(url + "/call_counter/clear")
	if err != nil {
		log.Printf("Request failed: %v", err)
	}
	defer resp.Body.Close()
}
