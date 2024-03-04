package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	Data int `json:"data"`

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
	fmt.Println("xxxxxxxxxxxxx",m)
	defer r.Body.Close()

	// why is this mutex necessary
	mutex.Lock()
	newBlock:= generateBlock(Blockchain[len(Blockchain)-1],m.Data)
	mutex.Unlock()

	if isBlockValid(newBlock,Blockchain[len(Blockchain)-1]){
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
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

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1!=newBlock.Index{
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash{
		return false
	}
	if calculateHash(newBlock)!=newBlock.Hash{
		return false
	}
	return true
}

func calculateHash(block Block) string {
	record:= strconv.Itoa(block.Index)+block.TimeStamp+ strconv.Itoa(block.Data)+block.PrevHash+block.Nonce
	h := sha256.New()

	h.Write([]byte(record))
	hashed:= h.Sum(nil)
	// why this ↓
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, Data int) Block {
	var newBlock Block
	t:= time.Now()

	newBlock.Index= oldBlock.Index+1
	newBlock.TimeStamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("/%x",i)
		newBlock.Nonce = hex

		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty){
			fmt.Println(calculateHash(newBlock), "do more work")
			time.Sleep(time.Second)
			continue
		}else{
			fmt.Println(calculateHash(newBlock),"Work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}

func isHashValid(hash string, difficlty int) bool {
	// what is this
	prefix:= strings.Repeat("0",difficlty)
	// 0,2= 00
	// 0,1= 0
	return strings.HasPrefix(hash, prefix)
}
