package plugins

import (
    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    "fmt"
    "strings"
)

func init() {
    RegisterCommand("tictactoe", TicTacToeHandler)
    RegisterCommand("quiz", QuizHandler)
    RegisterCommand("wordle", WordleHandler)
}

type TicTacToeGame struct {
    Board    [3][3]string
    Players  map[string]string
    Current  string
}

var activeGames map[string]*TicTacToeGame

func TicTacToeHandler(client *whatsmeow.Client, evt *events.Message) {
    chatID := evt.Chat.JID.String()
    
    if activeGames == nil {
        activeGames = make(map[string]*TicTacToeGame)
    }
    
    if game, exists := activeGames[chatID]; exists {
        // Handle ongoing game moves
        pos := lib.GetMessageText(evt)
        makeMove(client, evt, game, pos)
    } else {
        // Start new game
        game = &TicTacToeGame{
            Board:   [3][3]string{{" "," "," "},{" "," "," "},{" "," "," "}},
            Players: make(map[string]string),
            Current: evt.Sender.String(),
        }
        activeGames[chatID] = game
        lib.Reply(client, evt, "TicTacToe game started! Use numbers 1-9 to make your move")
    }
}

func QuizHandler(client *whatsmeow.Client, evt *events.Message) {
    category := lib.GetMessageText(evt)
    
    question, options, answer := lib.GetRandomQuiz(category)
    if question == "" {
        lib.Reply(client, evt, "Failed to get quiz question")
        return
    }
    
    responseText := fmt.Sprintf("Question: %s\n\nOptions:\n%s", 
        question, strings.Join(options, "\n"))
    
    lib.Reply(client, evt, responseText)
    
    // Store answer for validation when user responds
    lib.StoreQuizAnswer(evt.Chat.JID.String(), answer)
}

func WordleHandler(client *whatsmeow.Client, evt *events.Message) {
    guess := lib.GetMessageText(evt)
    chatID := evt.Chat.JID.String()
    
    if !lib.HasActiveWordle(chatID) {
        // Start new game
        word := lib.GetRandomWord()
        lib.StartWordle(chatID, word)
        lib.Reply(client, evt, "Wordle game started! Guess the 5-letter word")
        return
    }
    
    if len(guess) != 5 {
        lib.Reply(client, evt, "Please enter a 5-letter word")
        return
    }
    
    result, solved := lib.CheckWordleGuess(chatID, guess)
    if solved {
        lib.Reply(client, evt, fmt.Sprintf("Congratulations! The word was: %s", guess))
        lib.EndWordle(chatID)
    } else {
        lib.Reply(client, evt, result)
    }
}
