package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type UserInputReader interface {
	Read(secret bool) (string, error)
}

type StdinInputReader struct{}

func NewStdinInputReader() *StdinInputReader {
	return &StdinInputReader{}
}

func (reader StdinInputReader) Read(secret bool) (string, error) {
	if secret {
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		response := string(bytePassword)
		return response, nil
	}

	stdin_reader := bufio.NewReader(os.Stdin)
	response, err := stdin_reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return response, nil
}

func AreYouSure(message string, reader UserInputReader) bool {
	for {
		res := Ask(fmt.Sprintf("%s [y/n]", message), false, false, reader)
		return res == "y"
	}
}

func Ask(s string, required bool, secret bool, reader UserInputReader) string {
	var response string
	var err error
	for {
		fmt.Printf("%s: ", s)

		response, err = reader.Read(secret)
		response = strings.TrimSpace(response)

		if err != nil {
			continue
		}
		if len(response) == 0 && required {
			continue
		}
		break
	}

	return response
}
