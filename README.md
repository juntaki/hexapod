# Hexapod software

Run it on Raspberry Pi Zero W for [Hexapod](https://www.thingiverse.com/thing:3769750), and control the robot from your laptop

# Build and Setup

```
$ make
cd servo; GOARM=6 GOOS=linux GOARCH=arm go build -o ../robot .
cd controller; go build -o ../ct .
```

* Save Pub/Sub credential on project root as "cred.json"
* Run `echo -n 0 > seq.txt`

# Install

Send all files to robot. Change systemd setting for startup.

TODO

# Use Controller

```
# Subscribe heartbeat from robot
./ct heartbeat 

# move commands
./ct rotate
./ct walk
./ct arms 0 0 0 0 0 0 0 0 # <odd arms up-down> <even arms up-down> <arm1 rotate> <arm2> <arm3> <arm4> <arm5> <arm6>
```
