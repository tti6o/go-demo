package main

import (
	"log"
	"os/exec"
	"time"
)

func TimingCopy() {
	delay := 1*time.Second
	t := time.NewTicker(delay)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//fmt.Println("every 1 seconds executing")
			cmd := exec.Command("xcopy","Y:\\Library\\Application Support\\ManicTime\\Screenshots",
				"C:\\Users\\nickyluo\\AppData\\Local\\Finkit\\ManicTime\\Screenshots","/D","/E")
			err := cmd.Run()
			if err != nil {
				log.Printf("xcopy err=%v",err)
			}
		}
	}
}

func main() {
	log.Println("Start!")
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Printf("Process stopped - %+v", r)
	//	} else {
	//		log.Printf("Process stopped")
	//	}
	//	TimingCopy()
	//}()
	//TimingCopy()
	log.Println("End!")
}