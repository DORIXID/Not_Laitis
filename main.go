package main

import "fmt"

/*Работа с картами (map)

  Измени хранение команд: вместо []command используй map[string]command, где ключ — имя команды.
  Переделай addCommand, selectAllFrom, selectNFrom, чтобы они работали с картой.

*/

//CREATE TABLE command
type command struct {
	category    string
	args        []string
	description string
}

/*

var people = map[string]int{
    "Tom": 1,
    "Bob": 2,
    "Sam": 4,
    "Alice": 8,
}
fmt.Println(people)     // map[Tom:1 Bob:2 Sam:4 Alice:8]

*/

//MAIN
func main() {
	defer finish()                                            //DEFER - finish
	commands := make(map[string]command)                      //СОЗДЕМ, MAP С СТРУКТУРАМИ, ключем в которой служит имя команды
	commands = addCommand(commands, "Hello, world!", command{ //Создаём новую команду с помощью addCommand
		category:    "system",
		args:        []string{"echo Hello, world!"},
		description: "Выводит в командную строку Windows \"Hello world!\"",
	})
	commands = addCommand(commands, "dir", command{ //Создаём новую команду с помощью addCommand
		category:    "system",
		args:        []string{"dir"},
		description: "Выводит список файлов и папок текущей директории",
	})
	commands = addCommand(commands, "Ping Google", command{ //Создаём новую команду с помощью addCommand
		category:    "network",
		args:        []string{"ping google.com"},
		description: "Проверяет доступность Google.",
	})
	selectAllFrom(commands)
	filterByCategory(commands, "system")
}

//ФУНКЦИЯ ВЫБОРА команды - select i from command
func selectNFrom(mapp map[string]command, key string) {
	fmt.Printf("Имя команды: %s\n", key)
	fmt.Printf("Категория команды: %s\n", mapp[key].category)
	fmt.Printf("Аргументы: %s\n", mapp[key].args)
	fmt.Printf("Описание: %s\n", mapp[key].description)
	fmt.Printf("--------------------\n")
}

//ФУНКЦИЯ ВЫБОРА ВСЕХ КОМАНД - SELECT * FROM command
func selectAllFrom(mapp map[string]command) {
	i := 0
	for key, _ := range mapp {
		i++
		fmt.Printf("Команда №%d\n", i) //Индексация в командной строке
		selectNFrom(mapp, key)
	}
}

//ФУНКЦИЯ фильтрации команд по категориям
func filterByCategory(mapp map[string]command, category string) {
	fmt.Printf("\n--------------------\nФильрация по категории \"%s\":\n--------------------\n", category)
	for key, cmd := range mapp {
		if cmd.category == category {
			selectNFrom(mapp, key)
		}
	}
	fmt.Printf("✓-Фильтрация завершена-✓\n--------------------\n")
}

// ДОБАВЛЕНИЕ КОМАНДЫ В КАРТУ С КОМАНДАМИ
// берет на вход карту(map), структуру со значениями и название ключа для структуры в карте а возвращает обновлённую карту
func addCommand(mapp map[string]command, name string, cmd command) map[string]command {
	// Проверка на уникальность ключа
	if _, exists := mapp[name]; exists { //Здесь в if сначала выполняется инструкция - _, exists := mapp[name] - тут _ это значние, а exists это bool значение.
		fmt.Printf("\n\n❌--The entered name \"%s\" for the command is already in use--❌\n\n", name) //Оно указывает на то, получилась ли операция присвоения по данному ключу или нет, а потом уже смотрится на то true  или false у exists
		return mapp                                                                                 // Возвращаем неизмененную карту, если имя уже существует
	}
	mapp[name] = cmd
	return mapp // Если все хорошо добавляем в слайс новую структуру с командой
}

// DEFER-FUNC - выводит в терминал "Вывод завершён"
func finish() {
	fmt.Println("\n\n✓✓✓-Вывод завершён-✓✓✓")
}
