package main

import "os"

func main() {
	s := NewSwimmy(Args{
		Dir:      "./test",
		Procs:    1,
		Interval: 1,
		ApiKey:   os.Getenv("MACKEREL_API_KEY"),
		ApiBase:  "https://mackerel.io",
		Debug:    true,
	})
	s.Run()
}
