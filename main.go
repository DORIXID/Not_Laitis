package main

import "fmt"

/*Задание: Управление списком команд

Создай слайс структур для хранения нескольких команд. Реализуй две функции:

    addCommand(slice []command, cmd command) []command: добавляет новую команду в слайс и возвращает обновленный слайс.
    printCommands(slice []command): выводит список всех команд с их деталями.

Результат:

    Программа должна уметь:
        Добавлять команды в список.
        Выводить все команды из списка.
    Попробуй создать 2–3 команды (например, ping google.com и dir) и добавить их в список.
*/
//CREATE TABLE command
type command struct {
	name        string
	args        []string
	description string
}

//MAIN
func main() {
	defer finish()                           //DEFER - finish
	var commands []command = []command{}     //СОЗДЕМ, СЛАЙС С СТРУКТУРАМИ
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		name:        "Hello, world!",
		args:        []string{"echo Hello, world!"},
		description: "Выводит в командную строку Windows \"Hello world!\"",
	})
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		name:        "Hello, world!",
		args:        []string{"dir"},
		description: "Выводит список файлов и папок текущей директории",
	})
	commands = addCommand(commands, command{ //Создаём новую команду с помощью addCommand
		name:        "Ping Google",
		args:        []string{"ping google.com"},
		description: "Проверяет доступность Google.",
	})
	selectAllFrom(commands)
}

//ФУНКЦИЯ ВЫБОРА 1 команды - select N from command
func selectNFrom(cmd command) {
	fmt.Printf("Имя команды: %v\n", cmd.name)
	fmt.Printf("Аргументы: %v\n", cmd.args)
	fmt.Printf("Описание: %v\n", cmd.description)
	fmt.Printf("[--------------------]\n")
}

//ФУНКЦИЯ ВЫБОРА ВСЕХ КОМАНД - SELECT * FROM command
func selectAllFrom(slice []command) {
	for i, cmd := range slice {
		fmt.Printf("Команда №%d\n", i+1)
		fmt.Printf("Имя команды: %v\n", cmd.name)
		fmt.Printf("Аргументы: %v\n", cmd.args)
		fmt.Printf("Описание: %v\n", cmd.description)
		fmt.Printf("[--------------------]\n")
	}
}

// ДОБАВЛЕНИЕ КОМАНДЫ В СЛАЙС СТРУКТУР С КОМАНДАМИ
// берет на вход слайс и структуру со значениями, а возвращает обновлённый слайс через return и append.
// если попробовать что то типа slice := append(slice, cmd), то создастся новый слайс, так как он получается локальный, а настоящий не поменяется.
func addCommand(slice []command, cmd command) []command {
	for _, existingCmd := range slice {
		if existingCmd.name == cmd.name {
			fmt.Printf("\n\n❌--The entered name \"%s\" for the command is already in use--❌\n\n", existingCmd.name)
			return (slice)
		}
	}
	return append(slice, cmd)
}

// DEFER-FUNC - выводит "Вывод завершён"
func finish() {
	fmt.Println("\n\n✓✓✓Вывод завершён✓✓✓")
}
