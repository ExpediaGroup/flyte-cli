[![Build Status](https://travis-ci.org/HotelsDotCom/flyte-cli.svg?branch=master)](https://travis-ci.org/HotelsDotCom/flyte-cli)
# flyte/cli

`flyte/cli` is a command line client for flyte

## Make it
Build:
```
$ make build
```
Build & Install:
```
$ make install
```

## Use it
This is good place to start:
```
flyte [command]
```
The commands are:
```
help        Help about any command
test        Test step execution
version     Show the flyte version information
```

### Test command
Executes the step in the provided file. Test files MUST contain the
step, and trigger event definitions, and can optionally contain context and datastore
items as required. It should be in json or yaml format.

Example yaml file:
```
step:
  id: status
  event:
    packName: Slack
    name: ReceivedMessage
  criteria: "{{ Event.Payload.message|match:'^flyte status$' }}"
  context:
    UserID: "{{ Event.Payload.user.id }}"
  command:
    packName: Slack
    name: SendMessage
    input:
      channelId: "{{ Context.ChannelID }}"
      message: 'Hey <@{{ Context.UserID }}>, {{datastore(''message'')}}'
testData:
  event:
    pack:
      name: Slack
    event: ReceivedMessage
    payload:
      message: flyte status
      user:
        id: johnny
  context:
    ChannelID: '123'
  datastore:
    message: 'I''m up and running :run:'
```  
## Abuse it
Feel free to experiment and extend it by contributing back :relaxed:
