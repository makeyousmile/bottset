package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sort"
	"strings"
)

func main() {
	ctx, client := ConnectToDb()
	db := client.Database("bot420")

	Collections, err := db.ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		log.Print(err)
	}

	cityNames, catNames := SplitCollName(Collections)
	log.Print(getMap(Collections)["Барановичи"])

	bot, err := tgbotapi.NewBotAPI("1284532231:AAFDnFJGFS7IEpcgelRfGIgbEWW6az8FRnA")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Print(".")
	//ChatIds := make(map[int]UserCfg)
	for update := range updates {
		if update.Message != nil {
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Command() {
			case "start":

				var numericKeyboard tgbotapi.InlineKeyboardMarkup
				for _, city := range cityNames {
					row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(city, city))
					numericKeyboard.InlineKeyboard = append(numericKeyboard.InlineKeyboard, row)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите город")
				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)

			}
		}
		if update.CallbackQuery != nil {

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			data := update.CallbackQuery.Data
			log.Print(data)

			var CatKeyboard tgbotapi.InlineKeyboardMarkup

			for _, cat := range catNames {
				row := tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(cat, cat),
				)
				CatKeyboard.InlineKeyboard = append(CatKeyboard.InlineKeyboard, row)
			}
			//edit top text
			editedMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Выберите категорию")
			bot.Send(editedMsg)
			//edit body
			editedMsg2 := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, CatKeyboard)
			bot.Send(editedMsg2)

		}

	}

}

func SplitCollName(collections []string) ([]string, []string) {
	var city []string
	var cat []string

	sort.Strings(collections)

	for _, col := range collections {
		names := strings.Split(col, ":")

		if len(names) == 2 {
			city = append(city, names[0])
			cat = append(cat, names[1])
		}

	}

	return city, cat
}
func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func getMap(collections []string) map[string][]string {
	coll := make(map[string][]string)

	cities, categories := SplitCollName(collections)

	for i := 0; i < len(cities); i++ {
		var buf []string
		for n := i; cities[i] == cities[n]; n++ {
			buf = append(buf, categories[n])
			if n == len(cities)-1 {
				break
			}

		}
		i += len(buf)
		coll[cities[i-len(buf)]] = buf
	}

	return coll
}
