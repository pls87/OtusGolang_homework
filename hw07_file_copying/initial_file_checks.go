package main

import "os"

func initialFileChecks(params *CopyParams) error {
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
