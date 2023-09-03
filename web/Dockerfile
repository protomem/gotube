FROM node:18-alpine AS deps

WORKDIR /app

COPY package.json ./
COPY package-lock.json ./
RUN npm ci



FROM node:18-alpine AS builder

ARG NEXT_PUBLIC_APP_URL

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN NEXT_PUBLIC_APP_URL=$NEXT_PUBLIC_APP_URL npm run build



FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

USER nextjs

EXPOSE 3000

CMD ["npm", "start"]

