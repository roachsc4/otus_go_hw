package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const TimeLayout = "2006-01-02 03:04:05 -0700 MST"

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("current time:", time.Now().Format(TimeLayout))
	fmt.Println("exact time:", ntpTime.Format(TimeLayout))
}
