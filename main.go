package main

import "time"

func main() {
	scheduler := CheckScheduler{
		Checks: []CheckSettings{
			CheckSettings{
				Name:     "Ping google",
				Command:  "ping",
				Args:     []string{"-c 1", "www.google.com"},
				Env:      []string{""},
				Dir:      "",
				Interval: 5 * time.Second,
				Timeout:  5 * time.Second,
			}, CheckSettings{
				Name:     "Ping facebook",
				Command:  "ping",
				Args:     []string{"-c 1", "www.facebook.com"},
				Env:      []string{""},
				Dir:      "",
				Interval: 10 * time.Second,
				Timeout:  5 * time.Second,
			},
		},
	}

	scheduler.Start()

	noExit := make(chan bool)
	<-noExit

	// var d net.Dialer
	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// conn, err := d.DialContext(ctx, "tcp", "localhost:12345")
	// if err != nil {
	// 	log.Fatalf("Failed to dial: %v", err)
	// }
	// defer conn.Close()

	// if _, err := conn.Write([]byte("Hello, World!")); err != nil {
	// 	log.Fatal(err)
	// }
}
