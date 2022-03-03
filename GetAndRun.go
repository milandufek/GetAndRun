package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

/*
Simple program.
Get URL.
Parse its body as strings separated by ';'.
Run command on your system.
Profit.
*/

const HttpTimeout = 20 * time.Second

func main() {

	// parse args
	var url = flag.String("u", "", "URL")
	var period = flag.Int("p", 0, "Repeat time period.")
	var maxRepeat = flag.Int("m", 1, "Maximum number of repeats.")
	var onBackGround = flag.Int("b", 1, "Run command at background.")
	flag.Parse()

	var reqCount = 1

	if *url == "" {
		log.Fatal("URL is empty")
	}

	// proceed request(s)
	for true {
		body := httpGet(*url)
		if body != nil {
			runCmd(string(body), *onBackGround)
		}

		if *period <= 1 || *maxRepeat == reqCount {
			break
		}

		time.Sleep(time.Duration(*period) * time.Second)

		if *maxRepeat != 0 {
			reqCount++
		}
	}
}

func httpGet(url string) []byte {

	client := http.Client{Timeout: HttpTimeout}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("HTTP error [%s]\n", err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed [%s]\n", resp.Status)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		return nil
	}

	return body
}

func runCmd(input string, onBackGround int) {

	commands := strings.Split(input, ";")

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)

		slices := strings.Split(cmd, " ")
		bin := slices[0]
		args := slices[1:]

		if bin == "" {
			return
		}

		if onBackGround == 1 {
			log.Printf("Running command [%s] at background", cmd)
			err := exec.Command(bin, args...).Start()
			if err != nil {
				log.Println("Command failed")
			}

		} else {
			log.Printf("Running command [%s]", cmd)
			out, err := exec.Command(bin, args...).Output()
			if err != nil {
				log.Println("Command failed")
			}
			fmt.Printf("%s\n", out)
		}
	}
}
