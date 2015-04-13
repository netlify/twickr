# Twickr

> Send interesting tweets from Twitter to Slack.

Twickr connects to Twitter and listens for keywords you're interested in. Whenever a keyword is mentioned in a tweet, Twickr sends you a message over Slack.

## Installation

Download the relevant release for your server's OS and architecture and extract the files.

[Twickr Downloads](https://github.com/netlify/twickr/releases)

## Configuration

Before you can run Twickr, you need to setup a [Twitter application](https://apps.twitter.com/) and a [Slack Inbound Webhook](https://netlify.slack.com/services/new/incoming-webhook).

### Twitter Setup

1. Go to [Twitter's Application Management](https://apps.twitter.com/) and create a new app.
2. Fill out the name description and website and create the app.
3. Go to "Keys and Access Tokens" and create a new access token

### Slack Setup

1. Create a new Inbound Webhook for your Slack account [here](https://netlify.slack.com/services/new/incoming-webhook)

### The Configuration File

Create a `config.json` file in the same folder as the twickr executable. An example configuration file is included in the download.

1. Fill in the credentials for your new Twitter app and the URL from your Slack webhook.
2. Optionally specify which channel to send Twickr messages to
3. Enter a coma-separated list of keywords you want Twickr to keep an eye on

## Running

```bash
cd <path-to-your-twickr-install>
./twickr
```

Normally you'll run Twickr on a server (rather than on your laptop or desktop machine) where it can keep running continuously. An easy way to make sure it stays up after you close your SSH connection is to run it in a TMUX or Screen session.

## Netlify

Twickr is a small tool developed by [Netlify](https://www.netlify.com). Netlify builds, deploys & hosts static sites and apps.
