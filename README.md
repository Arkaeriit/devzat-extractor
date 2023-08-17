# Devzat-Extractor

On [Devzat](https://github.com/quackduck/devzat), people send a lot of messages and can have meaningful conversations. If you want to save those conversations, this plugin serves them over a web interface.

## In-chat usage

This plugins listen to every messages it receives. If you want to extract the messages sent in the last 30 minutes in the current room, use the command `extract 30m`. The plugin will reply with an URL such as `https://devzat.bobignou.red/timespan/main/1692250655/1692261454/extract.txt`. On it, you will find all the messages you wanted.

## Web API

The web API offers two routes:

* `/timespan/<room>/<from>/<to>/extract.txt`: extracts the messages from the given room sent between the two UNIX timestamps `from` and `to`.
* `/timespan-all/<from>/<to>/extract.txt`: extracts the messages from every room sent between the two UNIX timestamps `from` and `to`.

It will reply a 200 code if some messages are recovered, 204 when no messages have been found, 400 when the timestamps are invalid, and 404 if any other route is accesses.

## Admin usage

The plugin is made for a single-file executable. It is configured with the following environment variable.

|Variable name |Description                                        |Default                                                                     |
|--------------|---------------------------------------------------|----------------------------------------------------------------------------|
|`DEVZAT_HOST` |URL of the chat-room interface                     |`https://devzat.hackclub.com:5556`                                          |
|`DEVZAT_TOKEN`|Authentication token                               |Does not defaults to anything. The program panics if the token is not given.|
|`PORT`        |The port used to serve the web API                 |8080                                                                        |
|`HOST`        |The host name and protocol of the web API          |`http://localhost:8080`                                                     |
|`BANK_SIZE`   |The maximum number of messages the plugin remembers|1000                                                                        |

