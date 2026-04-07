# 🚀 TerasVPS

VPS Hosting Indonesia dengan Harga Terjangkau - Platform VPS berbasis Proxmox dengan integrasi otomatis.

## 📋 Fitur Utama

- ✅ **VM Management**: Buat, kelola, dan monitor VPS dari panel
- ✅ **Real-time Monitoring**: Chart CPU, Memory, dan Disk secara live
- ✅ **Billing Otomatis**: Invoice bulanan dengan auto-suspend
- ✅ **Support Ticket**: Buka tiket support dan komunikasi langsung
- ✅ **Payment Integration**: Pembayaran online via Xendit
- ✅ **Admin Panel**: Kelola user, VM, invoice, dan support ticket
- ✅ **SSH Key Management**: Upload SSH key untuk akses VM

## 🏗️ Teknologi

### Backend
- **Go** (Golang) - Backend API
- **Fiber** - Web framework
- **PostgreSQL** - Database
- **Redis** - Cache & sessions
- **GORM** - ORM

### Frontend
- **Astro** - Static site generator
- **React** - UI components
- **TailwindCSS** - Styling
- **Shadcn/ui** - UI components

### Infrastructure
- **Proxmox** - Virtualization platform
- **Docker** - Containerization
- **Nginx** - Reverse proxy

## 📦 Instalasi

### Prerequisites
- Docker & Docker Compose
- PostgreSQL (atau via Docker)
- Redis (atau via Docker)
- Proxmox Server

### 1. Clone Repository

```bash
git clone https://github.com/hilmimuharromi/teras-vps.git
cd teras-vps
```

### 2. Setup Environment Variables

Copy file `.env.example` dan sesuaikan konfigurasi:

**Backend:**
```bash
cd backend
cp .env.example .env
nano .env
```

**Frontend:**
```bash
cd frontend
cp .env.example .env
nano .env
```

### 3. Update Konfigurasi Penting

Di `backend/.env`, pastikan update:

```env
# Database password
DB_PASSWORD=your_secure_password_here

# JWT secret (gunakan string panjang & random)
JWT_SECRET=your_very_long_random_secret_key_change_this_in_production

# Proxmox credentials
PROXMOX_HOST=https://your-proxmox-server.com:8006
PROXMOX_USER=root@pam
PROXMOX_PASSWORD=your_proxmox_password

# Xendit API key
XENDIT_SECRET_KEY=your_xendit_secret_key_here

# Admin default
ADMIN_EMAIL=admin@terasvps.id
ADMIN_PASSWORD=Admin123!
```

### 4. Jalankan dengan Docker

```bash
# Dari root directory project
docker-compose up -d

# Cek status
docker-compose ps

# Lihat logs
docker-compose logs -f backend
```

### 5. Setup Database

```bash
# Jalankan migration
docker-compose exec backend go run migrations/migrate.go

# Atau via Docker build (jika ada seeder)
docker-compose exec backend go run main.go seed
```

### 6. Akses Aplikasi

- **Frontend**: http://localhost:4321
- **Backend API**: http://localhost:3000
- **Admin Panel**: http://localhost:4321/admin

Login admin default:
- Email: `admin@terasvps.id`
- Password: `Admin123!`

## 📁 Struktur Project

```
teras-vps/
├── backend/
│   ├── controllers/      # API handlers
│   ├── cron/            # Scheduled jobs (billing)
│   ├── database/        # Database migrations
│   ├── middleware/      # Auth, Admin, Logging
│   ├── models/          # Database models
│   ├── proxmox/         # Proxmox integration
│   ├── routes/          # API routes
│   ├── services/        # Business logic
│   ├── utils/           # Helper functions
│   ├── websocket/       # WebSocket server
│   └── main.go          # Entry point
├── frontend/
│   ├── src/
│   │   ├── pages/       # Halaman aplikasi
│   │   ├── components/  # UI components
│   │   └── layouts/     # Layout components
│   └── public/          # Static assets
├── deploy/              # Deployment configs
│   ├── nginx.conf       # Nginx config
│   └── terasvps.service # Systemd service
├── docker-compose.yml   # Docker setup
└── README.md           # Documentation
```

## 🔧 Konfigurasi Penting

### Proxmox Settings

Pastikan Proxmox server di-setup dengan:
- API user dengan permission
- Storage untuk VM
- Network bridge (vmbr0)
- Resource pools untuk isolasi

### Billing Automation

Cron jobs otomatis berjalan:
- Invoice generation: Setiap tanggal 1
- Auto-suspend: 7 hari setelah due date
- Auto-delete: 14 hari setelah due date

### Payment Gateway

Saat ini mendukung:
- **Xendit** (utama)
- Stripe (opsional - perlu di-setup)

## 🚀 Deployment

### Deploy ke Mini PC / VPS

```bash
# Clone repository
git clone https://github.com/hilmimuharromi/teras-vps.git
cd teras-vps

# Setup environment
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
nano backend/.env

# Build dan start
docker-compose up -d

# Setup SSL dengan Let's Encrypt
certbot certonly --nginx -d terasvps.id
```

### Nginx Configuration

Copy `deploy/nginx.conf` ke `/etc/nginx/sites-available/terasvps`:

```bash
sudo cp deploy/nginx.conf /etc/nginx/sites-available/terasvps
sudo ln -s /etc/nginx/sites-available/terasvps /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Systemd Service

Untuk auto-start saat boot:

```bash
sudo cp deploy/terasvps.service /etc/systemd/system/
sudo systemctl enable terasvps
sudo systemctl start terasvps
```

## 📊 API Endpoints

### Authentication
- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login user
- `GET /api/auth/me` - Get current user

### VM Management
- `GET /api/vms` - List semua VM user
- `GET /api/vms/:id` - Detail VM
- `POST /api/vms` - Buat VM baru
- `POST /api/vms/:id/start` - Start VM
- `POST /api/vms/:id/stop` - Stop VM
- `POST /api/vms/:id/reboot` - Reboot VM
- `DELETE /api/vms/:id` - Delete VM

### Billing
- `GET /api/billing/invoices` - List invoice
- `GET /api/billing/invoices/:id` - Detail invoice
- `POST /api/billing/pay/:id` - Bayar invoice

### Support
- `GET /api/support/tickets` - List tickets
- `POST /api/support/tickets` - Buat ticket
- `POST /api/support/tickets/:id/messages` - Kirim message

### Admin
- `GET /api/admin/stats` - Platform statistics
- `GET /api/admin/users` - List semua user
- `POST /api/admin/users/:id/suspend` - Suspend user

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Run specific test
go test ./controllers/vm_controller_test.go

# Frontend tests
cd frontend
npm test
```

## 📝 Development

### Backend Development

```bash
cd backend
go mod download
go run main.go
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

### WebSocket Testing

```bash
# Connect ke WebSocket server
wscat -c ws://localhost:3000/ws/vm/1?token=your_jwt_token
```

## 🔒 Security

- ✅ JWT authentication
- ✅ Password hashing (bcrypt)
- ✅ SQL injection prevention (GORM)
- ✅ CORS configuration
- ✅ Rate limiting (Fiber middleware)
- ✅ HTTPS ready

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Check PostgreSQL status
docker-compose ps postgres

# View logs
docker-compose logs postgres
```

### Proxmox API Error
```bash
# Verify Proxmox credentials
curl -k https://your-proxmox:8006/api2/json/version
```

### Redis Connection Error
```bash
# Check Redis status
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

## 📄 License

MIT License - Free untuk penggunaan personal & komersial

## 👨‍💻 Author

Hilmi Muharromi - https://github.com/hilmimuharromi

## 🤝 Kontribusi

Pull requests welcome! Silakan buat issue untuk bug report atau feature request.

---

**Made with ❤️ in Indonesia**
