package main

import (
    "os"
    "bufio"
    "fmt"
    "strconv"
    "rand"
    "time"
    "flag"
)

const emptyCell = 99
const boardSize = 3
const maxMoves  = boardSize * boardSize
const maxDepth  = 5

type Player byte
type TicTacToe struct {
    board [maxMoves]Player
    players [2]Player
    current Player
}

func (game *TicTacToe) PlayAt (row, col int, value Player) {
    game.board[row*boardSize+col] = value
}

func (game *TicTacToe) PlayedAt (row, col int) Player {
    return game.board[row*boardSize+col]
}

func (game *TicTacToe) Cell (num int) (row int, col int) {
    return num / boardSize, num % boardSize
}

func (game *TicTacToe) NextPlayer() Player {
    return ^game.current & 1
}

func (game *TicTacToe) Winner () (player Player) {
    for dim := 0; dim < 4; dim++ {
        for a := 0; a < boardSize; a++ {
            scores := [2]uint{0,0}
            for b := 0; b < boardSize; b++ {
                switch dim {
                    case 0: player = game.PlayedAt(a, b)
                    case 1: player = game.PlayedAt(b, a)
                    case 2: player = game.PlayedAt(b, b)
                    case 3: player = game.PlayedAt(b, boardSize-b-1)
                }
                if player == emptyCell {break}
                if scores[player]++; scores[player] == boardSize {
                    return player
                }
            }
            if dim > 1 {break}
        }
    }
    return emptyCell
}

func (game *TicTacToe) Copy () *TicTacToe {
    return &TicTacToe{game.board, game.players, game.current}
}

func (game *TicTacToe) Remaining () (moves []int) {
    moves = make([]int, 0, maxMoves)
    for i := 0; i < boardSize * boardSize; i++ {
        if game.board[i] == emptyCell {
            moves = append(moves, i)
        }
    }
    return moves
}

func (game *TicTacToe) Expand () []*TicTacToe {
    moves := game.Remaining()
    games := make([]*TicTacToe, len(moves))
    for i, move := range moves {
        row, col := game.Cell(move)
        games[i] = game.Copy()
        games[i].current = game.NextPlayer()
        games[i].PlayAt(row, col, games[i].current)
    }
    return games
}

type Outcome struct {
    move *TicTacToe
    score int
}

func (game *TicTacToe) Evaluate (response chan Outcome, depth int) Outcome {
    result := Outcome{game, 0}
    if winner := game.Winner(); winner != emptyCell {
        result.score = (maxMoves + 1) - depth
        // winner == game.current seems paradoxical, but at this point in the evaluation
        // game.current is the *other* player.
        if (winner == game.current) { result.score = -result.score }
        if response != nil {response <- result}
        return result
    }
    if depth == maxDepth {
        if response != nil {response <- result}
        return result
    }
    children := make(chan Outcome, maxMoves)
    dispatched := 0
    for _, move := range game.Expand() {
        go move.Evaluate(children, depth+1)
        dispatched++
    }
    best := make([]Outcome, 0, dispatched)
    for i := 0; i < dispatched; i++ {
        move := <- children
        // a concurrently-executed minimax search
        move.score = -move.score
        if len(best) == 0 || move.score == best[0].score {
            best = append(best, move)
        } else if len(best) > 0 && move.score > best[0].score {
            best = []Outcome{move}
        }
    }
    if len(best) > 0 { result.score = best[0].score }
    if response != nil {response <- result}
    if len(best) > 0 {
        return best[rand.Intn(len(best))]
    }
    return result
}

func (game *TicTacToe) BestMove () *TicTacToe {
    best := game.Evaluate(nil, 0)
    return best.move
}

func (game *TicTacToe) Display () {
    fmt.Print("\n")
    for row := 0; row < boardSize; row++ {
        fmt.Print(" ")
        for col := 0; col < boardSize; col++ {
            played := game.PlayedAt(row, col)            
            if played != emptyCell {
                fmt.Print(string(game.players[played]))
            } else {
                fmt.Print(row*boardSize+col)
            }
            if col < boardSize - 1 {
                fmt.Print(" | ")
            }
        }
        fmt.Print("\n")
        if row < boardSize - 1 {
            for col := 0; col < boardSize - 1; col++ {
                fmt.Print("---+")
            }
            fmt.Print("---\n")
        }
    }
    fmt.Print("\n")
}

func NewGame (players [2]Player) *TicTacToe {
    game := new(TicTacToe)
    game.players = players
    game.current = 1 // that way .NextPlayer() for the opening move is 0
    for row := 0; row < boardSize; row++ {
        for col := 0; col < boardSize; col++ {
            game.PlayAt(row, col, emptyCell)
        }
    }
    return game
}

func GameOver(game *TicTacToe) bool {
    if winner := game.Winner(); winner != emptyCell {
        fmt.Println(string(game.players[winner]), "wins!")
        return true
    }
    remains := game.Remaining()
    if len(remains) == 0 || (len(remains) == 1 && game.Expand()[0].Winner() == emptyCell) {
        fmt.Println("Stalemate!")
        return true
    }
    return false
}

func main() {
    player1 := flag.Bool("x", false, "Player is X; player goes first")
    player2 := flag.Bool("o", false, "Player is O; computer goes first")
    flag.Parse()

    rand.Seed(time.Nanoseconds())

    game := NewGame([2]Player{'X','O'}) 

    // if no choice has been made, pick randomly
    if !*player1 && (*player2 || rand.Intn(2) == 0) {
        fmt.Println("You are playing O. Enter a # and hit return to move.")
        game = game.BestMove()
    } else {
        fmt.Println("You are playing X. Enter a # and hit return to move.")
    }
    game.Display()

    in := bufio.NewReader(os.Stdin)

    for line, _, err := in.ReadLine();
        err == nil;
        line, _, err = in.ReadLine() {

        val, conv_err := strconv.Atoi(string(line))
        if conv_err != nil { continue }

        var row, col int = game.Cell(val)
        if game.PlayedAt(row, col) != emptyCell {continue}

        game.current = game.NextPlayer()
        game.PlayAt(row, col, game.current)
        game.Display()
        if GameOver(game) { break }

        game = game.BestMove()
        game.Display()
        if GameOver(game) { break }
    }
}
