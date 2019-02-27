package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("scraped. enter new line:")
		b, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("failed to read input; err= %q\n", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Add("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)

		w.Write(b)
		w.Write([]byte("\n"))

		fmt.Printf("wrote: %q\n", string(b))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("cannot run server; err= %q\n")
	}
}
