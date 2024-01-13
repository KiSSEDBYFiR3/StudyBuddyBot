package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	MaxDeposit     = 5
	DefaultTimeout = 30
)

type AuctionPlayer struct {
	u       *storage.DBUser
	deposit int
}

// auctions мапа всех дейсвуюзих аукционов
var auctions = make(map[int][]*AuctionPlayer)

// startAuctionExec предоставляет метод Exec для начала аукциона в чате.
type startAuctionExec string

// Exec: /start_auction [timeout] - запускает аукцион в чате, в котором указана данная команда.
// timeout - необязательный параметр времени окончания аукциона.
func (a startAuctionExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if _, ok := auctions[chat.ID]; ok {
		return &Response{message: msgAuctionIsStarted, method: sendMessageMethod}, nil
	}

	auctions[chat.ID] = make([]*AuctionPlayer, 0)

	timeout := getAuctionTimeout(inMessage)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		msg, err := p.finishAuction(chat.ID)
		if err != nil {
			log.Println("[ERROR] in goroutine /start_auction", err)
		}
		if msg != "" {
			_ = p.tg.SendMessage(chat.ID, msg, "", -1)
		}
	}()

	return &Response{message: fmt.Sprintf(msgStartAuction, MaxDeposit), method: sendMessageMethod, parseMode: telegram.Markdown}, nil
}

func getAuctionTimeout(inMessage string) int {
	strs := strings.Fields(inMessage)
	if len(strs) < 2 {
		return DefaultTimeout
	}
	timeout, err := strconv.Atoi(strs[1])
	if err != nil {
		return DefaultTimeout
	}
	if timeout < 10 || timeout > 60 {
		return DefaultTimeout
	}
	return timeout
}

// addDeposit предоставляет метод Exec для внесения депозита в ауцион.
type addDepositExec string

// Exec: /deposit {amount} - вносит депозит в текущий аукцион. Amount - обязательный параметр.
func (a addDepositExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := p.addDeposit(inMessage, user, chat)
	if err != nil {
		return nil, e.Wrap("can't exec /deposit", err)
	}
	mthd := sendMessageMethod

	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// addDeposit возвращает сообщание для телеграм чата, после команды /deposit {amount}.
func (p *Processor) addDeposit(inMessage string, user *telegram.User, chat *telegram.Chat) (string, error) {
	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	if _, ok := auctions[chat.ID]; !ok {
		return msgAuctionNotStarted, nil
	}

	strs := strings.Fields(inMessage)
	if len(strs) < 2 {
		return msgErrorDepositCmd, nil
	}

	deposit, err := strconv.Atoi(strs[1])
	if err != nil {
		return msgErrorDepositCmd, nil
	}

	var needAdd bool
	player := getPlayer(dbUser)
	if player == nil {
		player = &AuctionPlayer{
			u:       dbUser,
			deposit: 0,
		}
		needAdd = true
	}

	if !p.canDeposit(deposit, dbUser, player) {
		return msgErrorDeposit, nil
	}

	if needAdd {
		auctions[chat.ID] = append(auctions[chat.ID], player)
	}
	err = p.changeDickSize(dbUser, -deposit)
	if err != nil {
		return "", err
	}
	player.deposit += deposit

	return fmt.Sprintf(msgSuccessDeposit, deposit), nil
}

// canDeposit проверяет может ли участник положить столько см пениса в аукцион.
func (p *Processor) canDeposit(deposit int, user *storage.DBUser, player *AuctionPlayer) bool {
	dickSize := user.DickSize
	playerDeposit := player.deposit
	return deposit >= 1 && deposit+playerDeposit <= MaxDeposit && dickSize-deposit >= 1
}

// getPlayer возвращает игрока аукциона.
func getPlayer(user *storage.DBUser) *AuctionPlayer {
	chatID := user.ChatID
	players := auctions[chatID]

	for _, p := range players {
		if user.ID == p.u.ID {
			return p
		}
	}

	return nil
}

// finishAuctionExec предоставляет метод Exec для завершения аукциона в чате.
// Только для админов бота.
type finishAuctionExec string

// Exec: /finish_auction - запускает аукцион в чате, в котором указана данная команда.
func (a finishAuctionExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/finish_auction)", errors.New("can't do this cmd"))
	}

	message, err := p.finishAuction(chat.ID)
	if err != nil {
		return nil, e.Wrap("can't finish duel", err)
	}

	return &Response{message: message, method: sendMessageMethod}, nil
}

// finishAuction случайным образом выбирает победителя из всех игроков аукциона.
func (p *Processor) finishAuction(chatID int) (string, error) {
	if _, ok := auctions[chatID]; !ok {
		return msgAuctionNotStarted, nil
	}

	players := auctions[chatID]
	if len(players) == 0 {
		delete(auctions, chatID)
		return msgNotEnoughPlayers, nil
	}

	winner, reward := getAuctionWinnerAndReward(auctions[chatID])
	delete(auctions, chatID)

	for i := 3; i > 0; i-- {
		p.tg.SendMessage(chatID, fmt.Sprintf("До результата аукциона: %d!", i), "", -1)
		time.Sleep(1 * time.Second)
	}

	err := p.changeDickSize(winner, reward)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(msgWinner, winner.Username, reward), nil
}

// getAuctionWinnerAndReward случайным образом определяет победителя аукциона.
// Возвращает пользователя победителя и общий призовой фонд.
func getAuctionWinnerAndReward(players []*AuctionPlayer) (*storage.DBUser, int) {
	var reward int

	ids := make([]int, 0)

	for _, player := range players {
		reward += player.deposit
		for i := 1; i <= player.deposit; i++ {
			ids = append(ids, player.u.ID)
		}
	}

	winnerID := ids[rand.Intn(len(ids))]

	var winner *storage.DBUser
	for _, player := range players {
		if player.u.ID == winnerID {
			winner = player.u
		}
	}

	return winner, reward
}

// auctionExec предоставляет метод Exec для просмотра аукциона в чате.
type auctionExec string

// Exec: /auction - возвращает список всех участников текущего аукциона.
func (a auctionExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	var message string
	mthd := sendMessageMethod
	parseMode := telegram.Markdown

	if _, ok := auctions[chat.ID]; !ok {
		return &Response{message: msgAuctionNotStarted, method: mthd, parseMode: parseMode}, nil
	}

	message = getAuctionPlayers(chat.ID)

	return &Response{message: message, method: mthd, parseMode: parseMode}, nil
}

// getAuctionPlayers возвращает список текущих участников аукциона.
func getAuctionPlayers(chatID int) string {
	players := auctions[chatID]

	if len(players) == 0 {
		return msgZeroPlayers
	}

	message := "Текущие игроки аукциона:\n\n"
	reward := 0

	for _, p := range players {
		if p.deposit > 0 {
			reward += p.deposit
			message += fmt.Sprintf("%s:\n*8", p.u.FirstName+" "+p.u.LastName)
			for i := 0; i < p.deposit/5; i++ {
				message += "="
			}
			message += "=Ð*\n"
		}
	}
	message += fmt.Sprintf("\nОбщий фонд *%d см*!", reward)
	return message
}
