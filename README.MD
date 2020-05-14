# Narwhal Go Lib

Go version of Narwhal library, which is a library to allow you to programmatically perform certain docker 
actions.

Includes:
- Saving docker volume as tar ball
- Loading docker volume from tar ball
- Kill All Running Containers

# Getting Started
Import it as a library (go module)
```go
package main

import(
    "gitlab.com/kiringo/narwhal_lib"
)
```

### Create a narwhal instance  
Argument is whether its quiet. If true, it will suppress all console logs.
```go
func main(){
    narwhal := narwhal_lib.New(false)

}
``` 

### Loading Tarball into a docker named volume
`Load(<volume name>, <path to tar>)`
```go
err := narwhal.Load("my-volume","./data.tar.gz")
if len(err) > 0 {
    //handle error here
}
```

### Saving docker volume as a Tarball
`Save(<volume name>, <tarball name file name>, <relative path to save to>)`
```go
err := narwhal.Save("my-volume", "data", "./")
if len(err) > 0 {
    //handle error here
}
```

### Killing all docker Containers
`KillAll()`
```go
err := narwhal.KillAll()
if len(err) > 0 {
    //handle error here
}
```

### Stopping all docker Containers
`StopAll()`
```go
err := narwhal.KillAll()
if len(err) > 0 {
    //handle error here
}
```
### Remove all docker Containers
`RemoveAll()`
```go
err := narwhal.RemoveAll()
if len(err) > 0 {
    //handle error here
}
```

### Deploy a stack with swarm
`Deploy(<stack name>, <compose file name>)`
This deploy option comes with additional feature of auto-building images. In the `docker-compose.yml`, you can 
opt to add a `image` category, which allows you to build images before running the stack, allowing you to have the
`build` feature like with older `docker-compose` versions 

```go
err := narwhal.Deploy("app", "docker-compose.yml")
if len(err) > 0 {
    //handle error here
}
```

**Example docker-compose.yml** :
```yaml
version: '3.8'

services:
  webapp:
    image: rocket:latest
    ports:
      - 8000:8000
    deploy:
      replicas: 3

images:
  # images, you can add as many as you need 
  rokcet:latest:
    context: rocket # this can be excluded if its the current folder
    file: Dockerfile # this can also be excluded if the dockerfile is "Dockerfile"
```

### Deploy a stack with swarm in auto mode
`Deploy(<stack name>, <compose file name>, <unsafe>)`
This mode will automatically initialize swarm if its is not in swarm mode.

The `unsafe` option will automatically leave the swarm and re-init the swarm before if the deployment failed
initially. (hard reset).

```go
err := narwhal.DeployAuto("app", "docker-compose.yml", false)
if len(err) > 0 {
    //handle error here
}
```

## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.MD) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning 
We use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository.etting

## Authors
* [kirinnee](mailto:kirinnee97@gmail.com) 

## License
This project is licensed under MIT - see the [LICENSE.md](LICENSE.MD) file for details