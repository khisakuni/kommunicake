package main

import "github.com/khisakuni/kommunicake/app"

func main() {
	a, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	a.Run()
}
