package main

import "fmt"

type command struct { //CREATE TABLE command
	name        string
	args        []string
	discription string
}

func main() { //MAIN
	defer finish() //DEFER - finish
	var HelloWorld = command{
		name:        "Hello world!",
		args:        []string{"echo Hello, world!"},
		discription: "Выводит в командную строку Windows \"Hello world!\""}
	HelloWorld.selectallfrom()

}

func (cmd command) selectallfrom() { //МЕТОД select * from command
	fmt.Printf("\nИмя команды: %v", cmd.name)
	fmt.Printf("\nАргументы: %v", cmd.args)
	fmt.Printf("\nОписание: %v", cmd.discription)
}

func finish() { // DEFER-FUNC - выводит "Вывод завершён"
	fmt.Println("\nВывод завершён")
}
