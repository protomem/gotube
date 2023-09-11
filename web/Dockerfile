FROM node:alpine AS deps

WORKDIR /app
COPY package.json yarn.lock .
RUN yarn install



FROM node:alpine AS builder

ARG VITE_API_URL

WORKDIR /app
COPY . .
COPY --from=deps /app/node_modules ./node_modules
RUN VITE_API_URL=${VITE_API_URL} yarn build



FROM nginx:alpine AS runner

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html

