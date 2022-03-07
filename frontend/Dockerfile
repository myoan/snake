FROM node:lts-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY app.ts /app
COPY src /app/src
COPY tsconfig.json /app
COPY webpack.config.js /app
RUN npm run build

FROM nginx:stable-alpine 
COPY --from=builder /app/dist /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]