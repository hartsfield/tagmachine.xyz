# MachineTag is a re-write of TagMachine

NOTE: This program is still in alpha stages and is extremely unstable. Don't 
expect it to be bug free. 

TagMachine is a new type of social media website that aims to accept and 
embrace social media as a modern form of journalism. 

To run this program: 

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

Now visit `http://localhost:9999`, sign-up, and add some posts.
