package server

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/DNahar74/my-redis/store"
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
			s.DEL(key)
		}
	}

	if len(content) > 0 {
		// _, err := os.Open("./memory.aof")
		// if err != nil {
		file, err := os.Create("./memory.dat")
		if err != nil {
			fmt.Println("Error creating file")
			return
		}
		// fmt.Println("File created")
		// return
		// }
		_, err = io.WriteString(file, content)
		if err != nil {
			return
		}
	}
}
