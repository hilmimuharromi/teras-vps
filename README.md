# TerasVPS - VPS Panel for Selling Hosting

**Your own VPS panel built with Go Fiber + Astro**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Astro](https://img.shields.io/badge/Astro-FF5D01?style=flat&logo=astro&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

---

## 📋 Overview

TerasVPS is a self-hosted VPS management panel that allows you to sell VPS instances from your mini PC (Proxmox-based) to customers.

### Features

- ✅ User registration & authentication (JWT)
- ✅ Create/Manage VMs (Start, Stop, Reboot)
- ✅ Real-time monitoring (CPU, RAM, Disk)
- ✅ Billing system with invoice generation
- ✅ Payment gateway integration (QRIS - coming soon)
- ✅ Automatic suspend/delete for overdue invoices
- ✅ Backup system (create, restore)
- ✅ SSH keys management
- ✅ Admin dashboard
- ✅ VNC console access

---

## 🏗️ Architecture

```
teras-vps/
├── backend/          # Go Fiber API Server
├── frontend/         # Astro Frontend
├── docker/           # Docker Compose
├── scripts/          # Deployment Scripts
└── docs/             # Documentation
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Proxmox VE 7+

### Development Setup

```bash
# Clone repository
git clone https://github.com/hilmimuharromi/teras-vps.git
cd teras-vps

# Backend setup
cd backend
go mod download
cp .env.example .env
# Edit .env with your credentials

# Frontend setup
cd ../frontend
npm install
npm run dev
```

---

## 🐳 Docker Setup

```bash
# Start all services
docker-compose up -d

# Run database migrations
docker-compose exec backend go run migrations/main.go
```

---

## 📊 Tech Stack

### Backend
- **Framework:** Go Fiber v2
- **Database:** PostgreSQL + GORM
- **Cache:** Redis
- **Auth:** JWT
- **Proxmox:** proxmox-api-go

### Frontend
- **Framework:** Astro
- **Styling:** TailwindCSS
- **Components:** Shadcn/ui
- **Charts:** Recharts
- **Icons:** Lucide React

---

## 📖 Documentation

- [Deployment Guide](docs/DEPLOYMENT.md)
- [API Documentation](docs/API.md)
- [User Manual](docs/USER_MANUAL.md)

---

## 🔧 Configuration

See `.env.example` for all available environment variables.

### Key Environment Variables

```bash
# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=teras_vps
POSTGRES_PASSWORD=your_password
POSTGRES_DB=teras_vps

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your_secret_key

# Proxmox
PROXMOX_HOST=https://your-proxmox-ip:8006/api2/json
PROXMOX_USER=root@pam
PROXMOX_PASSWORD=your_proxmox_password
```

---

## 🤝 Contributing

Contributions are welcome! Please open an issue or PR.

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 👨‍💻 Author

**Hilmi Muharromi** - [@hilmi_muharromi](https://twitter.com/hilmi_muharromi)

---

## 🙏 Acknowledgments

- Built with [Go Fiber](https://docs.gofiber.io)
- Frontend powered by [Astro](https://astro.build)
- UI components from [Shadcn/ui](https://ui.shadcn.com)
