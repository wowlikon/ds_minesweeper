package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	ds "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var dg *ds.Session
var count, side int

func main() {
	var err error

	//Loading .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	side, _ = strconv.Atoi(os.Getenv("SIDE"))
	count, _ = strconv.Atoi(os.Getenv("COUNT"))

	dg, err = ds.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %s", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)
	dg.Identify.Intents = ds.IntentsAll
	err = dg.Open()
	defer dg.Close()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *ds.Session, m *ds.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	parts := strings.Split(m.Content, " ")
	switch parts[0] {
	case "!play":
		mines := count
		if len(parts) == 2 {
			param, err := strconv.Atoi(parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid parameter")
			} else if param <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Too few mines")
			} else if param >= 81 {
				s.ChannelMessageSend(m.ChannelID, "Too many mines")
			} else {
				mines = param
			}
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Mines: %d", mines))
		mm, _ := s.ChannelMessageSend(m.ChannelID, minesweeper(mines))
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)

		if mines == 80 {
			for _, i := range []int{-1, 3, 5, 8} {
				s.MessageReactionAdd(m.ChannelID, mm.ID, ItoE(i))
			}
		}
		break
	case "!rand":
		max_val := count
		if len(parts) == 2 {
			param, err := strconv.Atoi(parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid parameter")
			} else if param <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Must be >0")
			} else {
				max_val = param
			}
		}
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		rm, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Random value: %d/%d", rand.Intn(max_val), max_val))
		s.MessageReactionAdd(m.ChannelID, rm.ID, "🎲")
		break
	}
}

func messageReactionAdd(s *ds.Session, m *ds.MessageReactionAdd) {
	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("Error getting message:", err)
		return
	}
	if (msg.Author.ID == s.State.User.ID) && (m.UserID != s.State.User.ID) {
		if m.Emoji.Name == "🎲" {
			if !strings.HasPrefix(msg.Content, "Random value: ") {
				return
			}

			max_val := count
			part := strings.Split(msg.Content, "/")[1]
			roll, err := strconv.Atoi(part)
			if err == nil {
				max_val = roll
			}

			s.ChannelMessageDelete(m.ChannelID, m.MessageID)
			rm, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Random value: %d/%d", rand.Intn(max_val), max_val))
			s.MessageReactionAdd(m.ChannelID, rm.ID, "🎲")
		}
	}
}

func minesweeper(mines int) string {
	// Initialize the field
	field := make([][]int, side)
	for i := 0; i < side; i++ {
		field[i] = make([]int, side)
	}

	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			field[i][j] = 0
		}
	}

	//Generate
	for mines > 0 {
		x := rand.Intn(side)
		y := rand.Intn(side)
		if field[x][y] != -1 {
			field[x][y] = -1
			mines--
		}
	}

	//Calculate
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			if field[i][j] != -1 {
				count := 0
				for k := i - 1; k <= i+1; k++ {
					for l := j - 1; l <= j+1; l++ {
						if k >= 0 && k < side && l >= 0 && l < side {
							if field[k][l] == -1 {
								count++
							}
						}
					}
				}
				field[i][j] = count
			}
		}
	}

	//Print
	var output string
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			output += "||" + ItoE(field[i][j]) + "||"
		}
		output += "\n"
	}
	return output
}

func ItoE(value int) string {
	switch value {
	case 0:
		return "0️⃣"
	case 1:
		return "1️⃣"
	case 2:
		return "2️⃣"
	case 3:
		return "3️⃣"
	case 4:
		return "4️⃣"
	case 5:
		return "5️⃣"
	case 6:
		return "6️⃣"
	case 7:
		return "7️⃣"
	case 8:
		return "8️⃣"
	case 9:
		return "9️⃣"
	default:
		return "💣"
	}
}
