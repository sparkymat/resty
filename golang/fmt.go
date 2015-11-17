package golang

import (
	"io"
	"os"
	"os/exec"
)

func Fmt(inPath string, outPath string) error {
	args := []string{
		inPath,
	}
	cmd := exec.Command("gofmt", args...)
	cmdStdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer outFile.Close()

	io.Copy(outFile, cmdStdout)

	return nil
}
