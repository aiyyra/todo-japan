# # base image for the GO app
FROM golang:1.23

# # create a directory instances to be the working directory in docker
WORKDIR /usr/src/app

# # Copying the root from local directory to working dir in docker
# # copy go mod and sum file first to ensure no error
COPY go.mod go.sum ./
# # ensure go mod is working and verify all the dependencies/package
RUN go mod download && go mod verify

# # now we can copy all without any problem
COPY . .

# # build the app
RUN go build -v -o /usr/local/bin/app ./...

# # Specify Port
EXPOSE 8080

# # run the file with app cmd
CMD [ "app" ]
