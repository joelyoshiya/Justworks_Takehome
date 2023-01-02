# alpine chosen for small footprint
FROM golang:1.19.4-alpine

# set working directory
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

# copy go files, readme
COPY * .
# copy input and output directories
COPY input/ input/
COPY output/ output/

# compile the app
RUN go build -o Justworks_TakeHome

# use an entrypoint to run the app
ENTRYPOINT ["./Justworks_TakeHome"]

