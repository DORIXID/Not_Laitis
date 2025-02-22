package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

/*

 */

// CREATE TABLE command
type commandStruct struct {
	Category    string   `json:"Category"`
	Args        []string `json:"Args"`
	Description string   `json:"Description"`
}
type commandMap map[string]commandStruct

var commands commandMap

// MAIN
func main() {
	//defer finish()

	// Загрузка существующих команд из файла
	commands = loadCommandsFromFile("commands.json")

	// Добавление новых команд
	commands = addCommand(commands, "Hello, world!", commandStruct{
		Category:    "system",
		Args:        []string{"cmd", "/C", "echo Hello, world!"},
		Description: "Выводит в командную строку Windows \"Hello world!\"",
	})

	commands = addCommand(commands, "Fimoz", commandStruct{
		Category:    "media",
		Args:        []string{"cmd", "/C", "start", "Z:\\No_Laitis\\O_HET_FIMOZ.mp3", "&&", "start", "Z:\\No_Laitis\\ShUE_PPSh.gif"},
		Description: "SHUE PPsh",
	})

	//commands.Execute("Play_92")

	// Вывод всех команд
	//selectAllFrom(commands)

	// Фильтрация по категории
	//filterByCategory(commands, "system")

	// Сохранение обновленной карты команд в файл
	saveCommandsToFile("commands.json", commands)

	// Запуск выполнения приветственной команды
	executeCommand("Message")

	//Подключение к серверу и запуск прослушки и выполнения команд от сервера
	connectToWebSocket()

}

// Подключение к WebSocket серверу и обработка команд
func connectToWebSocket() {
	var conn *websocket.Conn
	var err error

	// Устанавливаем соединение с сервером
	for {
		conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8282/ws", nil)
		if err != nil {
			for {
				// Цикл попыток подключения к серверу
				log.Printf("❌ Ошибка подключения к WebSocket: %v", err)
				log.Println("Попытка переподключения через 10 секунд")
				// Отсчет 10 секунд
				for i := 1; i <= 10; i++ {
					time.Sleep(1 * time.Second)
					fmt.Printf("%d...", i)
				}
				fmt.Println("\nПопытка переподключения...")
				break
			}
		} else {
			break
		}
	}
	// Закрываем соединение в случае выхода из цикла ожидания сообщений(если программа завершится пользователем)
	defer conn.Close()
	log.Println("✔️ Подключено к WebSocket серверу")
	// Цикл ожидания сообщений от сервера
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка при получении сообщения:", err)
			break
		}
		command := string(message)
		log.Println("Получена команда:", command)
		// Выполняем команду
		executeCommand(command)
	}
}

// Выполнение переданной команды в cmd
func executeCommand(command string) {
	cmdInfo, ok := commands[command]
	if !ok {
		fmt.Printf("\nКоманда \"%s\" не найдена\n", command)
		return
	}
	if len(cmdInfo.Args) == 0 {
		fmt.Printf("В команде \"%s\" нет аргументов", command)
		return
	}
	log.Printf("Выполнение команды: %s", command)
	cmd := exec.Command(cmdInfo.Args[0], cmdInfo.Args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка выполнения команды: %v\n", err)
		log.Printf("Вывод команды: %s\n", output)
		return
	}
	log.Printf("Результат:\n%s\n", output)
}

// ФУНКЦИЯ для обработки входящего соединения. Из входящего запроса от сервера получает данные и передает их в command.Execute для выполнения
func handleConnection(conn net.Conn, commands commandMap) {
	// Создаем буферизованный читатель для чтения данных из соединения
	reader := bufio.NewReader(conn)
	// Читаем строку до символа новой строки
	commandFromServer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	// Удаляем символ новой строки из полученного сообщения
	commandFromServer = commandFromServer[:len(commandFromServer)-1]
	// Выводим полученное сообщение в консоль
	fmt.Println(commandFromServer)
	// Закрываем соединение
	conn.Close()
	// Выполняем принятую команду
	commands.Execute(commandFromServer)
}

// ФУНКЦИЯ для сохранения карты(map) с командами в json файле
func saveCommandsToFile(filename string, commands commandMap) {
	// Сериализация данных в json
	jsonData, err := json.MarshalIndent(commands, "", " ")
	if err != nil {
		fmt.Println("Ошибка при конертации данных в json: ", err)
		return
	}

	// Открытие файла для записи или создание файла для записи, если он еще не создан
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // os.O_WRONLY - только запись, os.O_CREATE - создание файла при его отсутствии, os.O_TRUNC - стирание данных файла
	if err != nil {
		fmt.Println("Ошибка при открытии файла: ", err)
		return
	}

	//Запись данных в файл
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Ошибка при записи данных в json: ", err)
		return
	}
	defer file.Close()
}

// ФУНКЦИЯ загрузки команд из json файла в карту(map)
func loadCommandsFromFile(filename string) commandMap {
	// Открытие файла для считывания
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла: ", err)
		return make(commandMap)
	}
	defer file.Close() //Закрываем файл

	// загрузка данных из файла в переменную jsonData
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка при чтении файла: ", err)
		return make(commandMap)
	}

	// Десреаилизация данных из jsonData
	var commands commandMap
	err = json.Unmarshal(jsonData, &commands)
	if err != nil {
		fmt.Println("Ошибка при десериализации двнных: ", err)
		return make(commandMap)
	}
	return commands //Возвращаем карту
}

// ФУНКЦИЯ ВЫБОРА команды - select i from command
func selectIFrom(mapp commandMap, key string) {
	fmt.Println("--------------------")
	fmt.Printf("Имя команды: %s\n", key)
	fmt.Printf("Категория: %s\n", mapp[key].Category)
	fmt.Printf("Аргументы: %s\n", mapp[key].Args)
	fmt.Printf("Описание: %s\n", mapp[key].Description)
	fmt.Println("--------------------")
}

// ФУНКЦИЯ ВЫБОРА ВСЕХ КОМАНД - SELECT * FROM command
func selectAllFrom(commands commandMap) {
	i := 0
	for key, _ := range commands {
		i++
		fmt.Printf("Команда №%d\n", i) //Индексация в командной строке
		selectIFrom(commands, key)
	}
}

// ФУНКЦИЯ фильтрации команд по категориям
func filterByCategory(mapp commandMap, category string) {
	fmt.Printf("\n--------------------\nФильрация по категории \"%s\":\n--------------------\n", category)
	for key, cmd := range mapp {
		if cmd.Category == category {
			selectIFrom(mapp, key)
		}
	}
	fmt.Printf("✓-Фильтрация завершена-✓\n--------------------\n")
}

// ДОБАВЛЕНИЕ КОМАНДЫ В КАРТУ С КОМАНДАМИ
// берет на вход карту(map), которую нужно изменить, затем берет название для ключа новой карты, затем берет данные для новой карты, а в конце возвращает обновлённую карту
func addCommand(mapp commandMap, name string, cmd commandStruct) commandMap {
	// Проверка на уникальность ключа
	if _, exists := mapp[name]; exists { //Здесь в if сначала выполняется инструкция - _, exists := mapp[name] - тут _ это значние, а exists это bool значение.
		fmt.Printf("\n\nThe \"%s\" command is already in use.\n\n", name) //Оно указывает на то, получилась ли операция присвоения по данному ключу или нет, а потом уже смотрится на то true  или false у exists
		return mapp                                                       // Возвращаем неизмененную карту, если имя уже существует
	}
	mapp[name] = cmd
	return mapp // Если все хорошо добавляем в слайс новую структуру с командой
}

// МЕТОД для выполнения команды
// Берет на вход карту с командами, а также название команды(ключ карты)
// Выполняет аргументы выбранной команды поочередности
func (mapp commandMap) Execute(key string) {
	// Проверка на существование ключа
	cmdInfo, ok := mapp[key]
	if !ok {
		fmt.Printf("\nКоманда \"%s\" не найдена\n", key)
		return
	}
	if len(cmdInfo.Args) == 0 {
		fmt.Printf("В команде \"%s\" нет аргументов", key)
		return
	}
	// Создание команды с аргументами и её выполнение
	cmd := exec.Command(cmdInfo.Args[0], cmdInfo.Args[1:]...) // Создание команды с процессом cmdInfo.Args[0] и флагами и аргументами процесса cmdInfo.Args[1:]...
	stdoutStderr, err := cmd.CombinedOutput()                 // Выполнение команды с помошью cmd.CombinedOutput(), а также как стандартный вывод, так и стандартный вывод ошибок в одно место(stdoutStderr),
	if err != nil {                                           // а если команда завершится с ошибкой, переменная err будет содержать информацию об ошибке.
		fmt.Printf("Ошибка при выполнении команды: %v\n", err)
		return
	}
	fmt.Printf("Результат команды \"%s\": \n%s\n", key, stdoutStderr)
}

// ИНТЕРФЕЙС для выполнимых структур
type Executable interface {
	Execute(key string)
}

///func (имя_параметра тип_получателя) имя_метода (параметры) (типы_возвращаемых_результатов){
///    тело_метода
///}

// DEFER-FUNCTION - выводит в терминал "Программа завершена" и закрывает файл json
func finish() {
	fmt.Println("\n\n✓✓✓-Программа завершена-✓✓✓")
}
