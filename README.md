# Neofy
Spotify Controller CLI App

# Table Of Contents
- [Why I Built This](#why-did-i-build-this)
- [Run Mock](#run-mock)
- [Requirements](#requirements)
- [Limitations](#limitations)
- [Usage](#usage)
- [Installation](#installation)
- [Future additions](#future-addtions)
- [Bugs](#bugs)
- [Contribution](#contribution)

# Why did I build this
I'm a developer that enjoys having music in the background.
I do development work mostly on my Mac, and I have Spotify
running in the background. I also have YouTube videos queued
in the background to help me or for leisure. When I click the
previous, pause/play, or next button on my keyboard, it seems
like my Mac selects at random to change my Spotify or YouTube.
My solution was Neofy, a Spotify CLI app that allows me to
control my Spotify.
Why a CLI app? I usually develop with tmux so I can quickly
switch to the app without interrupting my flow.
Why not just do ...? I had a problem that I wanted to
solve by coding. Not the most optimal solution, but the most fun for me.

# Run Mock
If you don't have Spotify Premium, don't want to make a Spotify app,
or just want to check out this project without the Spotify music 
controller. You can run a mock version of the app with the command below:
```bash
go run main.go -t
```

# Requirements
In order to run this CLI app properly, you will need 3 things:
* Spotify Premium
* Spotify App
* Linux based terminal
`NOTE` You will need to set up an application & copy the client ID & client secret.
For steps on how to do this, follow the Spotify [guide](https://developer.spotify.com/documentation/web-api/concepts/apps)

# Limitations
Neofy is built using the Spotify web API. The API
doesn't handle streaming, so you will need a Spotify client open.

# Installation
Currently, this app is under development, so ther is no current installation.
You can run the app as a Golang program.
With Spotify web credentials, you will need to add the following:
`.env` file in the root directory:
```
SPOTIFY_CLIENT_ID=<YOUR_CLIENT_ID>
SPOTIFY_CLIENT_SECRET=<YOUR_CLIENT_SECRET>
```
Once this has been added, you can just run:
```bash
# Note: In order to run the app Spotify must be playing a track inside a playlist
go run main.go
```
When you run the app it will redirect you to confirm access to Spotify on `localhost:8090`.
Once you accept this, you can return to the CLI.

# Usage
The CLI has 3 different modes: Player, Playlists, and Tracks.
`NOTE` The default mode is player
Player Key Binds:
* `<C-c>`: Exits app
* `u`: Switch to playlist mode
* `t`: Switch to track mode
* `s`: Toggles shuffle mode (on/off)
* `b`: Goes to previous song
* `p`: Plays song
* `x`: Pauses song
* `n`: Skips song
* `r`: Sets the next repeat mode (Flow: off -> context -> song)
* `-`: Lower the volume if applicable (Lowers by 10, range 0-100)
* `+`, `=`: Raises the volume if applicable (Raises by 10, range 0-100)
* `f`: Refreshes the current display data

Playlist Key binds:
* `<C-c>`, `<ESC>`: Switch to player mode
* `<C-u>`: Moves 10 rows up
* `<C-d>`: Moves 10 rows down
* `t`: Swtich to track mode
* `j`: Move Down
* `k`: Move Up
* `s`: Select Playlist

Tracks Key Binds:
* `<C-c>`, `<ESC>`: Switch to player mode
* `<C-u>`: Moves 10 rows up
* `<C-d>`: Moves 10 rows down
* `u`: Switch to playlist mode
* `j`: Move Down
* `k`: Move Up
* `s`: Play track

# Future additions
* Add auto-syncing when the track ends
* Add Syncing to Spotify (Tracks & playlists)
* Add support for non-tracks (podcasts)
* Customizable inputs
* Customizable window sizes
* Add Skimming for a track
* Add support for windows
* Add support to pick user devices

# Bugs
* Fix non-alphanumeric characters displaying width 2

# Contribution
This is a personal project. I'll update the app to suit
my needs. If there are any issues or suggestions, open a issue.

# Credit
* Inspiration for this project came from this [repo](#https://github.com/Rigellute/spotify-tui)
