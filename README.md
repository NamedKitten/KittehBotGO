## KittehBotGo

This is the GO port of my discord bot.

### Redis

You need a redis server on the host.
Make sure that you set "prefix" to the bots prefix and "token" to the token.

### Building from source

#### Linux/BSD/Others

Make sure to install the required golang and git packages.

First clone the repo using the command:

```sh
git clone https://github.com/NamedKitten/KittehBotGO.git --recursive --recurse-submodules
```

Then change directory to the `KittehBotGO` dir.
Next make sure to run `export GOPATH=$PWD/lib`.
It makes sure you use the required included librarys.
Then you can run `go build src/main.cpp` and if it gives no error then move onto the Usage section.


### Usage

```
./main --help
Usage of ./main:
  -redisDB int
    	DB ID for redis server.
  -redisIP string
    	IP for redis server. (default "localhost")
  -redisPassword string
    	Password for redis server.
  -redisPort int
    	Port for redis server. (default 6379)
  -runSetup
    	Run setup?
```

Make sure to use the correct redis settings.
If it is your first time running it then make sure to run with `-runSetup` to set up the bot.
