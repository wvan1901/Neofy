# Neofy
Spotify Controller Cli App

# Table Of Contents
- [Why I Built This](#why-did-i-build-this)
- [Usage](#usage)
- [Installation](#installation)
- [Future additions](#future-addtions)
- [Bugs](#bugs)
- [Contribution](#contribution)

# Why did I build this
I'm a developer that enjoys having music in the background.
I do development work mostly on my mac and I have spotify
running in the background, I also have youtube videos queued
in the background to help me or for lesuire. When I click the
previous, pause/play or next button on my keyboard it seems
like my mac selects at random to change my spotify or youtube.
My solutions was Neofy, a spotify cli app that allows me to
control my spotify.
Why a cli app? I usually develop with tmux so I can quickly
switch to the app.
Why not just do ... etc? I had a problem which I wnated to
solve by coding, not the most optimal solution by the most fun for me.

# Installation
TODO: Finish writing section\
To run this project you will need to have spotifys developer access.
You will need to set up a application & copy the client id & client secret.
With this you will need to add the following .env file in the root directory:
```
SPOTIFY_CLIENT_ID=<YOUR_CLIENT_ID>
SPOTIFY_CLIENT_SECRET=<YOUR_CLIENT_SECRET>
```
Once this has been added you can just run
```bash
go run main.go
```
When your run the app it will redirect you to confirm access to spoity,
once you accept this you can return to the cli.

# Usage
The Cli has 3 diffrent modes: Player, Playlists, Tracks
Player Keybinds:
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

Playlist Keybinds:
* `<C-c>`, `<ESC>`: Switch to player mode
* `t`: Swtich to track mode
* `j`: Move Down
* `k`: Move Up
* `s`: Select Playlist

Tracks Keybinds:
* `<C-c>`, `<ESC>`: Switch to player mode
* `u`: Swtich to playlist mode
* `j`: Move Down
* `k`: Move Up
* `s`: Play track

# Future additions
* Improve playlist & tracks navigation
* Add autosyncing when track ends
* Add Syncing to Spotify (Tracks & playlists)
* Customizable inputs
* Customizable window sizes
* Add Skimming for a track
* Add support for windows
* Add support to pick user devices

# Bugs
* Fix non alphanumeric characters displaying width 2

# Contribution
This is a personal project & I'll update the app to suit
my needs. If there are any issues or suggestions open a issue.

