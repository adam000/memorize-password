package main

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

const (
	minLength      = 4
	alwaysHashUpTo = 50
	paddingRune    = '\x00'
)

func main() {
	var password []byte
	var hashes incrementalHash
	// TODO if no password file exists...
	if true {
		// TODO get user password
		password = []byte("rosebud")
		hashesResult, err := makeIncrementalHash(password)
		hashes = hashesResult
		if err != nil {
			panic(err)
			// TODO
		}
		// TODO store the password
	} else {
		// TODO read existing password file (later: pass file as command line flag)
	}
	length := len(password)
	_ = length

	for index, partialHash := range hashes {
		// TODO prompt for password fragment (hint with length)
		fmt.Printf("Password? (%d) ", index+minLength)
		stdin := int(os.Stdin.Fd())
		passwordAttempt, err := terminal.ReadPassword(stdin)
		_ = passwordAttempt
		fmt.Println()
		if err != nil {
			panic(err)
			// TODO
		}
		// Compare password fragment to the hash we have
		if !partialHash.matches(passwordAttempt) {
			fmt.Println("\nTry again.")
			continue
		}
		// TODO use proper check
		if index+minLength == length {
			fmt.Println("You did it!")
			break
		}
	}
}

type partialHash struct {
	Hash []byte
	Salt []byte
}

func (h partialHash) matches(attempt []byte) bool {
	// TODO
	return true
}

func makePartialHash(password, salt []byte) (partialHash, error) {
	// TODO evaluate
	const (
		cost   = 16384
		r      = 1
		p      = 1
		keyLen = 32
	)

	hash, err := scrypt.Key(password, salt, cost, r, p, keyLen)
	return partialHash{hash, salt}, err
}

type incrementalHash []partialHash

// For a given password, hashes various lengths of it and returns a map of the length to its hash
func makeIncrementalHash(password []byte) (incrementalHash, error) {
	hashes := make(incrementalHash, 0, alwaysHashUpTo-minLength+1)

	paddedPassword := make([]byte, alwaysHashUpTo)
	copy(paddedPassword, password)

	for currentLength := minLength; currentLength <= alwaysHashUpTo; currentLength++ {
		// TODO make a salt
		salt := make([]byte, 32)
		hash, err := makePartialHash(paddedPassword[:currentLength], salt)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}

	return hashes, nil
}
