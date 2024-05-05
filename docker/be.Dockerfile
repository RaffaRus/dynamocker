# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD the BE
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.17-alpine as BUILD-BE

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY be-app/go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY be-app/ ./

RUN ls -ltu

# compile application
RUN make build && make test

##
## STEP 2 - DEPLOY
##
FROM scratch

WORKDIR /

#  copy binary
COPY --from=BUILD-BE /app/dynamocker /dynamocker

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8150

# command to be used to execute when the image is used to start a container
CMD [ "/dynamocker" ]