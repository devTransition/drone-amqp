# drone-amqp

## Using with drone:0.4

Image is available on Docker Hub

```
docker pull valichek/drone-amqp:0.4
```

It is not official plugin so use PLUGIN_FILTER env variable to allow it for Drone
```
PLUGIN_FILTER="plugins/* valichek/drone-amqp"
```

Add following to `.drone.yaml` to send amqp message when build finished

```yaml
notify:
  amqp:
    image: valichek/drone-amqp:0.4
    Connection: 
      Host: "192.168.99.100"
      Username: "guest"
      Password: "guest"
    Exchange: "rpc.test"
    Key: "routing.key.namespace"
    Publishing: 
      ContentType: "application/x-json"
      DeliveryMode: 1
    Template: >
      { "repo": "{{repo.full_name}}", "build": {{build.number}}, "commit": "{{build.commit}}"
```

Please, check [types.go](types.go) for available parameters.

## Dev notes

### Install vendor packages

Use https://github.com/govend/govend

Get dependencies: 
```sh
govend -v --prune
```

### Building image

Build binary:

```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO15VENDOREXPERIMENT=1 go build -ldflags '-s -w' -o drone-amqp
```

Build 
```sh
docker build --rm -t valichek/drone-amqp .
```

### Running

```sh
GO15VENDOREXPERIMENT=1 go run main.go types.go amqp.go < test.json
```

or

```sh
GO15VENDOREXPERIMENT=1 go run main.go types.go amqp.go <<EOF
{
  "repo": {
    "clone_url": "git://github.com/drone/drone",
    "owner": "drone",
    "name": "drone",
    "full_name": "drone/drone"
  },
  "system": {
    "link_url": "https://beta.drone.io"
  },
  "build": {
    "number": 22,
    "status": "success",
    "started_at": 1421029603,
    "finished_at": 1421029813,
    "message": "Update the Readme",
    "author": "johnsmith",
    "author_email": "john.smith@gmail.com",
    "event": "push",
    "branch": "master",
    "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
    "ref": "refs/heads/master"
  },
  "workspace": {
    "root": "/drone/src",
    "path": "/drone/src/github.com/drone/drone"
  },
  "vargs": {
    "Connection": {
      "Host": "192.168.99.100",
      "Username": "guest",
      "Password": "guest"
    },
    "Exchange": "some-exchange",
    "Key": "some-key",
    "Mandatory": true,
    "Publishing": {
      "Headers": {
        "header0": "header-zero",
        "header1": "header-one"
      },
      "ContentType": "application/x-json",
      "ContentEncoding": "",
      "DeliveryMode": 1,
      "Priority": 5,
      "CorrelationId": "Some-CorrelationId",
      "ReplyTo": "RPC",
      "Expiration": "Expiration-spec",
      "MessageId": "Some-MessageId",
      "Type": "Some-Type",
      "UserId": "Some-UserId",
      "AppId": "Some-AppId"
    },
    "Template": "{\"git_branch\": \"{{ build.branch }}\"}"
  }
}
EOF
```