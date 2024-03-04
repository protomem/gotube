FROM node:lts-alpine AS deps

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY ./package.json ./package-lock.json .
RUN npm install



FROM node:lts-alpine AS builder

WORKDIR /app

COPY . .
COPY --from=deps /app/node_modules ./node_modules

RUN npm run build



FROM nginx:stable-alpine AS runner

ENV NODE_ENV production

WORKDIR /app

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html



