package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	bolt "go.etcd.io/bbolt"

	"github.com/x1nchen/portainer-cli/cache"
)

func main() {
	homePath, err := os.UserHomeDir()
	panicIfError(err)
	s, err := cache.NewBoltStore(path.Join(homePath, ".portainer-cli"), "https://example.com")
	defer func() {
		err = s.Close()
		panicIfError(err)
	}()
	panicIfError(err)

	containers, err := s.ContainerService.FindAllContainers()
	panicIfError(err)
	fmt.Println("--- all containers ---")
	for _, container := range containers {
		fmt.Println(container.Names[0], container.EndpointName)
	}

	containers, err = s.ContainerService.FuzzyFindContainerByName("mt4cli")
	panicIfError(err)

	fmt.Println("--- fuzzy find containers ---")
	for _, container := range containers {
		fmt.Println(container.Names[0], container.EndpointName)
	}

	err = s.ContainerService.DB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("container"))

		_ = b.ForEach(func(k, v []byte) error {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			fmt.Printf("key=%s\n", k)
			return nil
		})
		return nil
	})

	fmt.Println("--- fuzzy ---")
	err = s.ContainerService.DB().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("container"))
		cursor := b.Cursor()
		// errcursor.Seek([]byte("mt4cli"))
		match := []byte("mt4cli")

		for k, _ := cursor.Seek(match); k != nil && bytes.Contains(k, match); k, _ = cursor.Next() {
			if bytes.Contains(k, match) {
				fmt.Println(string(k))
			}
		}

		return nil
	})

	panicIfError(err)

}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
