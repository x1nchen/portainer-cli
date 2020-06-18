package main

import (
	"fmt"
	"os"
	"path"

	"github.com/x1nchen/portainer-cli/cache"
)

func main() {
	homePath, err := os.UserHomeDir()
	panicIfError(err)
	s, err := cache.NewBoltStore(path.Join(homePath, ".portainer-cli"), "https://portainer.followme-internal.com")
	defer func() {
		err = s.Close()
		panicIfError(err)
	}()
	panicIfError(err)

	token, err := s.GetToken()
	panicIfError(err)

	fmt.Println("token:", token)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
