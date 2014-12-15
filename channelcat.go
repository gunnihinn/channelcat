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
	"os"
)

func main() {
	lines := make(chan string) // Hold our messages
	done := make(chan int)     // Syncronize goroutines

	scanner := bufio.NewScanner(os.Stdin)

	// Pipe input from scanner to a channel
	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
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
