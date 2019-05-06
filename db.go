package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func dbHandler() {

	// read OS env variables
	dbFilePath := "/tmp/restq.db"
	if len(os.Getenv("RESTQ_DB_FILE_PATH")) > 0 {
		dbFilePath = os.Getenv("RESTQ_DB_FILE_PATH")
	}

	dbFileUpdateInterval := 10
	if len(os.Getenv("RESTQ_DB_FILE_UPDATE_INTERVAL")) > 0 {
		dbFileUpdateInterval, _ = strconv.Atoi(os.Getenv("RESTQ_DB_FILE_UPDATE_INTERVAL"))
	}

	// dbFileUpdateInterval if set to 0 - do not save data on disk
	if dbFileUpdateInterval > 0 {

		// read current file if exists and size > 0
		jdata, err := ioutil.ReadFile(dbFilePath)
		if err == nil {
			err2 := json.Unmarshal(jdata, &QueueList)
			if err2 != nil {
				fmt.Println(err2)
			}
		} else {
			fmt.Println(err)
		}

		// file save loop
		for {
			time.Sleep(time.Second * time.Duration(dbFileUpdateInterval))
			mutex.Lock()
			jdata1, err := json.Marshal(QueueList)
			mutex.Unlock()
			if err == nil {
				err2 := ioutil.WriteFile(dbFilePath, jdata1, 0600)
				if err2 != nil {
					fmt.Println(err2)
				}
			} else {
				fmt.Println(err)
			}
		}

	}
}
