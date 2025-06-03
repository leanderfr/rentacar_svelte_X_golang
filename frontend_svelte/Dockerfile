

FROM node:22.14.0 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build
RUN npm prune --production

FROM node:22.14.0
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .
EXPOSE 3000
ENV NODE_ENV=production
CMD [ "node", "build" ]


#  docker build --rm -t rentacar_frontend .

#  docker run   -p 3000:3000 --name rentacar_frontend rentacar_frontend

