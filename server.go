package main

func main() {
	e := initRouter()
	e.Logger.Fatal(e.Start(":1323"))
}
