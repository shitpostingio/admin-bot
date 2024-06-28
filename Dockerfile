# 🇧🇷 powered

# Pull from our golang image with tdlib installed
FROM registry.gitlab.com/shitposting/golang:latest as builder

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'adminbot:x:65534:65534:adminbot:/:' > /user/passwd && \
    echo 'adminbot:x:65534:' > /user/group

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/gitlab.com/shitposting/admin-bot

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Compile adminbot
RUN make install

# Execution stage
FROM registry.gitlab.com/shitposting/tdlib:latest

# Dependencies
RUN apt update && apt install -y -qq \
    gperf 

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Set the workdir
WORKDIR /home/adminbot

# Copy the built file
COPY --from=builder /go/bin/admin-bot .

# Run the executable
CMD ["./admin-bot", "-config", "configs/admin-bot.toml"]
