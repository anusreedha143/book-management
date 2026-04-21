# STAGE 1: Build the Vue app using Node.js - named 'build'
# Sets the working directory to /app inside the container
# Uses a lightweight Node image to build the app, then copies the built files into a smaller Nginx image for serving
FROM node:20-alpine as build
WORKDIR /app

# Copy package.json and install dependencies
# Copies package.json + package-lock.json then 'run npm install' to install Vue and other dependencies in the container
# In this step we only copy the package.json files first and run npm install to take advantage of Docker's caching mechanism. 
# This way, if we change our source code but not our dependencies, Docker can reuse the cached layer with the installed node_modules,
# making subsequent builds faster.
COPY package*.json ./
RUN npm install

# Copy the rest of your Vue code/project and build it
# 'npm run build' compiles the Vue app into static files (HTML, CSS, JS) that can be served by Nginx.
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
