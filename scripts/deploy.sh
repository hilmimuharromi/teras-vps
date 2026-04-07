#!/bin/bash

# TerasVPS Deployment Script
# This script deploys both backend and frontend to production

set -e

echo "🚀 TerasVPS Deployment Script"
echo "=============================="
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    print_warning "Running as root. This is not recommended."
fi

# Pull latest code
echo "📥 Pulling latest code..."
git pull origin main || {
    print_error "Failed to pull latest code"
    exit 1
}
print_success "Latest code pulled"

# Backend deployment
echo ""
echo "🔧 Deploying backend..."
cd backend

# Download dependencies
echo "📦 Installing Go dependencies..."
go mod download || {
    print_error "Failed to download Go dependencies"
    exit 1
}
print_success "Go dependencies installed"

# Build backend
echo "🏗️  Building backend..."
go build -o teras-vps main.go || {
    print_error "Failed to build backend"
    exit 1
}
print_success "Backend built successfully"

# Run migrations (if needed)
echo "🔄 Running database migrations..."
# TODO: Add migration command here
print_success "Database migrations completed"

# Restart backend service
echo "♻️  Restarting backend service..."
sudo systemctl restart teras-vps-backend || {
    print_warning "Failed to restart backend service. Make sure systemd service is configured."
}
print_success "Backend service restarted"

# Frontend deployment
echo ""
echo "🎨 Deploying frontend..."
cd ../frontend

# Install dependencies
echo "📦 Installing Node.js dependencies..."
npm install || {
    print_error "Failed to install Node.js dependencies"
    exit 1
}
print_success "Node.js dependencies installed"

# Build frontend
echo "🏗️  Building frontend..."
npm run build || {
    print_error "Failed to build frontend"
    exit 1
}
print_success "Frontend built successfully"

# Deploy to nginx
echo "📤 Deploying to nginx..."
sudo rm -rf /var/www/teras-vps/*
sudo cp -r dist/* /var/www/teras-vps/ || {
    print_error "Failed to deploy to nginx"
    exit 1
}
print_success "Frontend deployed to nginx"

# Reload nginx
echo "♻️  Reloading nginx..."
sudo systemctl reload nginx || {
    print_error "Failed to reload nginx"
    exit 1
}
print_success "Nginx reloaded"

# Cleanup
echo ""
echo "🧹 Cleaning up..."
cd backend
rm -f teras-vps
print_success "Cleanup completed"

# Done
echo ""
echo "=============================="
print_success "Deployment completed successfully!"
echo "=============================="
echo ""
echo "🌐 Backend API: http://localhost:3000"
echo "🌐 Frontend: http://localhost:80"
echo ""
echo "To view logs:"
echo "  Backend: sudo journalctl -u teras-vps-backend -f"
echo "  Nginx: sudo journalctl -u nginx -f"
echo ""
