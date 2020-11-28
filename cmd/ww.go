package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ysakasin/wildwest/pkg/buffer"
	"github.com/ysakasin/wildwest/pkg/disk"
	"github.com/ysakasin/wildwest/pkg/storage"
)

func main() {
	disk, err := disk.New("./ww.db")
	if err != nil {
		fmt.Println("Can't open db file")
		os.Exit(1)
	}

	buf := buffer.NewFifo(disk, 10)
	defer func() {
		buf.Flush()
		disk.Close()
	}()

	s := storage.New(disk, buf)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() {
		query := scanner.Text()
		if strings.HasPrefix(query, "exit") {
			break
		} else if strings.HasPrefix(query, "get ") {
			key := strings.TrimPrefix(query, "get ")
			str, ok := s.Get(key)
			if ok {
				fmt.Println(str)
			} else {
				fmt.Println("Not found:", key)
			}
		} else if strings.HasPrefix(query, "put ") {
			strs := strings.SplitN(query, " ", 3)
			key := strs[1]
			value := strs[2]
			s.Put(key, value)
		}

		fmt.Print("> ")
	}
}
