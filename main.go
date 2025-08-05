package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Card struct {
	Word        string `json:"Word"`
	Translation string `json: "Translation"`
}

type StudyCard struct {
	Question string
	Answer   string
	Original Card
}

type Deck struct {
	Name  string
	Cards []Card
}

func NewDeck(name string, jsonFilePath string) (*Deck, error) {
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
		Name:  name,
		Cards: cards,
	}

	return deck, nil
}

func (d *Deck) createStudyCards(mode string) []StudyCard {
	var studyCards []StudyCard

	for _, card := range d.Cards {
		switch mode {
		case "en-ru":
			studyCards = append(studyCards, StudyCard{
				Question: card.Word,
				Answer:   card.Translation,
				Original: card,
			})
		case "ru-en":
			studyCards = append(studyCards, StudyCard{
				Question: card.Translation,
				Answer:   card.Word,
				Original: card,
			})
		case "both":
			studyCards = append(studyCards, StudyCard{
				Question: card.Word,
				Answer:   card.Translation,
				Original: card,
			})
			studyCards = append(studyCards, StudyCard{
				Question: card.Translation,
				Answer:   card.Word,
				Original: card,
			})
		}
	}

	return studyCards
}

func (d *Deck) selectMode() (mode string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Выберите режим обучения:")
	fmt.Println("1. Английский → Русский (en-ru)")
	fmt.Println("2. Русский → Английский (ru-en)")
	fmt.Println("3. Оба направления (both)")
	fmt.Print("Введите номер (1-3): ")

	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		return "en-ru"
	case "2":
		return "ru-en"
	case "3":
		return "both"
	default:
		fmt.Println("Неверный выбор, используем режим 'Английский → Русский'")
		return "en-ru"
	}
}

func (d *Deck) learnDeck() {
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)

	mode := d.selectMode()

	studyCards := d.createStudyCards(mode)

	total := len(studyCards)
	correct := 0
	incorrect := []StudyCard{}
	startTime := time.Now()

	rand.Shuffle(len(studyCards), func(i, j int) {
		studyCards[i], studyCards[j] = studyCards[j], studyCards[i]
	})

	modeNames := map[string]string{
		"en-ru": "Английский - Русский",
		"ru-en": "Русский - Английский",
		"both":  "Английский + Русский",
	}

	fmt.Printf("Учим колоду %s в режиме %s (%d карточек)\n", d.Name, modeNames[mode], total)
	fmt.Print("Введите enter для начала, q для выхода\n")
	fmt.Scan()

	for i, card := range studyCards {
		fmt.Printf("Карточка %d/%d\n", i+1, total)
		fmt.Printf("Слово: %s, введите перевод:\n", card.Question)

		scanner.Scan()
		answer := scanner.Text()

		if strings.ToLower(answer) == card.Answer {
			correct++
			fmt.Printf("Верно! %s = %s\n", card.Question, card.Answer)
		} else if strings.ToLower(answer) == "q" {
			break
		} else {
			fmt.Printf("Неправильно! Правильный ответ: %s. Ваш Ответ: %s", card.Answer, answer)
			incorrect = append(incorrect, card)
		}

		if len(incorrect) > 0 {
			fmt.Println("--- Повторяем ошибки ---")
			for _, card := range incorrect {
				fmt.Printf("Слово: %s, введите перевод:", card.Question)
				scanner.Scan()
				if strings.ToLower(scanner.Text()) == card.Answer {
					fmt.Printf("Верно! %s = %s", card.Question, card.Answer)
				}
			}
		}

	}

	duration := time.Since(startTime).Round(time.Second)
	accuracy := float64(correct) / float64(total) * 100
	fmt.Printf("\n🏁 Результат: %d/%d (%.0f%%) Длительность: %v\n, ", correct, total, accuracy, duration)
}

func showUserProfile() {
	fmt.Print("In progress")
}

func addDeck() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n➕ Добавление новой колоды")
	fmt.Print("Введите имя колоды для отображения (q для выхода): ")
	scanner.Scan()
	if strings.ToLower(strings.TrimSpace(scanner.Text())) == "q" {
		return
	}
	name := strings.TrimSpace(scanner.Text())

	fmt.Print("Введите путь к JSON файлу (например: decks/words.json): ")
	scanner.Scan()
	if strings.ToLower(strings.TrimSpace(scanner.Text())) == "q" {
		return
	}
	filePath := strings.TrimSpace(scanner.Text())

	deck, err := NewDeck(name, filePath)
	if err != nil {
		fmt.Printf("❌ Ошибка загрузки колоды: %v\n", err)
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		showMainMenu()
	}

	fmt.Printf("✅ Колода '%s' успешно загружена! (%d карточек)\n", deck.Name, len(deck.Cards))
	fmt.Print("Введите enter для продолжения: ")
	fmt.Scanln()
}

func processLearnDeck() {
	fmt.Println("\n📚 Выбор колоды для изучения")

	decksDir := "decks"
	if _, err := os.Stat(decksDir); os.IsNotExist(err) {
		fmt.Printf("❌ Папка '%s' не найдена. Создайте папку и добавьте JSON файлы с колодами.\n", decksDir)
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	var jsonFiles []string
	err := filepath.WalkDir(decksDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Ошибка чтения папки decks: %v\n", err)
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	if len(jsonFiles) == 0 {
		fmt.Printf("❌ JSON файлы не найдены в папке '%s'\n", decksDir)
		fmt.Println("Создайте файлы с колодами в формате JSON")
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	fmt.Println("Доступные колоды:")
	for i, file := range jsonFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	fmt.Printf("%d. Вернуться в главное меню\n", len(jsonFiles)+1)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Выберите номер колоды: ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	var selectedIndex int
	_, err = fmt.Sscanf(choice, "%d", &selectedIndex)
	if err != nil || selectedIndex < 1 || selectedIndex > len(jsonFiles)+1 {
		fmt.Println("❌ Неверный выбор")
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	if selectedIndex == len(jsonFiles)+1 {
		return
	}

	selectedFile := jsonFiles[selectedIndex-1]
	deckName := strings.TrimSuffix(filepath.Base(selectedFile), ".json")

	deck, err := NewDeck(deckName, selectedFile)
	if err != nil {
		fmt.Printf("❌ Ошибка загрузки колоды: %v\n", err)
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	if len(deck.Cards) == 0 {
		fmt.Println("❌ Колода пуста")
		fmt.Print("Введите enter для продолжения: ")
		fmt.Scanln()
		return
	}

	deck.learnDeck()
}

func showMainMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n" + strings.Repeat("=", 40))
		fmt.Println("🎓 Добро пожаловать в Anki-Like-App!")
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println("Выберите опцию:")
		fmt.Println("1. 📊 Профиль (В разработке)")
		fmt.Println("2. ➕ Добавить колоду")
		fmt.Println("3. 📚 Учить колоды")
		fmt.Println("4. 🚪 Выход")
		fmt.Print("Введите номер (1-4): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showUserProfile()
		case "2":
			addDeck()
		case "3":
			processLearnDeck()
		case "4":
			fmt.Println("👋 До свидания!")
			return
		default:
			fmt.Println("❌ Неверный выбор. Попробуйте еще раз.")
			fmt.Print("Введите enter для продолжения: ")
			fmt.Scanln()
		}
	}
}

func main() {
	showMainMenu()
}
