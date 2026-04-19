# STAGE 1: Build the Vue app using Node.js
FROM node:20-alpine as build
WORKDIR /app

# Copy package.json and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of your Vue code and build it
COPY . .
RUN npm run build


# STAGE 2: Serve the compiled app with Nginx
FROM nginx:alpine

# Copy the Nginx configuration file you already created
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy the compiled static files from Stage 1 into Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 for Kubernetes to route traffic to
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
