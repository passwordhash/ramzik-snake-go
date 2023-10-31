package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	tb "github.com/nsf/termbox-go"
)

func main() {
	input, err := os.Open("config.txt") // Открытие файла
	if err != nil {
		println(err)
		log.Fatal("Файла нет")
	}
	defer input.Close()
	configFile := bufio.NewScanner(input) // Инициализация сканера.
	configFile.Scan()
	deskSize, err := strconv.Atoi(configFile.Text())
	if err != nil {
		log.Fatal("Хопа! А я не могу прочитать, что написано в конфиг файле")
	}
	configFile.Scan()
	snakeSpeed, errSpeed := strconv.ParseFloat(configFile.Text(), 32) // При snakeSpeed >= 2 начинается сущий ад)
	if errSpeed != nil {
		log.Fatal("Хопа! А я не могу прочитать, что написано в конфиг файле")
	}
	// Создаём двумерный слайс
	playground := make([][]string, deskSize)
	for i := range playground {
		playground[i] = make([]string, deskSize)
	}
	spaceSymbol := "."
	// Заполняем двумерный слайс, сначала пробелы
	for i := 0; i < deskSize; i++ {
		for j := 0; j < deskSize; j++ {
			playground[i][j] = spaceSymbol
		}
	}
	snakeCord := make([][]int, 3) // Двумерный слайс змейки, каждый слайс содержит вертикальную и горизонтальную координату
	for i := range snakeCord {
		snakeCord[i] = make([]int, 2)
	}
	appleCord := make([]int, 2)
	snakeDirectionHorizontal := 1
	snakeDirectionVertical := 0
	gameOver := false
	applesEaten := 2
	keyboardErr := tb.Init()
	if keyboardErr != nil {
		panic(keyboardErr)
	}
	defer tb.Close()
	// Задаём координаты голове змеи и яблока
	for i := 0; i < applesEaten+1; i++ {
		for j := 0; j < 2; j++ {
			snakeCord[i][j] = deskSize / 2
		}
	}
	for i := 0; i < 2; i++ {
		appleCord[i] = rand.Intn(deskSize-1) + 0
	}
	// Если координаты головы змеи совпадают с яблоком, то перемещаем яблоко
	for appleCord[0] == snakeCord[0][0] && appleCord[1] == snakeCord[0][1] {
		appleCord[0] = rand.Intn(deskSize-1) + 0
		appleCord[1] = rand.Intn(deskSize-1) + 0
	}
	playground[snakeCord[0][1]][snakeCord[0][0]] = "*"
	playground[snakeCord[1][1]][snakeCord[1][0]] = "*"
	playground[snakeCord[2][1]][snakeCord[2][0]] = "*"
	playground[appleCord[0]][appleCord[1]] = "X"
	go readKey(&snakeDirectionHorizontal, &snakeDirectionVertical)
	for { // for {} == while True. Постоянный цикл
		// Координаты каждой клетки змейки кроме первой приравниваем к предыдущей
		// Первую клетку двигаем вперёд
		// Отрисовываем каждую клетку
		for i := 0; i < applesEaten; i++ {
			snakeCord[applesEaten-i][0], snakeCord[applesEaten-i][1] = snakeCord[applesEaten-i-1][0], snakeCord[applesEaten-i-1][1]
			playground[snakeCord[applesEaten-i][1]][snakeCord[applesEaten-i][0]] = "*"
		}
		playground[snakeCord[applesEaten][1]][snakeCord[applesEaten][0]] = spaceSymbol
		if snakeCord[0][1]+snakeDirectionVertical == -1 || snakeCord[0][1]+snakeDirectionVertical == 8 || snakeCord[0][0]+snakeDirectionHorizontal == -1 || snakeCord[0][0]+snakeDirectionHorizontal == 8 {
			gameOver = true
		}
		if !gameOver {
			snakeCord[0][1], snakeCord[0][0] = snakeCord[0][1]+snakeDirectionVertical, snakeCord[0][0]+snakeDirectionHorizontal
		}
		if playground[snakeCord[0][1]][snakeCord[0][0]] == "*" {
			gameOver = true
		}
		playground[snakeCord[0][1]][snakeCord[0][0]] = "*"
		// Захавал яблоко. Делаем новое
		if snakeCord[0][1] == appleCord[0] && snakeCord[0][0] == appleCord[1] {
			applesEaten = applesEaten + 1
			snakeCordAdd := []int{snakeCord[applesEaten-1][1] - snakeDirectionVertical, snakeCord[applesEaten-1][0] - snakeDirectionHorizontal}
			snakeCord = append(snakeCord, snakeCordAdd)
			appleCord[0] = rand.Intn(deskSize-1) + 0
			appleCord[1] = rand.Intn(deskSize-1) + 0
			for i := 0; i < applesEaten; i++ {
				for appleCord[0] == snakeCord[i][1] && appleCord[1] == snakeCord[i][0] {
					// Если новые координаты яблока совпадают с телом змеи, то яблоко нужно пересоздать
					appleCord[0] = rand.Intn(deskSize-1) + 0
					appleCord[1] = rand.Intn(deskSize-1) + 0
				}
			}
			playground[appleCord[0]][appleCord[1]] = "X"
		}
		for k := 0; k < deskSize; k++ { // Вывод матрицы в терминал
			for l := 0; l < deskSize; l++ {
				fmt.Print(playground[k][l], " ")
			}
			fmt.Println()
		}
		if gameOver {
			for k := 0; k < deskSize; k++ {
				fmt.Printf("\033[1A\033[K")
			}
			fmt.Println("Game Over")
			break
		} else {
			time.Sleep(time.Duration(snakeSpeed) * time.Millisecond)
			for k := 0; k < deskSize; k++ {
				fmt.Printf("\033[1A\033[K")
			}
		}
	}
}

func readKey(horizAddress *int, vertAddress *int) { // Чтение инпута с клавиатуры. Ненавижу
	for {
		event := tb.PollEvent()
		switch {
		case event.Ch == 'a':
			if *horizAddress == 0 {
				*horizAddress = -1
				*vertAddress = 0
			}
		case event.Ch == 's':
			if *vertAddress == 0 {
				*horizAddress = 0
				*vertAddress = 1
			}
		case event.Ch == 'd':
			if *horizAddress == 0 {
				*horizAddress = 1
				*vertAddress = 0
			}
		case event.Ch == 'w':
			if *vertAddress == 0 {
				*horizAddress = 0
				*vertAddress = -1
			}
		}
	}
}
