package main

import (
	"time"
	"log"
)

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %ss", name, elapsed.Seconds())
}

// example:
/*func factorial(n *big.Int) (result *big.Int) {
    defer timeTrack(time.Now(), "factorial")
    // ... do some things, maybe even return under some condition
    return n
}*/
