# kubefirst-watcher

This is part of the [kubefirst watcher operator](https://github.com/kubefirst/kubefirst-watcher-operator). 

This tool will help to detect a combination of status of the cluster to hold action to be executed only after given state. 


# Notes

Updateting images at the moment:
```bash 
export GIT_SHA=`git rev-parse --short HEAD`
docker build -f build/Dockerfile .  -t k1test:$GIT_SHA
docker image tag k1test:$GIT_SHA 6zar/k1test:$GIT_SHA
docker image tag k1test:latest 6zar/k1test:latest
docker image push 6zar/k1test:latest
docker image push 6zar/k1test:$GIT_SHA
echo $GIT_SHA
```

This images need to be migrated to kubefirst docker registry
