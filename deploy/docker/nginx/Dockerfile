FROM node:24-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json ./
COPY web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

FROM nginx:1.27-alpine

COPY deploy/docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=web-builder /app/web/dist /usr/share/nginx/html
