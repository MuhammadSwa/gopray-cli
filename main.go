package main

import (
	"fmt"
	"log"
	"myThing/gopray"
	"os"
)

func main() {
	a, err := gopray.MakeApp()
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) > 1 {

		switch os.Args[1] {
		case "list":
			a.ListAllPrayers()

		case "next":
			fmt.Println(a.TimeLeftForNextPrayer())

		case "date":
			fmt.Println(a.HijriDate())

		// case "help":
		// 	help()
		default:
			fmt.Fprintln(os.Stderr, "Please provide a valid argument")
			// help()
		}
		// } else {
		// 	help()
		// }
	}
}
