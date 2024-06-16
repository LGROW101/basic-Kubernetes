# Use the official Node.js image as the base image
FROM node:18.17.0

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the Next.js application
RUN npm run build

# Expose the port on which the application will run
EXPOSE 3000

# Run the application
CMD ["npm", "start"]
