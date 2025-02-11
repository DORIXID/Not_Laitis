package main

import "fmt"

/*Задание: Расширенные структуры

  Добавь в command новое поле category string (например, "system", "network" и т. д.).
  Реализуй фильтрацию команд по категории (функция filterByCategory(commands []command, category string)).

*/

//CREATE TABLE command
type command struct {
	category    string
	name        string
	args        []string
	description string
}

//MAIN
func main() {
	defer finish()                           //DEFER - finish
	var commands []command = []command{}     //СОЗДЕМ, СЛАЙС С СТРУКТУРАМИ
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		category:    "system",
		name:        "Hello, world!",
		args:        []string{"echo Hello, world!"},
		description: "Выводит в командную строку Windows \"Hello world!\"",
	})
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		category:    "system",
		name:        "dir",
		args:        []string{"dir"},
		description: "Выводит список файлов и папок текущей директории",
	})
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		category:    "network",
		name:        "Ping Google",
		args:        []string{"ping google.com"},
		description: "Проверяет доступность Google.",
	})
	selectAllFrom(commands)
	filterByCategory(commands, "system")
}

//ФУНКЦИЯ ВЫБОРА команды - select i from command
func selectNFrom(cmd command) {
	fmt.Printf("Категория команды: %v\n", cmd.category)
	fmt.Printf("Имя команды: %v\n", cmd.name)
	fmt.Printf("Аргументы: %v\n", cmd.args)
	fmt.Printf("Описание: %v\n", cmd.description)
	fmt.Printf("--------------------\n")
}

//ФУНКЦИЯ ВЫБОРА ВСЕХ КОМАНД - SELECT * FROM command
func selectAllFrom(slice []command) {
	for i, cmd := range slice {
		fmt.Printf("Команда №%d\n", i+1) //Индексация в командной строке
		selectNFrom(cmd)
	}
}

//ФУНКЦИЯ фильтрации команд по категориям
func filterByCategory(slice []command, category string) {
	fmt.Printf("\n--------------------\nФильрация по категории \"%s\":\n--------------------\n", category)
	for _, cmd := range slice {
		if cmd.category == category {
			selectNFrom(cmd)
		}
	}
	fmt.Printf("✓-Фильтрация завершена-✓\n--------------------\n")
}

// ДОБАВЛЕНИЕ КОМАНДЫ В СЛАЙС СТРУКТУР С КОМАНДАМИ
// берет на вход слайс и структуру со значениями, а возвращает обновлённый слайс через return и append.
// если попробовать что то типа slice := append(slice, cmd), то создастся новый слайс, так как он получается локальный, а настоящий не поменяется.
func addCommand(slice []command, cmd command) []command {
	for _, existingCmd := range slice {
		if existingCmd.name == cmd.name {
			fmt.Printf("\n\n❌--The entered name \"%s\" for the command is already in use--❌\n\n", existingCmd.name)
			return (slice) // Возвращаем не обновленный слайс, если обнаружено повторение(return завершает выполнение функции)
		}
	}
	return append(slice, cmd) // Если все хорошо добавляем в слайс новую структуру с командой
}

// DEFER-FUNC - выводит "Вывод завершён"
func finish() {
	fmt.Println("\n\n✓✓✓-Вывод завершён-✓✓✓")
}
