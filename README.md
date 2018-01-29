# Go Discord Bot
This bot is a work in progress, and will probalby never finish being developed. The idea is that people can make a pull request and add any features they feel like. Create an issue, let's discuss it! Have an idea, but don't know what to do? Give it a try!

## How do I setup Go Discord Bot?
First you'll need [Go](https://golang.org/doc/install) installed. Next up, grab this package.
```
go get -v github.com/LetsLearnCommunity/godiscordbot ./..
```

Open the directory.
```
cd $GOPATH/src/github.com/LetsLearnCommunity/godiscordbot/cmd/client/
```

Run the application
```
go run main.go -token [discord bot token] -music-channel-id [channel id]
```

Too obtain the channel id, go to your settings, click on appearance and enable the developer ui. Once you've done that, right click on the channel you want the bot to play in.