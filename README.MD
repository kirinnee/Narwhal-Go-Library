# Narwhal Go Lib

Go version of Narwhal library, which is a library to allow you to programmatically save and load docker 
volumes as tarballs


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
    narwhal := Narwhal{false}

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

## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.MD) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning 
We use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository.etting

## Authors
* [kirinnee](mailto:kirinnee97@gmail.com) 

## License
This project is licensed under MIT - see the [LICENSE.md](LICENSE.MD) file for details