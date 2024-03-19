# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.17-alpine as BUILD

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./

# compile application
RUN go build -o /dynamocker

##
## STEP 2 - DEPLOY
##
FROM scratch

WORKDIR /

COPY --from=BUILD /dynamocker /dynamocker

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8888

# command to be used to execute when the image is used to start a container
CMD [ "/dynamocker" ]