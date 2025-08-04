package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Card struct {
	Word         string
	Transalation string
}

type Deck struct {
	Name       string
	Cards      []Card
	TotalCards int
}

func NewDeck(name string, jsonFilePath string, tc int) (*Deck, error) {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var cards []Card
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return nil, err
	}

	deck := &Deck{
		Name:       name,
		Cards:      cards,
		TotalCards: tc,
	}

	return deck, nil
}

func (d *Deck) learnDeck() {
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)

	total := len(d.Cards)
	correct := 0

	shuffledCards := make([]Card, len(d.Cards))
	copy(shuffledCards, d.Cards)
	rand.Shuffle(len(shuffledCards), func(i, j int) {
		shuffledCards[i], shuffledCards[j] = shuffledCards[j], shuffledCards[i]
	})

	fmt.Printf("Учим колоду %s (%d карточек)", d.Name, d.TotalCards)
	fmt.Print("Введите enter для начала")
	fmt.Scan()

	for i, card := range shuffledCards {
		fmt.Printf("Карточка %d/%d\n", i+1, total)
		fmt.Printf("Слово: %s, введите перевод:\n", card.Word)

		scanner.Scan()
		answer := scanner.Text()

		if answer == card.Transalation {
			correct++
			fmt.Printf("Верно! %s = %s", card.Word, card.Transalation)
		} else {
			fmt.Printf("Неправильно! Правильный ответ: %s. Ваш Ответ: %s", card.Transalation, answer)
		}

	}
	fmt.Printf("\n🏁 Результат: %d/%d (%.0f%%)\n", correct, total, float64(correct)/float64(total)*100)
}

func main() {

}
