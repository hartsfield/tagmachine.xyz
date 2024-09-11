![Untitled drawing (7)](https://github.com/hartsfield/machineTag/assets/30379836/87bfbd09-ed33-4584-8ca9-d1a7034ee9ab)
# MachineTag is a re-write of TagMachine

NOTE: This program is still in alpha stages and is extremely unstable. Don't 
expect it to be bug free. 

TagMachine is a new type of social media website that aims to accept and 
embrace social media as a modern form of journalism. 

 - The supposedly "stable" version can be seen at [the main site](https://tagmachine.xyz/).
 - Updates to this new version can be seen at [the beta site](http://beta.tagmachine.xyz/).
 - The old version can be seen at [the old site](http://old.tagmachine.xyz/).

## Features implemented so far:
 - User Authentication
 - Posts
 - Replies
 - Webm/mp4/jpeg/png/gif upload support
 - Ranked posts
 - Two sorted views (ranked/new)
 - Pagification

SLOC = 1601 as of this commit

## Features implemented in the old version that need to be ported:
 - Tagging
 - Filter by tag
 - Mentions
 - Follow user
 - User pages
     - following stream
     - suggestions

## Features coming soon:
 - User favorites
 - Maybe votes
 - Maybe mp3 support
 - Maybe sorting algorithms

# Development
## Requirements
  - Go programming environment
  - Redis (v7+)
  - Only tested on Linux (Debian & RedHat based)

## Instructions for running

Clone the repo, and run the following (assuming redis is running on the default port):

    go mod tidy
    chmod +x autoload.sh
    ./autoload.sh tm2 9999 [random phrase here]

Now visit `http://localhost:9999`, sign-up, and add some posts.

## Bolt Architecture

An explanation of what Bolt architecture is can be found in [the git repository for the Bolt software](https://github.com/hartsfield/bolt)

Basically the front end is built using components composed of regular 
HTML/css/JavaScript, but they're executed via Go html templating, so Go 
template directives can be used. 

 - These components are mostly self-contained in directories located in `internal/components`.
 - Once created, these components can be added to a `page` located in 
`internal/pages`, or they can be used as a sub-component in another component.
 - Routes are registered in `routes.go`
 - View/Model structs are located in `viewdata.go`

![Screenshot from 2023-08-26 19-23-38](https://github.com/hartsfield/machineTag/assets/30379836/6fa734ad-2dfb-4387-8f24-d8386acec19c)
