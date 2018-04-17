package main

import (
	"fmt"
	"log"
	"os"
)

//
// This script is to be placed on the swarm manager in /home/ubuntu/scripts
//
// it must then  be placed in cron to run after all jobs should have completed.
// for example: 1 6 * * * /home/ubuntu/scripts/jobstatus.sh >/dev/null 2>&1

// setup-work in case it is required
// mkdir -p /home/ubuntu/jobstatus/metrics >/dev/null 2>&1

var htmlFile *os.File
var err error
var cronjobs map[string]int
var bytesWritten int
var regel string

func writestatus(cronjobs map[string]int) {
	htmlFile, err := os.Create("/home/ubuntu/jobstatus/metrics/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	// Write map to buffer
	for key, value := range cronjobs {
		fmt.Println("%s %i\n", key, value)
		len, err := htmlFile.WriteString(fmt.Sprintf("%s %d\n", key, value))
		if err != nil {
			log.Fatal(err)
			fmt.Printf("\nLength: %d bytes", len)
		}
	}

}

func main() {
	cronjobs := make(map[string]int)
	cronjobs["etl_status"] = 0
	cronjobs["backupswarm_status"] = 0
	cronjobs["backupportainer_status"] = 0
	os.MkdirAll("/home/ubuntu/jobstatus/metrics", 0755)

	writestatus(cronjobs)
	/*
	   ### ETL job

	   SUCCESS=0

	   # Check if etl-log is up-to-date

	   if [ `find "/home/ubuntu/etl.log" -mmin +720` ]
	   then
	       # file is more than 12 hours old and therefore not up-to-date
	       SUCCESS=0
	   else
	      tail -100 /home/ubuntu/etl.log | grep "Estimated remaining time: 0 m.*recommendations" >/dev/null 2>&1
	      if [ $? -eq 0 ]
	      then
	        # We found the line. Indicates successful end
	        # Let's double-check for errors
	        grep -i traceback /home/ubuntu/etl.log >/dev/null 2>&1
	        if [ $? -ne 1 ]
	        then
	          # We found a traceback message. indicating failure
	          SUCCESS=0
	        else
	          SUCCESS=1
	        fi
	      else
	        SUCCESS=0
	      fi
	   fi

	   # Write it to file to serve

	   if [ $SUCCESS -eq 1 ]
	   then
	     # All good
	     setstatus_ok etl_status
	   else
	     # Bad stuff has happened
	     setstatus_ko etl_status
	   fi


	   ### END ETL job

	   ### Swarm metadata backup
	   SUCCESS=0
	   # Check if backup-log is up-to-date

	   if [ `find "/tmp/backupswarm.log" -mmin +720` ]
	   then
	       # file is more than 12 hours old and therefore not up-to-date
	       SUCCESS=0
	   else
	      grep "Backup succeeded" /tmp/backupswarm.log >/dev/null 2>&1
	      if [ $? -eq 0 ]
	      then
	        # We found the line. Indicates successful end
	        SUCCESS=1
	      else
	        SUCCESS=0
	      fi
	   fi

	   # Write it to file to serve

	   if [ $SUCCESS -eq 1 ]
	   then
	     # All good
	     setstatus_ok backupswarm_status
	   else
	     # Bad stuff has happened
	     setstatus_ko backupswarm_status
	   fi

	   ### END swarm metadata backup


	   ### Portainer metadata backup
	   SUCCESS=0
	   # Check if backup-log is up-to-date

	   if [ `find "/tmp/backupportainer.log" -mmin +720` ]
	   then
	       # file is more than 12 hours old and therefore not up-to-date
	       SUCCESS=0
	   else
	      grep "Backup succeeded" /tmp/backupportainer.log >/dev/null 2>&1
	      if [ $? -eq 0 ]
	      then
	        # We found the line. Indicates successful end
	        SUCCESS=1
	      else
	        SUCCESS=0
	      fi
	   fi

	   # Write it to file to serve

	   if [ $SUCCESS -eq 1 ]
	   then
	     # All good
	     setstatus_ok backupportainer_status
	   else
	     # Bad stuff has happened
	     setstatus_ko backupportainer_status
	   fi

	   ### END portainer metadata backup
	*/
}
