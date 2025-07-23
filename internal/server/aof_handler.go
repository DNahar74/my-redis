package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DNahar74/my-redis/internal/command"
	"github.com/DNahar74/my-redis/internal/resp"
	"github.com/DNahar74/my-redis/internal/store"
)

func handleAOF(s *store.Store) {
	file, err := os.OpenFile("./commands.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file")
	}
	for {
		time.Sleep(1 * time.Second)
		writeToFile(file, s)
	}
}

func writeToFile(file *os.File, s *store.Store) {
	content := ""

	select {
	case cmd := <-s.AOFChan:
		content += cmd
	default:
		return
	}

	_, err := file.WriteString(content + "\n#\n")
	if err != nil {
		fmt.Println("Error writing commands to AOF file :: ", err)
	}
}

func restoreStorage() error {
	fileB, err := os.ReadFile("./commands.aof")
	if err != nil {
		if os.IsNotExist(err) {
			// file doesn't exist, silently skip
			return nil
		}
		return err
	}

	content := string(fileB)
	cmds := strings.Split(content, "\n#\n")

	for i, c := range cmds[:len(cmds)-1] {
		cmd, err := resp.Deserialize(c)
		if err != nil {
			fmt.Println("Error in cmd", i)
			return err
		}
		_, err = command.HandleCommands(cmd)
		if err != nil {
			fmt.Println("Error in restoring storage. Cmd:", i)
			return err
		}
	}

	return nil
}
