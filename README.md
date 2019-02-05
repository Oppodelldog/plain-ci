![alt text](https://github.com/Oppodelldog/plain-ci/raw/master/webview/assets/images/logo.png)
> just another continuous integration server

So what are the minimalistic features of a ci server?

* clone git repository
* checkout a commit
* execute a build script (build/ci.sh)
* trigger some API before and after a build

So that's what this project is about.

### run
```go run cmd/main.go```

Per default it listens to port 12345 and has 2 built-in triggers.

### trigger build
a simple post is enough to trigger a build
```shell
curl -X POST \
  http://localhost:12345/hook/simple \
  -H 'content-type: application/json' \
  -d '{
  "URL": "https://github.com/yourGithubUserName/yourRepo.git",
  "REV": "fd1245758904390c4b474bb7842b473b6cbe92f0"
}'
```

### list queue
just the build queue
```shell
curl -X GET \
  http://localhost:12345/queue \
  -H 'accept: application/json'
```

### abort a build
```shell
curl -X GET \
  http://localhost:12345/queue/abort/171e21a7b6244253bf7c6e23d429e739
```

### github integration
Github integration allow to automatically start a build for each push to github.
When the build job finishes it sends a commit status notification to github indicating
the status of the build.  

create a token for plain-ci
* https://github.com/settings/tokens
when starting plain-ci, set the token in **env var** ```PLAIN_CI_GITHUB_TOKEN```

create a Webhook for your github project
+ http://yourIpAddress:12345/hook/github
+ Content-Type: application/json   
+ Secret does not matter right now :-)   
+ Only push notifications are supported   
