# restq
Simple REST queue daemon

## Protocol

### pull
POST /URI_PREFIX/queue_name
```
{
  "action": "pull", 
  "message": {                                // optional
    "ttl": 0,                               // optional
}
```
`message` is not mandatory field. 
If `ttl` is not set, message is will be removed from queue!

### push
POST /URI_PREFIX/queue_name
```
{
  "action": "push",
  "message": {
    "body": "anything could be set here",   // mandatory field
    "ttl": 10,                              // optional, by default 10
}
```

### ack
POST /URI_PREFIX/queue_name
```
{
  "action": "ack",
  "message": {
    "uuid": "251d36de-6ce7-11e9-a923-1681be663d3e",   // mandatory field
}
```
`Acknowlege` message had been handled by consumer and should be removed from the queue.

### ext
POST /URI_PREFIX/queue_name
```
{
  "action": "ext",
  "message": {
    "uuid": "251d36de-6ce7-11e9-a923-1681be663d3e",   // mandatory field
    "ttl": 10                                         // mandatory field
}
```
`Extend` message locked time.


## Configuration
Environment variables:
* `RESTQ_BIND_IP=0.0.0.0` listen on specified IP address
* `RESTQ_BIND_PORT=8080` bind to port
* `RESTQ_PREFIX_URI=/` URI prefix
* `RESTQ_DB_FILE_PATH=/tmp/restq.db` flush current queue data to file
* `RESTQ_DB_FILE_UPDATE_INTERVAL=10` flush interval in seconds
* `RESTQ_GARBAGE_CLEANER_INTERVAL=10` in seconds, regular job to clean up all closed messages from the queue
* `RESTQ_MESSAGE_EXPIRE_DAYS=2` time-to-live for all not updated messages, if message is not handled during expiration interval, it will be removed from queue

## To do 

* generate UUIDs with go
* add GET requests
* work with some external DB (f.e. BoltDB)

## Internal message structure
```
{
    "action": "pull",                // pull, push, ack, ext
    "message": {
        "uuid": "251d36de-6ce7-11e9-a923-1681be663d3e", 
        "status": "open",            // open/locked/closed
        "created": epochtime,
        "modified": epochtime,
        "expires": epochtime,
        "ttl": int,                  // seconds
        "body": "a bla-bla-bla"      // text
    }
}
```


