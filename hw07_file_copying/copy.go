package main

import (
	"errors"
	"io"
	"os"
)

var chunkSize int64 = 0x20000 // 128KB

func makeCopy(from io.Reader, to io.Writer, limit int64, progress chan int64) (err error) {
	chunk, remaining := chunkSize, limit
	var count int64
	for remaining > 0 && !errors.Is(err, io.EOF) {
		if remaining < chunk {
			chunk = remaining
		}

		count, err = io.CopyN(to, from, chunk)
		if err != nil && !errors.Is(err, io.EOF) {
			break
		}

		remaining -= count
		progress <- count
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}

	close(progress)
	return err
}

func postError(finish chan error, err error, postAnyway bool) bool {
	post := err != nil || postAnyway
	if post {
		finish <- err
		close(finish)
	}
	return post
}

func cp(params *CopyParams, progress chan int64, finish chan error) {
	from, err := os.Open(params.from)
	if postError(finish, err, false) {
		return
	}
	defer from.Close()
	from.Seek(params.offset, io.SeekStart)

	to, err := os.Create(params.to)
	if postError(finish, err, false) {
		return
	}
	defer to.Close()

	postError(finish, makeCopy(from, to, params.limit, progress), true)
}
