package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
)

type CopyParams struct {
	from, to      string
	limit, offset int64
}

var params = &CopyParams{}

func init() {
	flag.StringVar(&params.from, "from", "", "file to read from")
	flag.StringVar(&params.to, "to", "", "file to write to")
	flag.Int64Var(&params.limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&params.offset, "offset", 0, "offset in input file")
}

func initialChecks(params *CopyParams) error {
	stat, _ := os.Stat(params.from)
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if params.offset < 0 {
		params.offset = 0
	}

	if stat.Size() <= params.offset {
		return ErrOffsetExceedsFileSize
	}

	if params.limit <= 0 || params.limit > stat.Size()-params.offset {
		params.limit = stat.Size() - params.offset
	}
	return nil
}

func main() {
	flag.Parse()

	if err := initialChecks(params); err != nil {
		fmt.Println(err)
		return
	}

	progress, finish := make(chan int64), make(chan error)

	go Copy(params, progress, finish)

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
