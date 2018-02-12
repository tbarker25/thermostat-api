package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func myDeepEquals(got, want interface{}) error {
	gotFile, err := ioutil.TempFile("", "got.json")
	if err != nil {
		panic(err)
	}
	defer gotFile.Close()
	defer os.Remove(gotFile.Name())

	wantFile, err := ioutil.TempFile("", "want.json")
	if err != nil {
		panic(err)
	}
	defer wantFile.Close()
	defer os.Remove(wantFile.Name())

	gotEncoder := json.NewEncoder(gotFile)
	gotEncoder.SetIndent("", "  ")
	err = gotEncoder.Encode(got)
	if err != nil {
		panic(err)
	}

	wantEncoder := json.NewEncoder(wantFile)
	wantEncoder.SetIndent("", "  ")
	err = wantEncoder.Encode(want)
	if err != nil {
		panic(err)
	}

	wantFile.Close()
	gotFile.Close()

	cmd := exec.Command("diff", gotFile.Name(), wantFile.Name())
	output, err := cmd.CombinedOutput()

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("got: <, want: >\n%s", output)
		}
		panic(err)
	}

	return nil
}

func myBytesEquals(got, want []byte) error {
	if !bytes.Equal(got, want) {
		return fmt.Errorf("got: %s\nwant: %s", hex.Dump(got), hex.Dump(want))
	}
	return nil
}
