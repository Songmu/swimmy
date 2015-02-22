package main

import "os"

func main() {
	s := NewSwimmy(Args{
		Dir:      "./test",
		Procs:    1,
		Interval: 1,
		APIKey:   os.Getenv("MACKEREL_API_KEY"),
		APIBase:  "https://mackerel.io",
		Debug:    false,
	})
	s.Run()
}
