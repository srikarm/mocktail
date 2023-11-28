FROM node:alpine
RUN mkdir -p /code
WORKDIR /code
ADD mocktail-dashboard /code
RUN npm install -g  npm@10.2.4 --legacy-peer-deps
RUN npm install --save  react-scripts react-bootstrap bootstrap --legacy-peer-deps
RUN yarn run build && yarn cache clean
CMD [ "npm", "start" ]
EXPOSE 3001
