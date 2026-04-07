# TerasVPS - Project Setup Summary

**Date:** April 7, 2026
**Project Name:** teras-vps
**Architecture:** Go Fiber (Backend) + Astro (Frontend)

---

## ✅ What Has Been Created

### Project Structure
```
teras-vps/                          ← Monorepo root
├── backend/                        # Go Fiber API Server
│   ├── main.go                     # Entry point
│   ├── go.mod                      # Go module definition
│   ├── go.sum                      # Dependencies lock
│   ├── .env.example                # Environment variables template
│   ├── Dockerfile                  # Backend Docker image
│   ├── config/                     # Configuration
│   │   └── config.go               # Config loader
│   ├── database/                   # Database layer
│   │   ├── postgres.go             # PostgreSQL connection
│   │   ├── redis.go                # Redis connection
│   │   └── migrations/
│   │       ├── 001_init.sql        # Database schema
│   │       └── 002_seed.sql        # Seed data
│   ├── models/                     # Data models
│   │   ├── user.go                 # User model
│   │   ├── vm.go                   # VM model
│   │   ├── plan.go                 # Plan model
│   │   ├── invoice.go              # Invoice model
│   │   ├── transaction.go          # Transaction model
│   │   ├── backup.go               # Backup model
│   │   ├── ssh_key.go              # SSH key model
│   │   ├── audit_log.go            # Audit log model
│   │   └── models.go               # Models import
│   ├── controllers/                # API controllers
│   │   ├── auth_controller.go      # Auth endpoints
│   │   ├── vm_controller.go        # VM endpoints
│   │   ├── billing_controller.go   # Billing endpoints
│   │   ├── user_controller.go      # User endpoints
│   │   └── admin_controller.go     # Admin endpoints
│   ├── routes/                     # Route definitions
│   │   └── routes.go               # All routes setup
│   ├── services/                   # Business logic (TODO)
│   ├── middleware/                 # Middleware (TODO)
│   ├── utils/                      # Utilities (TODO)
│   └── proxmox/                    # Proxmox client (TODO)
│
├── frontend/                       # Astro Frontend
│   ├── package.json                # Node dependencies
│   ├── astro.config.mjs            # Astro configuration
│   ├── tailwind.config.cjs         # TailwindCSS configuration
│   ├── Dockerfile                  # Frontend Docker image
│   ├── nginx.conf                  # Nginx configuration
│   ├── src/
│   │   ├── layouts/
│   │   │   └── Layout.astro       # Base layout
│   │   ├── pages/
│   │   │   └── index.astro        # Landing page
│   │   ├── dashboard/              # Dashboard pages (TODO)
│   │   ├── components/             # React components (TODO)
│   │   ├── lib/                    # API client (TODO)
│   │   └── styles/
│   │       └── global.css         # Global styles
│   └── public/                    # Static assets
│       └── images/                 # Images
│
├── scripts/                        # Deployment scripts
│   └── deploy.sh                   # Deployment script
│
├── docker-compose.yml              # Docker Compose config
├── .gitignore                      # Git ignore rules
├── README.md                       # Project documentation
└── PROJECT_SETUP_SUMMARY.md        # This file
```

---

## 📊 Database Schema

### Tables Created:
- ✅ `users` - User accounts
- ✅ `plans` - Pricing plans
- ✅ `vms` - Virtual machines
- ✅ `invoices` - Billing invoices
- ✅ `transactions` - Payment transactions
- ✅ `vm_backups` - VM backups
- ✅ `ssh_keys` - SSH keys
- ✅ `audit_logs` - Audit logs

### Default Plans (Seed Data):
- **Starter:** 1 Core / 1GB RAM / 20GB Disk - Rp 50.000/month
- **Standard:** 2 Core / 2GB RAM / 40GB Disk - Rp 100.000/month
- **Premium:** 4 Core / 4GB RAM / 80GB Disk - Rp 200.000/month

---

## 🔧 Configuration

### Environment Variables (.env.example)
```bash
# Server
PORT=3000
FRONTEND_URL=http://localhost:4321

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
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24

# Proxmox
PROXMOX_HOST=https://your-proxmox-ip:8006/api2/json
PROXMOX_USER=root@pam
PROXMOX_PASSWORD=your_password
PROXMOX_NODE=proxmox

# Business Logic
MAX_VMS_PER_USER=5
SUSPEND_DAYS=7
DELETE_DAYS=14
MAX_BACKUPS_PER_VM=3

# Currency
CURRENCY=IDR
```

---

## 🚀 API Endpoints (Defined in routes)

### Auth
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get current user

### VMs
- `GET /api/v1/vms` - List user VMs
- `POST /api/v1/vms` - Create new VM
- `GET /api/v1/vms/:id` - Get VM details
- `PUT /api/v1/vms/:id` - Update VM
- `DELETE /api/v1/vms/:id` - Delete VM
- `POST /api/v1/vms/:id/start` - Start VM
- `POST /api/v1/vms/:id/stop` - Stop VM
- `POST /api/v1/vms/:id/reboot` - Reboot VM
- `GET /api/v1/vms/:id/stats` - Get VM stats
- `POST /api/v1/vms/:id/backup` - Create backup
- `GET /api/v1/vms/:id/backups` - List backups

### Billing
- `GET /api/v1/billing/invoices` - List invoices
- `GET /api/v1/billing/invoices/:id` - Get invoice
- `POST /api/v1/billing/invoices/:id/pay` - Pay invoice
- `GET /api/v1/billing/plans` - List plans

### SSH Keys
- `GET /api/v1/ssh-keys` - List SSH keys
- `POST /api/v1/ssh-keys` - Add SSH key
- `DELETE /api/v1/ssh-keys/:id` - Delete SSH key

### User
- `GET /api/v1/user/profile` - Get profile
- `PUT /api/v1/user/profile` - Update profile
- `PUT /api/v1/user/password` - Change password

### Admin
- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/vms` - List all VMs
- `GET /api/v1/admin/stats` - Get platform stats
- `POST /api/v1/admin/suspend-user/:id` - Suspend user
- `POST /api/v1/admin/unsuspend-user/:id` - Unsuspend user

---

## 🐳 Docker Deployment

### Start All Services:
```bash
docker-compose up -d
```

### Stop All Services:
```bash
docker-compose down
```

### View Logs:
```bash
docker-compose logs -f
```

---

## 📝 Next Steps (TODO)

### Phase 1: Authentication (Week 2)
- [ ] Implement JWT middleware
- [ ] Implement bcrypt password hashing
- [ ] Implement register logic
- [ ] Implement login logic
- [ ] Implement logout logic
- [ ] Add role-based access control

### Phase 2: Proxmox Integration (Week 3)
- [ ] Implement Proxmox service
- [ ] Create VM in Proxmox
- [ ] Delete VM from Proxmox
- [ ] Start/Stop/Reboot VM
- [ ] Get VM stats from Proxmox

### Phase 3: Billing System (Week 4)
- [ ] Implement invoice generation
- [ ] Integrate payment gateway (QRIS)
- [ ] Implement auto-suspend logic
- [ ] Implement auto-delete logic

### Phase 4: Frontend Development (Ongoing)
- [ ] Build auth pages (Login, Register)
- [ ] Build dashboard layout
- [ ] Build VM list page
- [ ] Build VM detail page
- [ ] Build billing pages
- [ ] Add Shadcn/ui components

### Phase 5: Real-time Monitoring (Week 5)
- [ ] Implement WebSocket
- [ ] Build real-time charts (Recharts)
- [ ] Add live status updates

### Phase 6: Additional Features (Week 6)
- [ ] SSH keys management
- [ ] Backup system
- [ ] Support ticket system
- [ ] Admin dashboard

### Phase 7: Deployment (Week 7)
- [ ] Deploy to production (mini PC)
- [ ] Setup SSL certificate
- [ ] Configure Nginx
- [ ] Create systemd services
- [ ] Test production environment

---

## 🎯 Business Rules Implemented

- ✅ **Max VMs per user:** 5 (configurable)
- ✅ **Suspend after:** 7 days overdue
- ✅ **Delete after:** 14 days overdue
- ✅ **Max backups per VM:** 3
- ✅ **Currency:** IDR
- ✅ **Language:** English (UI)

---

## 📚 Documentation Files to Create

- [ ] `docs/DEPLOYMENT.md` - Detailed deployment guide
- [ ] `docs/API.md` - Complete API documentation
- [ ] `docs/USER_MANUAL.md` - User manual for customers
- [ ] `docs/ADMIN_GUIDE.md` - Admin guide

---

## ✨ Features to Implement Later (Post-Launch)

- [ ] Payment gateway integration (Xendit QRIS)
- [ ] Mobile app (React Native)
- [ ] Developer API
- [ ] Multiple data center locations
- [ ] Domain + SSL add-ons
- [ ] Advanced backup scheduling
- [ ] VM marketplace (templates)
- [ ] Bandwidth limits
- [ ] IP whitelisting

---

## 📞 Support

For questions or issues, contact:
- **Author:** Hilmi Muharromi
- **Email:** hilmi@example.com
- **GitHub:** https://github.com/hilmimuharromi/teras-vps

---

**Status:** ✅ Project structure created successfully!
**Next:** Start Phase 1 - Authentication implementation
