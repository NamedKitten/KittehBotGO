## KittehBotGo

This is the GO port of my discord bot.

### Redis

You need a redis server on the host.
Make sure that you set "prefix" to the bots prefix and "token" to the token.

### Installing from binary packages.

We provide binary packages for these following systems:

- Linux
    - amd64
    - i386
    - arm5
    - arm6
    - arm7
    - arm64
    - mips
    - mipsle
    - mips64
    - mips64le
    - ppc64
    - ppc64le
- OpenBSD
    - amd64
    - i386
    - arm5
    - arm6
    - arm7
- FreeBSD
    - amd64
    - i386
    - arm5
    - arm6
    - arm7
- NetBSD
    - amd64
    - i386
    - arm5
    - arm6
    - arm7
- DragonFly BSD
    - amd64

The releases are located [here](https://github.com/NamedKitten/KittehBotGo/releases).

### Installing from source.

#### Linux/BSD/Others

Make sure to install golang v1.8 or above.

Then create a `gopath` dir.    

Next make sure to run `export GOPATH=$PWD/gopath`.

Then run `go get -u github.com/golang/dep/cmd/dep` to download, install and compile dep to allow you to download KittehBotGO's required libraries.

Then run `go get -u github.com/jteeuwen/go-bindata` to download go-bindata.
go-bindata is used to store the assets for the built in dashboard inside the binary so that you don't neeed a folder full of assets wherever the bots executable is.

After that run `go install github.com/jteeuwen/go-bindata/go-bindata` to install and compile go-bindata.

Next run `go get -u github.com/NamedKitten/KittehBotGo` to download the bot's source.

Then navagate to `$GOPATH/src/github.com/NamedKitten/KittehBotGo` and run `$GOPATH/bin/dep ensure` to download the needed dependencies.

Next run `go generate` as this generates the file which the bot reads which contains all the assets.

Then you can navagate to `$GOPATH/bin` and run `go install github.com/NamedKitten/KittehBotGo`.

If you got no errors and there is a `KittehBotGo` file in the current directory then you have successfully compiled the bot.


### Usage

```
./KittehBotGo -h
Usage of ./KittehBotGo:
  -redisDB int
    	DB ID for redis server.
  -redisIP string
    	IP for redis server. (default "localhost")
  -redisPassword string
    	Password for redis server.
  -redisPort int
    	Port for redis server. (default 6379)
  -runDashboard
    	Run dashboard? (default true)
  -runSetup
    	Run setup?
  -version
    	Print version and exit.

```

Make sure to use the correct redis settings.
If it is your first time running it then make sure to run with `-runSetup` to set up the bot.
