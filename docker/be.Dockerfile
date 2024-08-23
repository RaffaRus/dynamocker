# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD the BE
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.22-alpine AS build-be

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY be-app/go.mod /app/

# download Go modules and dependencies
RUN go mod tidy

# copy directory files i.e all files ending with .go
COPY be-app/ /app

RUN ls -ltu

# compile application
RUN go build -o build/dynamocker cmd/main.go

# test application
RUN go test ./... -v -p 1 

RUN mkdir -p dynamocker/bin
RUN mkdir -p dynamocker/mocks

RUN mv build/dynamocker dynamocker/bin/


##
## STEP 2 - DEPLOY
##
FROM scratch

ENV BE_PORT=8150

WORKDIR /

#  copy binary
COPY --from=build-be /app/dynamocker .

# tells Docker that the container listens on specified network ports at runtime
EXPOSE ${BE_PORT}

# command to be used to execute when the image is used to start a container
CMD [ "/bin/dynamocker" ]