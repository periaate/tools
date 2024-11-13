package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
)

type fauxRC struct{ io.Reader }

func (f fauxRC) Close() error { return nil }

func main() {
	args := fsio.Args()
	if len(args) == 0 {
		clog.Fatal("not enough args")
	}

	b, err := os.ReadFile(args[0])
	if err != nil {
		clog.Fatal("could't open file", "err", err)
	}

	f := fauxRC{bytes.NewReader(b)}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	defer speaker.Close()

	n, err := strconv.Atoi(args[1])
	if err != nil {
		clog.Fatal("couldn't parse string", "err", err)
	}

	dur := time.Duration(n) * time.Minute
	tim := time.Now().Add(dur)
	for time.Now().Before(tim) {
		fmt.Printf("\r%v", time.Until(tim))
		time.Sleep(16 * time.Millisecond)
	}
	fmt.Printf("\r0.0s%v\r", strings.Repeat(" ", 10))

	speaker.Play(streamer)
	time.Sleep(5 * time.Second)
}
