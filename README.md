# drone-amqp

## install vendor packages

Use https://github.com/govend/govend

Get dependencies running: 
```sh
govend -v --prune
```

## build

```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO15VENDOREXPERIMENT=1 go build -ldflags '-s -w' -o drone-amqp
```

```sh
docker build --rm -t local/drone-amqp .
```

## Run as drone plugin

Use PLUGIN_FILTER env variable to allow running of `local/drone-amqp` docker image
```
PLUGIN_FILTER="plugins/* local/*"
```


## testing
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