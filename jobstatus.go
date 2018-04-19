package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

//
// The compiled version of this program is to be placed on the swarm manager in /home/ubuntu/scripts
//
// it must then  be placed in cron to run after all jobs should have completed.
// for example: 21 6 * * * /home/ubuntu/scripts/jobstatus >/dev/null 2>&1

var htmlFile *os.File
var err error
var cronjobs map[string]int
var bytesWritten int
var regel string
var logToCheck string

func writestatus(cronjobs map[string]int) {
	htmlFile, err := os.Create("/home/ubuntu/jobstatus/metrics/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	// Write map to buffer
	for key, value := range cronjobs {
		len, err := htmlFile.WriteString(fmt.Sprintf("%s %d\n", key, value))
		if err != nil {
			log.Fatal(err)
			fmt.Printf("\nLength: %d bytes", len)
		}
	}
}

func miniGrep(path string, searchtext string) (int, error) {
	returncode := 0
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	re := regexp.MustCompile(searchtext)
	for _, regel := range lines {
		match := re.FindString(regel)
		if len(match) > 0 {
			returncode = returncode + 1
		}
	}
	return returncode, scanner.Err()
}

func main() {
	cronjobs := make(map[string]int)
	cronjobs["etl_status"] = 0
	cronjobs["backupswarm_status"] = 0
	cronjobs["backupportainer_status"] = 0
	os.MkdirAll("/home/ubuntu/jobstatus/metrics", 0755)

	//	   ### ETL job

	// Check if etl-log is up-to-date
	var (
		fileInfo os.FileInfo
		err      error
	)
	logToCheck = "/home/ubuntu/etl.log"
	fileInfo, err = os.Stat(logToCheck)
	if err != nil {
		log.Fatal(err)
	}
	modTime := fileInfo.ModTime()
	curTime := time.Now()
	diff := curTime.Sub(modTime)
	var matches int
	if diff > (time.Duration(720) * time.Minute) {
		fmt.Printf("%s is older than 12 hours (%s)\n", logToCheck, diff)
	} else {
		fmt.Printf("%s is newer than 12 hours (%s)\n", logToCheck, diff)
		// it is recent enough; let's check for expected lines
		matches, err = miniGrep(logToCheck, "Estimated remaining time: 0 m.*recommendations")
		if err != nil {
			log.Fatal(err)
		}
		if matches > 0 {
			log.Printf("%d matches found!\n", matches)
			// looking good so far. Let's check for Python traceback messages too
			matches, err = miniGrep(logToCheck, "traceback")
			if err != nil {
				log.Fatal(err)
			}
			if matches > 0 {
				log.Printf("%d traceback messages found. Not good.\n", matches)
			} else {
				cronjobs["etl_status"] = 1
				log.Printf("No traceback messages found. Looks okay.\n")
			}
		}
	}

	// Check portainer backup
	logToCheck = "/tmp/backupportainer.log"
	fileInfo, err = os.Stat(logToCheck)
	if err != nil {
		log.Fatal(err)
	}
	modTime = fileInfo.ModTime()
	curTime = time.Now()
	diff = curTime.Sub(modTime)
	if diff > (time.Duration(720) * time.Minute) {
		fmt.Printf("%s is older than 12 hours (%s)\n", logToCheck, diff)
	} else {
		fmt.Printf("%s is newer than 12 hours (%s)\n", logToCheck, diff)
		// it is recent enough; let's check for expected lines
		matches, err = miniGrep(logToCheck, "Backup succeeded")
		if err != nil {
			log.Fatal(err)
		}
		if matches > 0 {
			log.Printf("%d matches found!\n", matches)
			cronjobs["backupportainer_status"] = 1
			log.Printf("Portainer backup looks okay.\n")
		}
	}

	// Check swarm backup

	logToCheck = "/tmp/backupswarm.log"
	fileInfo, err = os.Stat(logToCheck)
	if err != nil {
		log.Fatal(err)
	}
	modTime = fileInfo.ModTime()
	curTime = time.Now()
	diff = curTime.Sub(modTime)
	if diff > (time.Duration(720) * time.Minute) {
		fmt.Printf("%s is older than 12 hours (%s)\n", logToCheck, diff)
	} else {
		fmt.Printf("%s is newer than 12 hours (%s)\n", logToCheck, diff)
		// it is recent enough; let's check for expected lines
		matches, err = miniGrep(logToCheck, "Backup succeeded")
		if err != nil {
			log.Fatal(err)
		}
		if matches > 0 {
			log.Printf("%d matches found!\n", matches)
			cronjobs["backupswarm_status"] = 1
			log.Printf("Swarm backup looks okay.\n")
		}
	}
	// Save the findings
	writestatus(cronjobs)
}
