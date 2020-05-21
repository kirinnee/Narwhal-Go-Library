# Narwhal Go Lib

Go version of Narwhal library, which is a library to allow you to programmatically perform certain docker 
actions.

Includes:
- Saving docker volume as tar ball
- Loading docker volume from tar ball
- Kill All Running Containers
- Stopping All Running Containers
- Remove all Containers
- Remove Images 
- Filtering Images
- Deploying via docker swarm 

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

stack: "stack-name" # you can have an alternative stack-name, which in the method, you can pass in "" to ignore

services:
  webapp:
    image: rocket:latest
    ports:
      - 8000:8000
    deploy:
      replicas: 3

# images, you can add as many as you need
images:

  # This builds at 'rocket' folder and uses 'rocket/df' as the Dockerfile 
  rocket:latest: 
    context: rocket 
    file: df 

  # This builds at 'golang' folder and uses 'golang/Dockerfile' as the Dockerfile
  golang:latest: 
    context: golang

  # This builds at '.' folder and uses './dfile' as the Dockerfile
  node:latest: 
    file: dfile
  
  # This builds at 'dotnet' folder and uses 'dotnet/Dockerfile' as the Dockerfile
  dotnet:latest: dotnet 
  
  # This builds at '.' folder and uses './Dockerfile' as the Dockerfile
  ror:latest:
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

### Build and Run Workspace 
`Run(<build context>, <dockerfile>, <image name>, <container name>)`
To ask the daemon to decide container name, please use `"""` (empty string) for 
`container name` field.
```go
err := narwhal.Run(".","Dockerfile","server:5", "")
if len(err) > 0 {
    //handle error here
}
``` 

### Get the list of Images 
`Images(<filter.....>)`
Get a list of images with the applied filter

returns: `images`, `remaining`, `errors`  
`images` - The list of image object after filtering
`remaining` - The remaining input passed in that is not filters, in order.
`errors` - Errors that occured

For the list of possible filters, look at [filtering](#Filtering)

```go
images, remain, err := narwhal.Images("ref=narwhal/*","dangling=true","label=5","before=6d")
if len(err) > 0 {
    //handle error here
}
```

### Remove Images
`RemoveImage(<filter.....>)`
Remove the images that pass the applied filters

For the list of possible filters, look at [filtering](#Filtering)
```go
 err := narwhal.RemoveImage("ref=narwhal/*","dangling=true","label=5","before=6d")
if len(err) > 0 {
    //handle error here
}
```

# Filtering

There are a few filter methods, all of which comes in the `<key>=<value>` format.

Filter works like the && operator, where images must fulfill **ALL** of the filters
applied

## Reference
`ref=` or `reference=` 

You can use this filter multiple times

This can filter images based on their image name and tag. Wildcards apply:

Images:
```
narwhal/a:0
narwhal/a:1
narwhal/b:0
narwhal/b:1
docker:latest
```
Examples:
```
ref=narwhal/* => narwhal/a:0,  narwhal/a:1,  narwhal/b:0,  narwhal/b:1,    
ref=docker* => docker
ref=narwhal/a:* => narwhal/a:0, narwhal/a:1
ref=narwhal/*:0 => narwhal/a:0, narwhal/b:0
```

## Label
`label=`

You can only use this filter once

This can filter images based on the present of the label, or the label with the exact 
value

Examples:
```
label=a  => all image with the label with key 'a'
label=b=5 => all images with the label 'b=5'
```

## Dangling
`dangling=`
You can use this filter once

This filter where the image is dangling

Examples:
```
dangling=true
```

## Time
`before=` and `after=` and `from= to=`

You can use this filter multiple times

Filter does a time query

### Type
There are a few types of time queries

#### Before
It will pick images that is before the time specified

Examples:
```yaml
# before the day 2020/02/05
before=2020/02/05 

# before 24h ago
before=24h  

# before the time of the specified image (alpine:latest) is created
before=alpine:latest 
```

#### After
It will pick images that is after the time specified

Examples:
```yaml
# after the day 2020/02/05
after=2020/02/05 

# everything is the last 24h 
after=24h 

# after the time of the specified image (alpine:latest) is created
after=alpine:latest 
```

#### From-To
It will pick images between the two times specified
Examples:
```yaml
# from the day 2020/02/05 to 24h ago 
"from=2020/02/05 to=24h"

# from 24 hour 30min 20s ago too the time where the image 'alpine latest' is created
"from=24h30m20s to=alpine:latest 
```
### Time Format
Time can come in many form, relative, absolute or based on images

#### Relative time
You can specify time relative to the current time. 

|short hand | reference |
|---|---|
|`h`|hour|
|`m`|minute|
|`s`|second|
|`ms`|millisecond|
|`us`| microsecond|
|`ns`|nanosecond|

Examples
```yaml
# Before 2 hours 20s ago
before=2h20s

# The past 2 min 50 microsecond and 60 nano seconds
after=2m50us60ns

# From 2days ago to yesterday
from=48h to=24h
```

#### Image-based Time
This allows you to use the reference the time of the image is 
created at

```yaml
# The image 
# server:0.1 created at 1 June 2020
# server:1.0 created at 10 June 2020

# before 1st June 2020
before=server:0.1

# after 10th June 2020
after=server:1.0
```

#### Absolute time format
There are multiple formats of this time, but do note if you 
do not provide the timezone, it will try to use your system's 
time zone. 

Format of dates are in `mm/dd`, never the other way round. 

Accepted formats:

```yaml
May 8, 2009 5:57:51 PM
oct 7, 1970
oct 7, '70
oct. 7, 1970
oct. 7, 70
Mon Jan  2 15:04:05 2006
Mon Jan  2 15:04:05 MST 2006
Mon Jan 02 15:04:05 -0700 2006
Monday, 02-Jan-06 15:04:05 MST
Mon, 02 Jan 2006 15:04:05 MST
Tue, 11 Jul 2017 16:28:13 +0200 (CEST)
Mon, 02 Jan 2006 15:04:05 -0700
Thu, 4 Jan 2018 17:53:36 +0000
Mon Aug 10 15:44:11 UTC+0100 2015
Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)
September 17, 2012 10:09am
September 17, 2012 at 10:09am PST-08
September 17, 2012, 10:10:09
October 7, 1970
October 7th, 1970
12 Feb 2006, 19:17
12 Feb 2006 19:17
7 oct 70
7 oct 1970
03 February 2013
1 July 2013
2013-Feb-03
 
#  mm/dd/yy
3/31/2014
03/31/2014
08/21/71
8/1/71
4/8/2014 22:05
04/08/2014 22:05
4/8/14 22:05
04/2/2014 03:00:51
8/8/1965 12:00:00 AM
8/8/1965 01:00:01 PM
8/8/1965 01:00 PM
8/8/1965 1:00 PM
8/8/1965 12:00 AM
4/02/2014 03:00:51
03/19/2012 10:11:59
03/19/2012 10:11:59.3186369

# yyyy/mm/dd
2014/3/31
2014/03/31
2014/4/8 22:05
2014/04/08 22:05
2014/04/2 03:00:51
2014/4/02 03:00:51
2012/03/19 10:11:59
2012/03/19 10:11:59.3186369

# Chinese
2014年04月08日

#   yyyy-mm-ddThh
2006-01-02T15:04:05+0000
2009-08-12T22:15:09-07:00
2009-08-12T22:15:09
2009-08-12T22:15:09Z

#  yyyy-mm-dd hh:mm:ss
2014-04-26 17:24:37.3186369
2012-08-03 18:31:59.257000000
2014-04-26 17:24:37.123
2013-04-01 22:43
2013-04-01 22:43:22
2014-12-16 06:20:00 UTC
2014-12-16 06:20:00 GMT
2014-04-26 05:24:37 PM
2014-04-26 13:13:43 +0800
2014-04-26 13:13:43 +0800 +08
2014-04-26 13:13:44 +09:00
2012-08-03 18:31:59.257000000 +0000 UTC
2015-09-30 18:48:56.35272715 +0000 UTC
2015-02-18 00:12:00 +0000 GMT
2015-02-18 00:12:00 +0000 UTC
2015-02-08 03:02:00 +0300 MSK m=+0.000000001
2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001
2017-07-19 03:21:51+00:00
2014-04-26
2014-04
2014
2014-05-11 08:20:13,787

# mm.dd.yy
3.31.2014
03.31.2014
08.21.71
2014.03
2014.03.30

# yyyymmdd and similar
20140601
20140722105203
# unix seconds, ms, micro, nano
1332151919
1384216367189
1384216367111222
1384216367111222333
```


## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.MD) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning 
We use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository.etting

## Authors
* [kirinnee](mailto:kirinnee97@gmail.com) 

## License
This project is licensed under MIT - see the [LICENSE.md](LICENSE.MD) file for details 