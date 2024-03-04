package main

import (
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

func handleGetBlockchain() {

}

func handleWriteBlock() {

}

func respondWithJSON() {

}

func isBlockValid() bool {

}

func calculateHash() string {

}

func generateBlock() Block {

}

func isHashValid() bool {

}
