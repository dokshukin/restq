package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	mutex = &sync.Mutex{}
)

type request struct {
	Action  string  `json:"action"`
	Message message `json:"message"`
}

type response struct {
	Body    string `json:"body"`
	UUID    string `json:"uuid"`
	Created int64  `json:"created"`
	TTL     int64  `json:"ttl"`
}

func postQueueHandler(c *gin.Context) {
	queue := c.Param("queue")

	var req request
	err := c.BindJSON(&req)
	if err == nil {
		switch req.Action {

		case "pull":
			res, err := messagePull(queue, req)
			if err == nil {
				// send 200 response
				c.JSON(http.StatusOK, res)
			} else {
				// send 204 response
				c.Status(http.StatusNoContent)
			}

		case "push":
			// backgroud job to push message into queue
			go messagePush(queue, req)
			// send 201 response
			c.Status(http.StatusCreated)

		case "ack":
			err := messageAck(queue, req)
			if err == nil {
				// send 200 response
				c.Status(http.StatusOK)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "unsupported request",
				})
			}

		case "ext":
			err := messageExtend(queue, req)
			if err == nil {
				// send 200 response
				c.Status(http.StatusOK)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "unsupported request",
				})
			}

		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "unsupported action",
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported request",
		})
	}
}

func messagePull(queue string, req request) (res response, err error) {
	if len(QueueList[queue]) == 0 {
		err = errors.New("queue doesn't exist")
		return
	}

	flagMsgExists := false
	for key, msg := range QueueList[queue] {
		if msg.Status == statusOpen {
			mutex.Lock()
			res = response{
				Body:    msg.Body,
				UUID:    msg.UUID,
				Created: msg.Created.Unix(),
				TTL:     req.Message.TTL,
			}
			t := time.Now()
			if req.Message.TTL == 0 {
				QueueList[queue][key].Status = statusClosed
			} else {
				// todo: make a transaction
				QueueList[queue][key].Status = statusLocked
				QueueList[queue][key].Modified = t
				QueueList[queue][key].Expires = t.Add(
					time.Duration(req.Message.TTL) * time.Second)
				QueueList[queue][key].TTL = req.Message.TTL
			}
			mutex.Unlock()
			flagMsgExists = true
			break
		}
	}
	if flagMsgExists == false {
		err = errors.New("no messsage in the queue")
		return
	}
	return
}

func messagePush(queue string, req request) {
	// assemble message
	u, _ := genUUID()
	t := time.Now()
	msg := message{
		Body:     req.Message.Body,
		Created:  t,
		Modified: t,
		Expires:  t,
		TTL:      req.Message.TTL,
		Status:   statusOpen,
		UUID:     u,
	}

	// push msg to queue
	mutex.Lock()
	QueueList[queue] = append(QueueList[queue], msg)
	mutex.Unlock()
}

func messageExtend(queue string, req request) (err error) {
	if len(req.Message.UUID) > 0 {
		flagMsgExists := false
		for k, v := range QueueList[queue] {
			if v.UUID == req.Message.UUID {
				mutex.Lock()
				if req.Message.TTL > 0 {
					QueueList[queue][k].TTL = req.Message.TTL
				}
				QueueList[queue][k].Expires = time.Now().Add(
					time.Duration(QueueList[queue][k].TTL) * time.Second)
				mutex.Unlock()
				flagMsgExists = true
				break
			}
		}
		if flagMsgExists != true {
			err = errors.New("no message with uuid in the queue")
		}
	} else {
		err = errors.New("Unknown message uuid")
	}
	return
}

func messageAck(queue string, req request) (err error) {
	if len(req.Message.UUID) > 0 {
		flagMsgExists := false
		for k, v := range QueueList[queue] {
			if v.UUID == req.Message.UUID {
				mutex.Lock()
				QueueList[queue][k].Status = statusClosed
				mutex.Unlock()
				flagMsgExists = true
				break
			}
		}
		if flagMsgExists != true {
			err = errors.New("no message with uuid in the queue")
		}
	} else {
		err = errors.New("Unknown message uuid")
	}
	return
}

func queueCleaner() {
	// TODO move it to config
	msgExpireDays := defaultMsgExpireDays
	if len(os.Getenv(envMessageExpireDays)) > 0 {
		msgExpireDays, _ = strconv.Atoi(
			os.Getenv(envMessageExpireDays))
	}
	garbageCleanerInterval := defaultGarbageCleanerInterval // sec
	if len(os.Getenv(envGarbageCleanerInterval)) > 0 {
		garbageCleanerInterval, _ = strconv.Atoi(
			os.Getenv(envGarbageCleanerInterval))
	}
	msgAckExpireSec := defaultLockTimeoutSec // sec
	if len(os.Getenv(envMessageLockTimeoutSec)) > 0 {
		msgAckExpireSec, _ = strconv.Atoi(os.Getenv(envMessageLockTimeoutSec))
	}

	for {
		time.Sleep(time.Duration(garbageCleanerInterval) * time.Second)
		for queue := range QueueList {
			for key := 0; key < len(QueueList[queue]); key++ {
				if QueueList[queue][key].Status == statusClosed {
					// remove closed messages
					removeFromQueue(queue, key)
					key++
				} else if msgExpireDays > 0 &&
					(QueueList[queue][key].Expires.Unix()+int64(msgExpireDays*60*60*24) < time.Now().Unix()) {
					// remove expired messages
					removeFromQueue(queue, key)
					key++
				} else if QueueList[queue][key].Status == statusLocked &&
					QueueList[queue][key].Expires.Unix() < time.Now().Unix() {
					// unlock non ACK message after timeout
					mutex.Lock()
					QueueList[queue][key].Status = statusOpen
					QueueList[queue][key].Expires = time.Now().
						Add((time.Duration(msgAckExpireSec)))
					mutex.Unlock()
				}
			}
		}
	}
}

func genUUID() (uuid string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func removeFromQueue(queue string, key int) {
	mutex.Lock()
	QueueList[queue] = append(QueueList[queue][:key],
		QueueList[queue][key+1:]...)
	mutex.Unlock()
	return
}
