package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

/*

 */

// Определение структуры для хранения команды из json файла
type commandStruct struct {
	Category    string   `json:"Category"`
	Args        []string `json:"Args"`
	Description string   `json:"Description"`
}

// Определение карты для хранения структур с командами из json файла
type commandMap map[string]commandStruct

// Переменная с картой со структурами, хранящими команды и их атрибуты
var commands commandMap

// Структура для данных регистрации или авторизации
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MAIN
func main() {
	//defer finish()

	// URL путь для регистрации и логина
	registerURL := "http://185.72.144.59:80/register"
	loginURL := "http://185.72.144.59:80/login"
	// URL для подключения к WebSocket серверу
	wsURL := "ws://185.72.144.59:80/ws"

	// При первом запуске файла создается пустой файл
	if _, err := os.OpenFile("commands.json", os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		fmt.Println("Ошибка при создании файла commands.json: ", err)
		return
	}

	// Загрузка существующих команд из файла
	commands = loadCommandsFromFile("commands.json")

	// Добавление новых команд
	commands = addCommand(commands, "Start message", commandStruct{
		Category:    "system",
		Args:        []string{"cmd", "/C", "msg * \"Выполнение команд работает корректно!\""},
		Description: "Выводит на экран сообщение об успешном запуске программы",
	})

	commands = addCommand(commands, "offComp", commandStruct{
		Category:    "system",
		Args:        []string{"cmd", "/C", "shutdown /s /t 60"},
		Description: "Завершает работу пк через 60 секунд",
	})

	//commands.Execute("Play_92")

	// Вывод всех команд
	//selectAllFrom(commands)

	// Фильтрация по категории
	//filterByCategory(commands, "system")

	// Сохранение обновленной карты команд в файл
	saveCommandsToFile("commands.json", commands)

	// Запуск выполнения приветственной команды
	//executeCommand("Message")

	//Загружаем токен
	token, err := loadString(".token")
	if err != nil {
		var resp string
		log.Println("⚠️ Токен не найден, требуется авторизация или регистрация")
		for {
			fmt.Print("\nВведите \"Р\", если хотите зарегистрироваться или \"Л\", если хотите авторизоваться: ")
			fmt.Fscan(os.Stdin, &resp)
			var login string
			var password string
			if resp == "Р" {
				fmt.Print("\nВведите логин: ")
				fmt.Fscan(os.Stdin, &login)
				fmt.Print("\nВведите пароль: ")
				fmt.Fscan(os.Stdin, &password)
				if err := registerUser(registerURL, login, password); err != nil {
					fmt.Printf("Ошибка при регистрации: %v", err)
					continue
				}
				fmt.Println("Вы успешно зарегистрировались, теперь необходимо авторизоваться с помощью вашего логина и пароля")
			} else if resp == "Л" {
				fmt.Print("\nВведите логин: ")
				fmt.Fscan(os.Stdin, &login)
				fmt.Print("\nВведите пароль: ")
				fmt.Fscan(os.Stdin, &password)
				if err := loginUser(loginURL, login, password); err != nil {
					fmt.Printf("Ошибка при авторизации: %v", err)
					continue
				}
				fmt.Println("Вы успешно вошли в систему")
				log.Println("🔑 Загружен токен:", token)
				// Загружаем из файла rKey и генерируем ссылку для удаленной активации команд
				rKey, err := loadString(".rKey")
				if err != nil {
					log.Println("Ошибка при загрузке rKey: ", err)
					return
				}
				fmt.Printf("Ваша ссылка для удалённого запуска команд: \nhttp://185.72.144.59:80/run?user=%s&cmd=ВашаКоманда\n", rKey)
				//Подключаемся по ws
				connectToWS(wsURL)
			} else {
				fmt.Print("\nНеверно введен ответ")
			}
		}
	} else {
		log.Println("🔑 Загружен токен:", token)
		// Загружаем из файла rKey и генерируем ссылку для удаленной активации команд
		rKey, err := loadString(".rKey")
		if err != nil {
			log.Println("Ошибка при загрузке rKey: ", err)
			return
		}
		fmt.Printf("Ваша ссылка для удалённого запуска команд: \nhttp://185.72.144.59:80/run?user=%s&cmd=ВашаКоманда\n", rKey)
		//Подключаемся по ws
		connectToWS(wsURL)
	}
}

// Функция для отправки запроса на регистрацию
func registerUser(url string, username string, password string) error {
	// Создаем структуру с данными регистрации
	reqData := RegisterRequest{
		Username: username,
		Password: password,
	}

	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	// Отправляем POST-запрос на сервер
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	// Закрываем ответ от сервера при завершении функции для избежания утечки памяти
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("ошибка регистрации: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Сохраняем randomKey в файл
	randomKey := string(body)
	fmt.Println("Регистрация успешна, получен randomKey:", randomKey)
	if err := saveString(randomKey, ".rKey"); err != nil {
		return fmt.Errorf("Ошибка сохранения rKey: %v", err)
	}

	return nil
}

// Функция для отправки запроса на автоизацию
func loginUser(url string, username string, password string) error {
	// Создаем структуру с данными авторизации(используя структуру для регистрации)
	reqData := RegisterRequest{
		Username: username,
		Password: password,
	}

	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	// Отправляем POST-запрос на сервер
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	//Закрываем(удаляем) соединение от сервера при завершении функции для избежания утечки памяти
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка Авторизации: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Десериализуем ответ в структуру
	var loginResp map[string]string
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("ошибка десериализации ответа: %v", err)
	}

	// Сохраняем токен
	token, ok := loginResp["token"]
	if !ok {
		return fmt.Errorf("токен не найден в ответе")
	}
	if err := saveString(token, ".token"); err != nil {
		return fmt.Errorf("ошибка сохранения токена: %v", err)
	}

	// Сохраняем токен
	rKey, ok := loginResp["rKey"]
	if !ok {
		return fmt.Errorf("токен не найден в ответе")
	}
	if err := saveString(rKey, ".rKey"); err != nil {
		return fmt.Errorf("ошибка сохранения токена: %v", err)
	}

	fmt.Println("Авторизация успешно пройдена, токен сохранен")
	return nil
}

// Сохранение string в файл
func saveString(varString string, filename string) error {
	// Создаем или записываем в существующий файл varString, преобразуя его в байтовый список
	return os.WriteFile(filename, []byte(varString), 0600) // Доступ только владельцу
}

// Загрузка string из файла
func loadString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Функция для подключения к WebSocket серверу с использованием JWT токена
func connectToWS(url string) {
	// Заголовки для аутентификации
	token, err := loadString(".token")
	if err != nil {
		log.Println("Ошибка при чтении токена: ", err)
		return
	}
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token) // Я кста не знаю есть ли смысл добавлять Bearer т к на стороне сервера он все равно игнорируется
	// беконечный цикл подключения к серверу и обработки команд при удачном подключении

	for {
		//Пробуем подключиться
		//Если подключиться не удается, то через 10 секунд запускаем цикл заново и пытаемся подключиться
		conn, _, err := websocket.DefaultDialer.Dial(url, headers)
		if err != nil {
			log.Printf("❌ Ошибка подключения к WebSocket: %v", err)
			log.Println("Попытка переподключения через 5 секунд")
			for i := 4; i != 0; i-- {
				fmt.Printf("%d...", i)
				time.Sleep(1 * time.Second)
			}
			fmt.Println("")
			continue
		}

		log.Println("✔️ Успешно подключено к WebSocket серверу")

		// Чтение сообщений от сервера
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Ошибка при получении сообщения:", err)
				conn.Close() // Явно закрываем перед новой попыткой подключения
				break
			}
			command := string(message)
			log.Println("Получена команда:", command)
			go executeCommand(command)
		}
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

// DEFER-FUNCTION - выводит в терминал "Программа завершена" и закрывает файл json
func finish() {
	fmt.Println("\n\n✓✓✓-Программа завершена-✓✓✓")
}
