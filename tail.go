package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hpcloud/tail"
	"golang.org/x/net/context"
)

func tailFile(ctx context.Context, file string, poll bool, dest *os.File) {
	defer wg.Done()

	s, err := os.Stat(file)
	if err != nil {
		log.Fatalf("unable to stat %s: %s", file, err)
	}

	t, err := tail.TailFile(file, tail.Config{
		Follow: true,
		ReOpen: true,
		Poll:   poll,
		Logger: tail.DiscardingLogger,
		Pipe:   s.Mode()&os.ModeNamedPipe != 0,
	})
	if err != nil {
		log.Fatalf("unable to tail %s: %s", "foo", err)
	}

	// main loop
	for {
		select {
		// if the channel is done, then exit the loop
		case <-ctx.Done():
			t.Stop()
			t.Cleanup()
			return
		// get the next log line and echo it out
		case line := <-t.Lines:
			if line != nil {
				fmt.Fprintln(dest, line.Text)
			}
		}
	}
}
