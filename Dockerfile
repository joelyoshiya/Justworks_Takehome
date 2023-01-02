# alpine chosen for small footprint
FROM golang:1.19.4-alpine

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY * .
COPY input/* ./input/
COPY output/* ./output/

# Need to set the user arguments as environment variables
ENV input_file_name=data_raw_1.csv
ENV output_file_name=output.csv

# compile the app
RUN go build -o Justworks_TakeHome

# run the app
CMD ["./Justworks_TakeHome", ${input_file_name}, ${output_file_name}]