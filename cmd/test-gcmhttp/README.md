# test FCM/HTTP

Command runs test of FCM/HTTP client based on [GCM client library](https://github.com/gamegos/gcmlib) from [Gamegos](https://github.com/gamegos).

## Usage

Configuration is following 12-factor app methodology and is using ENV as config source.

```bash
#!/usr/bin/env bash

export APP_LOG_LEVEL="all"
export APP_APPS_FILE="data/fcm-apps.json"
export APP_INSTANCES_FILE="data/fcm-inst.json"
export APP_GCM_ENDPOINT="http://localhost:8080"

./test-gcmhttp
```

## Example of test session

Test is done on OSX using docker image.

Launch server

    $ docker run --rm -i \
        -p 8080:8080 \
        -e "APP_SERVICE=fcm" \
        -e "APP_APPS_FILE=/srv/data/fcm-apps.json" \
        -e "APP_INSTANCES_FILE=/srv/data/fcm-inst.json" \
        -e "APP_LOG_LEVEL=all" \
        -v /Mac/go/src/github.com/szpakas/fakepushprovider/tmp:/srv/data \
        --name fakepushprovider fakepushprovider
    
    {"msg":"Parsed config from env => {Service:fcm AppsFile:/srv/data/fcm-apps.json InstancesFile:/srv/data/fcm-inst.json HTTPHost:0.0.0.0 HTTPPort:8080 LogLevel:all APNSCertFile: APNSKeyFile:}","level":"debug","ts":1465243462160068356,"fields":{}}
    {"msg":"starting","level":"info","ts":1465243462160091975,"fields":{}}
    {"msg":"import: apps","level":"info","ts":1465243462164974066,"fields":{"service":"FCM"}}
    {"msg":"import: instances","level":"info","ts":1465243462166711422,"fields":{"service":"FCM"}}
    {"msg":"import:instances:report => {Succeeded:20 Failed:0 Failures:[]}","level":"debug","ts":1465243462167045352,"fields":{"service":"FCM"}}
    {"msg":"storage:report => map[apps:total:2 apps:id=appId-1:id:appId-1 apps:id=appId-1:apiKey:apiKey-1 apps:id=appId-1:senderId:senderId-1 apps:id=appId-2:id:appId-2 apps:id=appId-2:apiKey:apiKey-2 instances:total:count:20 apps:id=appId-2:senderId:senderId-2 apps:id=appId-1:instances:total:10 apps:id=appId-2:instances:total:10]","level":"debug","ts":1465243462167097765,"fields":{"service":"FCM"}}
    {"msg":"start listening","level":"info","ts":1465243462167117342,"fields":{"service":"FCM","host":"0.0.0.0","port":8080}}

Launch test GCM client

    docker exec -it fakepushprovider \
        /bin/bash -c \
        "APP_APPS_FILE=/srv/data/fcm-apps.json APP_INSTANCES_FILE=/srv/data/fcm-inst.json APP_GCM_ENDPOINT=http://localhost:8080 /srv/test-gcmhttp"
    
    {"msg":"Parsed config from env => {AppsFile:/srv/data/fcm-apps.json InstancesFile:/srv/data/fcm-inst.json GCMEndpoint:http://localhost:8080 LogLevel:all}","level":"debug","ts":1465243469339848480,"fields":{}}
    {"msg":"starting","level":"info","ts":1465243469339915578,"fields":{}}
    {"msg":"import:apps","level":"info","ts":1465243469341939032,"fields":{}}

Client side message exchange

    {"msg":"message:sent","level":"debug","ts":1465243469344547628,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-1-regId-1","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469345190016,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-2-regId-2","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469345711995,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-3-regId-3","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469346363670,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-4-regId-8","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469346803693,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-5-regId-3","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469347392318,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-6-regId-2","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469348161941,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-7-regId-2","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469348955731,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-8-regId-6","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469349565041,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-9-regId-3","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469350064750,"fields":{"mobile:app:id":"appId-1","mobile:instance:token":"appId-1-instanceId-10-regId-1","response:success:count":0,"response:failure:count":1}}
    {"msg":"message:sent","level":"debug","ts":1465243469350548652,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-1-regId-4","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469351024983,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-2-regId-1","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469351560940,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-3-regId-4","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469352115510,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-4-regId-1","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469352577672,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-5-regId-1","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469353400966,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-6-regId-7","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469354110654,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-7-regId-3","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469354641980,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-8-regId-8","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469355166986,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-9-regId-2","response:success:count":1,"response:failure:count":0}}
    {"msg":"message:sent","level":"debug","ts":1465243469355720707,"fields":{"mobile:app:id":"appId-2","mobile:instance:token":"appId-2-instanceId-10-regId-3","response:success:count":0,"response:failure:count":1}}

Server side message exchange

    {"msg":"request:responded","level":"debug","ts":1465243469344220347,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469344825626,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469345439479,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469345932407,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469346574581,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469347160902,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469347771037,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469348666359,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469349311579,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469349905867,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":0,\"failure\":1,\"canonical_ids\":0,\"results\":[{\"error\":\"DEVICE_UNREGISTERED\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469350393529,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469350855454,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469351390457,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469351891063,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469352403690,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469353229338,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469353864725,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469354502222,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469354954706,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"m:1234\"}]}\n"}}
    {"msg":"request:responded","level":"debug","ts":1465243469355483328,"fields":{"service":"FCM","res:status":200,"res:body":"{\"multicast_id\":1234567890,\"success\":0,\"failure\":1,\"canonical_ids\":0,\"results\":[{\"error\":\"DEVICE_UNREGISTERED\"}]}\n"}}
