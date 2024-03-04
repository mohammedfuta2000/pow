package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const difficulty = 1

type Block struct {
	Index      int
	TimeStamp  string
	Data       int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

var Blockchain []Block

type Message struct{
	Data int

}
var mutex = &sync.Mutex{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	go func() {
		t := time.Now()
		genesisBlock:=Block{}
		genesisBlock = Block{
			Index:      0,
			TimeStamp:  t.String(),
			Data:       0,
			Hash:       calculateHash(genesisBlock),
			PrevHash:   "",
			Difficulty: difficulty,
			Nonce:      ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())

}

func run() error {
	mux:=makeMuxRouter()
	httpPort:= os.Getenv("PORT")
	log.Println("the http server is running and listening on port",httpPort)
	s := &http.Server{
		Addr: ":"+httpPort,
		Handler: mux,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
		MaxHeaderBytes: 1<<20,
	}

	if err:= s.ListenAndServe(); err!=nil{
		return err
	}
	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter:= mux.NewRouter()
	
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter,r *http.Request) {
	// bytes, err:= json.Marshal(Blockchain)
	bytes, err:= json.MarshalIndent(Blockchain,""," ")
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	var m Message

	if err:= json.NewDecoder(r.Body).Decode(&m);err!=nil{
		respondWithJSON(w,r,http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()


}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type","application/json")
	response, err:= json.MarshalIndent(payload,""," ")
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal server error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func isBlockValid() bool {

}

func calculateHash() string {

}

func generateBlock() Block {

}

func isHashValid() bool {

}
