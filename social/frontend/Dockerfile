FROM node:lts-alpine as front

WORKDIR /app

COPY package*.json .

RUN --mount=type=cache,target=/app/node_modules \
    npm install

COPY . .

RUN --mount=type=cache,target=/app/node_modules \
    npm run build

FROM moeryomenko/ngxjs:latest

COPY --from=front /app/public /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

ENTRYPOINT [ "nginx", "-g", "daemon off;" ]
