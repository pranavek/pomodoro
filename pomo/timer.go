package pomo

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
)

func Run() {

	pomoCount := 0
	carryOn := true

	for carryOn == true {
		fmt.Println("Starting pomodoro timer (25 minutes)")

		time.Sleep(25 * time.Minute)
		fmt.Println("End of pomodoro interval")

		err := beeep.Alert("Pomodoro", "End of Pomodoro", "assets/information.png")
		if err != nil {
			panic(err)
		}

		err = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}

		pomoCount += 1

		fmt.Println("Check Marks:", pomoCount)

		if pomoCount == 4 {
			fmt.Println("Take a long break - 30 minutes")
			time.Sleep(30 * time.Minute)
			pomoCount = 0
		} else {
			fmt.Println("Take a short break - 5 minutes")
			time.Sleep(5 * time.Minute)
		}

		//Ask for input to set carryon as true or false
	}
	fmt.Println("Good bye!")
}
