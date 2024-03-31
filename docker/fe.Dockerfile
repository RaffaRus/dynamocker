# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD the UI
##

# specify the image used to build the UI
FROM node:20-alpine3.18 as BUILD-UI

# create a working directory inside the image
WORKDIR /app

# copy the source of the project
COPY fe-app/ /app/

# install required packages
RUN npm install

# compile application
RUN npm run build --configuration=production

# ##
# ## STEP 2 - DEPLOY
# ##
FROM nginx:1.21-alpine

WORKDIR /

#  copy js project
COPY --from=BUILD-UI /app/dist/ /usr/share/nginx/html

# copy nginx config
COPY docker/nginx-dynamocker.conf /etc/nginx/conf.d/default.conf

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8151
