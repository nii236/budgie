package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token    string
	dbAdd    *sql.Stmt
	dbCheck  *sql.Stmt
	dbTopTen *sql.Stmt
	dev      bool
)

// Record is a single recorded transaction
type Record struct {
	Desc      string
	Amount    int
	CreatedAt time.Time
}

func init() {
	flag.BoolVar(&dev, "d", false, "Dev mode")
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func db() {

	database, err := sql.Open("sqlite3", "./budgie.db")
	if err != nil {
		panic(err)
	}

	if dev {
		database.Exec("DROP TABLE records")
	}
	stmt, err := database.Prepare(`
CREATE TABLE IF NOT EXISTS records
(
	id INTEGER PRIMARY KEY, 
	description TEXT, 
	cents INTEGER, 
	user TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`)

	stmt.Exec()
	dbAdd, err = database.Prepare("INSERT INTO records (description, cents, user) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}

	dbCheck, err = database.Prepare("SELECT IFNULL(SUM(cents), 0) FROM records")
	if err != nil {
		panic(err)
	}

	dbTopTen, err = database.Prepare("SELECT description, cents, created_at FROM records ORDER BY created_at DESC LIMIT 10;")
	if err != nil {
		panic(err)
	}
}

func main() {
	db()
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(help)
	dg.AddHandler(ping)
	dg.AddHandler(add)
	dg.AddHandler(sub)
	dg.AddHandler(check)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!help") {
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, `ping - test connection
add - add money
sub - subtract money
help - this text
`)
	if err != nil {
		fmt.Println(err)
	}

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!ping") {
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
	if err != nil {
		fmt.Println(err)
	}
}

func add(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!add") {
		return
	}

	fields := strings.Fields(m.Content)
	if len(fields) != 3 {
		s.ChannelMessageSend(m.ChannelID, "Bad args: !add <description> <amount>")
		return
	}
	amount, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		fmt.Println(err)
	}
	_, err = dbAdd.Exec(fields[1], int(amount*100), m.Author.ID)
	if err != nil {
		fmt.Println(err)
	}
	var total int
	err = dbCheck.QueryRow().Scan(&total)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you %s, your total balance is now $%.2f", m.Author.Username, compliment(), float64(total)/100))
	if err != nil {
		fmt.Println(err)
	}
}
func sub(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!sub") {
		return
	}
	fields := strings.Fields(m.Content)

	if len(fields) != 3 {
		s.ChannelMessageSend(m.ChannelID, "Bad args: !sub <description> <amount>")
		return
	}
	amount, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		fmt.Println(err)
	}
	_, err = dbAdd.Exec(fields[1], int(amount*-100), m.Author.ID)
	if err != nil {
		fmt.Println(err)
	}

	var total int
	err = dbCheck.QueryRow().Scan(&total)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you %s, your total balance is now $%.2f", m.Author.Username, compliment(), float64(total)/100))
	if err != nil {
		fmt.Println(err)
	}
}

func check(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!check") {
		return
	}
	fields := strings.Fields(m.Content)

	if len(fields) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Bad args: !check")
		return
	}

	var total int
	err := dbCheck.QueryRow().Scan(&total)
	if err != nil {
		fmt.Println(err)
		return
	}

	records := []*Record{}
	rows, err := dbTopTen.Query()
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		rec := &Record{}
		err = rows.Scan(&rec.Desc, &rec.Amount, &rec.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return
		}
		records = append(records, rec)
	}

	list := ""
	for _, r := range records {
		list += fmt.Sprintf("%s %s %.2f\n", r.CreatedAt.Format("2006-01-02"), r.Desc, float64(r.Amount)/100)
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Total Balance: $%.2f", float64(total)/100))
	s.ChannelMessageSend(m.ChannelID, "Last 10 transactions:")
	s.ChannelMessageSend(m.ChannelID, list)
}
