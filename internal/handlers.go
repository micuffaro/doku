package doku

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "math/rand"
    "net/http"
    "time"
)

// TODO: Figure out better storage, like mongodb
// sudokuStorage stores Sudoku objects
var sudokuStorage = make(map[string]*Sudoku)

func dumpJSONToSession(b io.ReadCloser) Session {

    body, _ := ioutil.ReadAll(b)

    var session Session
    err := json.Unmarshal(body, &session)
    if err != nil {
        fmt.Println(err)
    }

    return session
}

// Session is used to map the json data for the session.
type Session struct {
    Grid   string `json:"grid"`
    Hash   string `json:"hash"`
    Solved bool   `json:"solved"`
    Error  string `json:"error"`
}

// NewSudokuHandler initializes a sudokuHandler
func NewSudokuHandler() http.Handler { return sudokuHandler{} }

type sudokuHandler struct{}

func (h sudokuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    session := dumpJSONToSession(r.Body)

    rand.Seed(time.Now().UnixNano())
    session.Hash = randomString(5)

    s := NewSudoku(session.Grid)
    sudokuStorage[session.Hash] = s

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}

// NewDisplayHandler initializes a displayHandler
func NewDisplayHandler() http.Handler { return displayHandler{} }

type displayHandler struct{}

func (h displayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    session := dumpJSONToSession(r.Body)

    if _, ok := sudokuStorage[session.Hash]; !ok {
        session.Error = "Sudoku not found"
    } else {
        s := sudokuStorage[session.Hash]
        // TODO: This won't work, grid is a map :)
        session.Grid = s.DisplayString()
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}

// NewSolveHandler initializes a solveHandler
func NewSolveHandler() http.Handler { return solveHandler{} }

type solveHandler struct{}

func (h solveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    session := dumpJSONToSession(r.Body)

    if _, ok := sudokuStorage[session.Hash]; !ok {
        session.Error = "Sudoku not found"
    } else {
        s := sudokuStorage[session.Hash]
        s.Solve()
    }
}

// NewStateHandler initializes a stateHandler
func NewStateHandler() http.Handler { return stateHandler{} }

type stateHandler struct{}

func (h stateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    session := dumpJSONToSession(r.Body)

    if _, ok := sudokuStorage[session.Hash]; !ok {
        session.Error = "Sudoku not found"
    } else {
        s := sudokuStorage[session.Hash]
        session.Solved = s.isSolved()
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(session)
}
