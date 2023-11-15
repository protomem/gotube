FROM node:lts as deps

WORKDIR /app
COPY package.json package-lock.json ./

RUN npm install --frozen-lockfile



FROM node:lts as builder

WORKDIR /app

COPY . .
COPY --from=deps /app/node_modules ./node_modules

RUN npm run build



FROM node:lts as runner

WORKDIR /app
ENV NODE_ENV development

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

CMD ["npm", "run", "start"]
