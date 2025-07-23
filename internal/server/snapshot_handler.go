package server

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/DNahar74/PulseDB/internal/store"
)

func handleMemoryState(s *store.Store) {
	for {
		time.Sleep(10 * time.Second)
		handleFile(s)
	}
}

func handleFile(s *store.Store) {
	content := ""

	for key, value := range s.Items {
		if value.Expiry.IsZero() || value.Expiry.After(time.Now()) {
			content += key + " => \t"
			str, err := value.Value.Serialize()
			if err != nil {
				return
			}
			content += str
		} else {
			err := s.DEL(key)
			if err != nil {
				fmt.Println("error deleting in background worker:", err)
				return
			}
		}
	}

	if len(content) > 0 {
		file, err := os.Create("./memory.dat")
		if err != nil {
			fmt.Println("Error creating file")
			return
		}
		_, err = io.WriteString(file, content)
		if err != nil {
			return
		}
	}
}
