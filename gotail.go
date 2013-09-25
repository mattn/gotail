package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
	"os"
	"sync"
)

var colors = []ct.Color{
	ct.Green,
	ct.Cyan,
	ct.Magenta,
	ct.Yellow,
	ct.Blue,
	ct.Red,
}
var ci int

var mutex sync.Mutex

func relay(color ct.Color, in io.Reader, out io.Writer) (err error) {
	var n int
	b := make([]byte, 4096)
	for {
		n, err = in.Read(b)
		if err != nil {
			break
		}
		mutex.Lock()
		ct.ChangeColor(color, false, ct.None, false)
		_, err = out.Write(b[:n])
		ct.ResetColor()
		mutex.Unlock()
		if err != nil {
			break
		}
	}
	return
}

func tail(in io.Reader) {
	var err error
	color := colors[ci]
	if ci++; ci >= len(colors) {
		ci = 0
	}
	for {
		err = relay(color, in, os.Stdout)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, "gotail: " + err.Error())
				break
			}
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		tail(os.Stdin)
	} else {
		for _, arg := range os.Args[1:] {
			var in io.Reader
			if arg != "-" {
				fin, err := os.Open(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "gotail: " + err.Error())
					os.Exit(1)
				}
				defer fin.Close()
				fin.Seek(0, os.SEEK_END)
				in = fin
			} else {
				in = os.Stdin
			}
			go tail(in)
		}
		select {}
	}
}
