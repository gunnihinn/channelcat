# Channelcat

Channelcat is a cat clone that uses channels to print its input to
standard output.  It's a silly thing that's teaching me how to use
goroutines and channels.

## Installation

    go get github.com/gunnihinn/channelcat
    go install github.com/gunnihinn/channelcat

## Use

    channelcat <file1> ... <fileN>
    <a Unix utility> | channelcat

## Problems

Channelcat fails in an interesting way when it can't open one of the
input files we pass it; it prints a few lines of the files before that
one, then exits with an error before finishing to print all the files.

This is a consequence of the naive concurrent way that channelcat is
built: It's internal pipeline is in three pieces:

    Scanners -> Input lines -> Output

That is, channelcat uses a goroutine to it first loops through its input
files, make a Scanner for each one and feed that scanner into a channel.
Another goroutine then receives from the scanner channel, loops through
the input lines of each scanner, and sends them to an output channel.
Finally, one last goroutine then loops over that channel and prints
whatever it finds there.

Thus, when one of the input files can't be opened, the error signaling
this gets processed while the previous file is still being printed and
it gets cut short by the fatal error message we emit.
