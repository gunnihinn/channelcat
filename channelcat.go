/*
channelcat is a cat clone that uses channels to pipe strings from
standard input to standard output.

We use a scanner to read from a file or standard input.
We send the text from the scanner into a channel input.
We read from channel input and print it to standard output.
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanners := make(chan *bufio.Scanner) // Hold input files
	lines := make(chan string)            // Hold our messages
	done := make(chan int)                // Syncronize goroutines

	// If called with arguments, scan each file and feed scanner into
	// channel. Otherwise scan standard input and feed scanner into
	// channel.
	if len(os.Args) > 1 {
		go func() {
			for _, arg := range os.Args[1:] {
				file, err := os.Open(arg)
				if err != nil {
					log.Fatalf("Can't open %s: %s\n", arg, err)
				}
				scanners <- bufio.NewScanner(file)
			}
			close(scanners)
		}()
	} else {
		go func() {
			scanners <- bufio.NewScanner(os.Stdin)
			close(scanners)
		}()
	}

	// Pipe input from scanner to a channel
	go func() {
		for scanner := range scanners {
			for scanner.Scan() {
				lines <- scanner.Text()
			}
			if scanner.Err() != nil {
				log.Fatalf("Error during scanning: %s\n", scanner.Err())
			}
		}
		close(lines)
	}()

	// Read from a channel and print contents to Stdout
	go func() {
		for line := range lines {
			fmt.Println(line)
		}
		done <- 1
	}()

	// Don't quit until everything has been printed
	<-done
}
