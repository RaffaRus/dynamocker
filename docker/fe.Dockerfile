# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD the UI
##

# specify the image used to build the UI
FROM node:20-alpine3.18 AS build-ui

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

ENV FE_PORT=8151

# Used to redirect the requests toward the beckend
ENV BE_PORT=8150

WORKDIR /

#  copy js project
COPY --from=build-ui /app/dist/ /usr/share/

# copy nginx config
COPY docker/nginx-dynamocker.conf /etc/nginx/conf.d/default.conf

# tells Docker that the container listens on specified network ports at runtime
EXPOSE ${FE_PORT}

# replace the BE PORT
RUN echo 'find . -type f -name "main.*.js" -print0 | xargs -0 sed -i "s/localhost:8150/localhost:$BE_PORT/g"' > /dynamocker-entrypoint.sh

CMD ["/bin/sh", "-c", "sh dynamocker-entrypoint.sh && nginx -g 'daemon off;'"]
