FROM node:lts-alpine

ENV NODE_ENV production
WORKDIR /usr/src/app
COPY . .
RUN npm ci --only=production
EXPOSE 3000
CMD ["npm", "start"]
