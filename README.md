# WallBoard - Ephemeral - No Signup - GPLv3 - Message Board

NOTE: This program is still in alpha stages and is extremely unstable. Don't 
expect it to be bug free. 

TagMachine is a new type of social media website that aims to accept and 
embrace social media as a modern form of journalism. 

TagMachine requires Redis (6.0.16+) and has only been tested on Linux servers.

To run this program: 

* clone the repository and `cd` into the project directory
* Start redis (generally `redis-server &`)
* run this command with your personal environment variables for the `hmac` sample
secret and testing password:
</a>

    hmacss=YOUR_SECRET_PHRASE testPass=YOUR_TESTING_PASS go run .

This will start TagMachine but the website will fail to load until you add 
data. You can add test data using another progam I'm creating as a test suite 
called [TagBot](https://github.com/hartsfield/TagBot).

# New Style

![Screenshot from 2023-08-26 19-23-38](https://github.com/hartsfield/machineTag/assets/30379836/6fa734ad-2dfb-4387-8f24-d8386acec19c)


# Requirements
## For Development
  - Go programming environment
  - Redis (v7+)
  - Only tested on Linux

# Instructions

Clone the repo, and run the following (assuming redis is running on the default port):

    go mod tidy
    chmod +x autoload.sh
    ./autoload.sh tm2 9999 [random phrase here]

Now visit `http://localhost:9999` and add some posts:

[example.webm](https://github.com/hartsfield/WallBoard/assets/30379836/326f0e8f-607c-468d-a657-3b294094a340)
