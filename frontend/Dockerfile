FROM node:22.10.0-alpine AS frontend

COPY ./ ./

RUN npm install
RUN npm run build

EXPOSE 3000

CMD npm start