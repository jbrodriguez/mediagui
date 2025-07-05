# Deployment Guide

This guide covers how to deploy mediagui using Docker and GitHub Container Registry (ghcr.io).

## Table of Contents

1. [GitHub Setup](#github-setup)
2. [Linux Box Setup](#linux-box-setup)
3. [Environment Configuration](#environment-configuration)
4. [Security Notes](#security-notes)
5. [Troubleshooting](#troubleshooting)

## GitHub Setup

### 1. Repository Settings

The GitHub Actions workflow is already configured to build and push Docker images to GitHub Container Registry (ghcr.io) automatically.

### 2. Enable GitHub Packages

GitHub Container Registry is enabled by default for all repositories. No additional setup is required.

### 3. Deploy to GitHub

Simply push your code to the main branch:

```bash
git push origin main
```

The GitHub Actions workflow will automatically:
- Build the React UI
- Build the Go binary with proper version tagging
- Create a multi-stage Docker image
- Push the image to `ghcr.io/your-username/mediagui`

### 4. Make Package Public (Optional)

To make your Docker image publicly accessible:

1. Go to your repository on GitHub
2. Click on "Packages" in the right sidebar
3. Click on the mediagui package
4. Click "Package settings"
5. Scroll down to "Change visibility"
6. Select "Public"

## Linux Box Setup

### 1. Install Docker & Docker Compose

For Ubuntu/Debian systems:

```bash
# Update package index
sudo apt update

# Install Docker
sudo apt install docker.io docker-compose-plugin

# Add your user to the docker group
sudo usermod -aG docker $USER

# Log out and log back in, or run:
newgrp docker

# Verify installation
docker --version
docker compose version
```

### 2. Authentication to GitHub Container Registry

You need to authenticate with GitHub Container Registry to pull images. Choose one of these methods:

#### Method A: Personal Access Token (Recommended)

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate a new token with `read:packages` scope
3. Log in to the registry:

```bash
echo "your_token_here" | docker login ghcr.io -u your-username --password-stdin
```

#### Method B: GitHub CLI

```bash
# Install GitHub CLI
sudo apt install gh

# Login to GitHub
gh auth login

# Configure Docker to use GitHub CLI
gh auth setup-git
```

### 3. Deploy on Linux Box

#### Create deployment directory

```bash
mkdir -p ~/mediagui
cd ~/mediagui
```

#### Download deployment files

```bash
# Download docker-compose.yml
curl -o docker-compose.yml https://raw.githubusercontent.com/your-username/mediagui/main/docker-compose.yml

# Download environment template
curl -o .env.example https://raw.githubusercontent.com/your-username/mediagui/main/.env.example
```

#### Configure environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit the environment file
nano .env
```

#### Deploy the service

```bash
# Pull the latest image
docker compose pull

# Start the service
docker compose up -d

# Check status
docker compose ps
docker compose logs -f mediagui
```

### 4. Updates and Maintenance

#### Update to latest version

```bash
# Pull latest image
docker compose pull

# Restart with new image
docker compose up -d

# Clean up old images
docker image prune -f
```

#### View logs

```bash
# View live logs
docker compose logs -f mediagui

# View recent logs
docker compose logs --tail 100 mediagui
```

#### Stop/Start service

```bash
# Stop
docker compose stop

# Start
docker compose start

# Restart
docker compose restart
```

## Environment Configuration

### Required Environment Variables

Create a `.env` file in your deployment directory with these variables:

```bash
# Media folders to scan (comma-separated paths)
MEDIA_FOLDERS=/media/movies,/media/tv,/media/documentaries

# Unraid hosts to monitor (comma-separated)
UNRAID_HOSTS=unraid.local:8080,backup-server.local:8080

# TMDB API key for movie metadata (required)
TMDB_KEY=your_tmdb_api_key_here

# User agent for API requests
USER_AGENT=mediagui/1.0

# User and group IDs for file permissions
USER_ID=1000
GROUP_ID=1000
```

### Optional Environment Variables

```bash
# Custom data directory (default: Docker volume)
DATA_DIR=/path/to/data

# Log level (debug, info, warn, error)
LOG_LEVEL=info

# Server port (default: 7623)
PORT=7623
```

### Volume Mounting

The docker-compose.yml file includes volume mounts for:

- `mediagui_data:/data` - Application data (database, config)
- `./config:/app/config:ro` - Configuration files (optional)

To mount your media folders, add them to the docker-compose.yml:

```yaml
volumes:
  - mediagui_data:/data
  - /path/to/your/movies:/media/movies:ro
  - /path/to/your/tv:/media/tv:ro
  - /path/to/your/documentaries:/media/documentaries:ro
```

### Getting a TMDB API Key

1. Go to [The Movie Database (TMDB)](https://www.themoviedb.org/)
2. Create a free account
3. Go to Settings → API
4. Request an API key
5. Use the "API Key (v3 auth)" in your `.env` file

## Security Notes

### GitHub Token Permissions

When creating a Personal Access Token for GitHub Container Registry:

- **Required scopes**: `read:packages`
- **Optional scopes**: `write:packages` (only if you plan to push images)
- Store the token securely and never commit it to version control

### Environment Variables

- Never commit your `.env` file to version control
- Use strong, unique values for sensitive configuration
- Consider using Docker secrets for production deployments

### File Permissions

The Docker container runs as a non-root user (UID 1000 by default). Ensure your media folders are readable by this user:

```bash
# Check current permissions
ls -la /path/to/media

# Fix permissions if needed
sudo chmod -R 755 /path/to/media
sudo chown -R 1000:1000 /path/to/media
```

### Network Security

- mediagui runs on port 7623 by default
- Consider using a reverse proxy (nginx, Traefik) for SSL termination
- Restrict access using firewall rules if needed

## Troubleshooting

### Common Issues

#### Permission Denied Errors

```bash
# Check if user is in docker group
groups $USER

# Add user to docker group if missing
sudo usermod -aG docker $USER
newgrp docker
```

#### Container Won't Start

```bash
# Check logs for errors
docker compose logs mediagui

# Check if port is already in use
sudo netstat -tlnp | grep :7623

# Check if image was pulled correctly
docker images | grep mediagui
```

#### Authentication Issues

```bash
# Test GitHub registry authentication
docker pull ghcr.io/your-username/mediagui:latest

# Re-authenticate if needed
docker logout ghcr.io
echo "your_token" | docker login ghcr.io -u your-username --password-stdin
```

#### File Permission Issues

```bash
# Check container user
docker compose exec mediagui id

# Check volume permissions
docker compose exec mediagui ls -la /data

# Fix volume permissions
sudo chown -R 1000:1000 /path/to/data
```

#### Build Issues

If you're building locally and encounter issues:

```bash
# Clean build (remove cache)
docker compose build --no-cache

# Build with verbose output
docker compose build --progress=plain

# Check if UI built correctly
docker compose run --rm mediagui ls -la /app
```

### Log Analysis

Common log patterns to look for:

- **Startup errors**: Check the first few lines of logs
- **Permission errors**: Look for "permission denied" messages
- **Network errors**: Check for connection refused or timeout errors
- **Configuration errors**: Look for environment variable validation messages

### Getting Help

If you encounter issues:

1. Check the logs: `docker compose logs mediagui`
2. Verify your environment configuration
3. Check file permissions
4. Ensure all required environment variables are set
5. Test network connectivity to your media folders and Unraid hosts

### Performance Tips

- Use local storage for the database volume for better performance
- Consider using SSD storage for the application data
- Monitor disk space usage regularly
- Use appropriate logging levels (avoid debug in production)

## Production Deployment

For production deployments, consider:

1. **Reverse Proxy**: Use nginx or Traefik for SSL termination
2. **Monitoring**: Set up log aggregation and monitoring
3. **Backups**: Regular backups of the data volume
4. **Updates**: Automated update strategies
5. **Security**: Regular security updates and vulnerability scanning

### Example nginx Configuration

```nginx
server {
    listen 80;
    server_name mediagui.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name mediagui.example.com;

    ssl_certificate /etc/ssl/certs/mediagui.crt;
    ssl_certificate_key /etc/ssl/private/mediagui.key;

    location / {
        proxy_pass http://localhost:7623;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

This completes the deployment guide for mediagui. The setup process is straightforward and should get you up and running quickly.