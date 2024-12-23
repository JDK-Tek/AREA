FROM node

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

RUN npm run build
RUN npm install -g serve
CMD npm start