# Stage 1: build static assets
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# Empty VITE_*_SERVICE_URL makes API calls relative (routed through nginx proxy).
ARG VITE_VEHICLE_SERVICE_URL=""
ENV VITE_VEHICLE_SERVICE_URL=$VITE_VEHICLE_SERVICE_URL
ARG VITE_ALERT_SERVICE_URL=""
ENV VITE_ALERT_SERVICE_URL=$VITE_ALERT_SERVICE_URL
RUN npm run build

# Stage 2: serve with nginx
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
