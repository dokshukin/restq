package main

const (
	// constant statuses of messages in queue
	statusOpen   = 0
	statusLocked = 1
	statusClosed = 2

	// default message expire (0 - never)
	defaultMsgExpireDays = 0

	defaultLockTimeoutSec         = 3600 // sec = 1 hour
	defaultGarbageCleanerInterval = 10   // sec

	// ENV variables
	envMessageExpireDays      = "RESTQ_MESSAGE_EXPIRE_DAYS"
	envGarbageCleanerInterval = "RESTQ_GARBAGE_CLEANER_INTERVAL"
	envMessageLockTimeoutSec  = "RESTQ_MESSAGE_LOCK_TIMEOUT_SECONDS"
)
