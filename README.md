# History bot for slack
Receive and search slack history for all of your organization's channels

Built with [slack](https://github.com/nlopes/slack) and [bleve](http://blevesearch.com)


## Prerequisites

1. [Go](https://golang.org/)
2. [Make](https://www.gnu.org/software/make/)

## Installation

```
mkdir -p $GOPATH/src/github.com/maddevsio/
cd $GOPATH/src/github.com/maddevsio
git clone https://github.com/maddevsio/slack_history_bot
cd slack_history_bot
make depends
make
```

Or golang way

```
mkdir -p $GOPATH/src/github.com/maddevsio/
cd $GOPATH/src/github.com/maddevsio
git clone https://github.com/maddevsio/slack_history_bot
cd slack_history_bot
go get -v
go build -v
go install
```

## Configure

1. Obtain slack api key here https://api.slack.com/bot-users
2. export SLACK_TOKEN=YOUR_TOKEN

## Run

```
./slackhistory_bot
```
