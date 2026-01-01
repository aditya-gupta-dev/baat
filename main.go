package main

func main() {
	config := DefaultConfig()

	room := CreateRoom(&config)
	room.StartListening()
}
