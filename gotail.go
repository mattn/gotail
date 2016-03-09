package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/daviddengcn/go-colortext"
	goenc "github.com/mattn/go-encoding"
	"golang.org/x/text/encoding"
)

var colors = []ct.Color{
	ct.White,
	ct.Green,
	ct.Cyan,
	ct.Magenta,
	ct.Yellow,
	ct.Blue,
	ct.Red,
}
var ci int

var mutex sync.Mutex
var decoder *encoding.Decoder
var enc = flag.String("e", "", "Decode encoding")

func tail(in io.Reader, out io.Writer, follow bool) (err error) {
	color := colors[ci]
	if ci++; ci >= len(colors) {
		ci = 0
	}
	buf := bufio.NewReader(in)
	for {
		b, _, err := buf.ReadLine()
		if len(b) > 0 {
			if decoder != nil {
				var d []byte
				if d, err = decoder.Bytes(b); err == nil {
					b = d
				}
			}
			mutex.Lock()
			ct.ChangeColor(color, false, ct.None, false)
			fmt.Fprintln(out, string(b))
			ct.ResetColor()
			mutex.Unlock()
		}
		if err != nil {
			if err != io.EOF || !follow {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	return err
}

func main() {
	flag.Parse()

	if *enc != "" {
		decoder = goenc.GetEncoding(*enc).NewDecoder()
	}

	if flag.NArg() == 0 {
		if err := tail(os.Stdin, os.Stdout, false); err != nil {
			fmt.Fprintln(os.Stderr, "gotail: "+err.Error())
		}
	} else {
		var wg sync.WaitGroup
		for _, arg := range flag.Args() {
			var in io.Reader
			if arg != "-" {
				fin, err := os.Open(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "gotail: "+err.Error())
					os.Exit(1)
				}
				defer fin.Close()
				fin.Seek(0, os.SEEK_END)
				in = fin
			} else {
				in = os.Stdin
			}
			wg.Add(1)
			go func(in io.Reader) {
				if err := tail(in, os.Stdout, true); err != nil {
					fmt.Fprintln(os.Stderr, "gotail: "+err.Error())
				}
				wg.Done()
			}(in)
		}
		wg.Wait()
	}
}
