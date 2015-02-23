package main

import "os"

func main() {
	// TODO
	s := newSwimmy(args{
		service:  "test1",
		dir:      "./test",
		procs:    1,
		interval: 1,
		apiKey:   os.Getenv("MACKEREL_API_KEY"),
		apiBase:  "https://mackerel.io",
		// debug:    false,
	})
	s.swim()
}
