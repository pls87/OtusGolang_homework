package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/cheggaaa/pb/v3"
)

type CopyParams struct {
	from, to      string
	limit, offset int64
}

var (
	params                   = CopyParams{}
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func init() {
	flag.StringVar(&params.from, "from", "", "file to read from")
	flag.StringVar(&params.to, "to", "", "file to write to")
	flag.Int64Var(&params.limit, "limit", 0, "limit of bytes to cp")
	flag.Int64Var(&params.offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if err := initialFileChecks(&params); err != nil {
		fmt.Println(err)
		return
	}

	progress, finish := make(chan int64), make(chan error)

	go cp(params, progress, finish)

	bar := pb.Start64(params.limit)
	for {
		select {
		case status := <-finish:
			if status != nil {
				fmt.Println("Error occurred: ", status)
			}
			bar.Finish()
			return
		case delta := <-progress:
			bar.Add64(delta)
		}
	}
}
