package main

func main() {
	server := NewRouterServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
