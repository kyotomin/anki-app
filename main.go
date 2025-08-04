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
	incorrect := []Card{}
	startTime := time.Now()

	shuffledCards := make([]Card, len(d.Cards))
	copy(shuffledCards, d.Cards)
	rand.Shuffle(len(shuffledCards), func(i, j int) {
		shuffledCards[i], shuffledCards[j] = shuffledCards[j], shuffledCards[i]
	})

	fmt.Printf("Учим колоду %s (%d карточек)", d.Name, d.TotalCards)
	fmt.Print("Введите enter для начала, q для выхода")
	fmt.Scan()

	for i, card := range shuffledCards {
		fmt.Printf("Карточка %d/%d\n", i+1, total)
		fmt.Printf("Слово: %s, введите перевод:\n", card.Word)

		scanner.Scan()
		answer := scanner.Text()

		if answer == card.Transalation {
			correct++
			fmt.Printf("Верно! %s = %s", card.Word, card.Transalation)
		} else if answer == "q" {
			break
		} else {
			fmt.Printf("Неправильно! Правильный ответ: %s. Ваш Ответ: %s", card.Transalation, answer)
			incorrect = append(incorrect, card)
		}

		if len(incorrect) > 0 {
			fmt.Println("--- Повторяем ошибки ---")
			for _, card := range incorrect {
				fmt.Printf("Слово: %s, введите перевод:", card.Word)
				scanner.Scan()
				if scanner.Text() == card.Transalation {
					fmt.Printf("Верно! %s = %s", card.Word, card.Transalation)
					correct++
				}
			}
		}

	}

	duration := time.Since(startTime).Round(time.Second)
	fmt.Printf("\n🏁 Результат: %d/%d (%.0f%%) Длительность: %v\n, ", correct, total, float64(correct)/float64(total)*100, duration)
}

func main() {

}
