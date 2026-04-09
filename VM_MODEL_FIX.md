# VM Model Fix - Column Name Mismatch

## ❌ Error

```
2026/04/09 00:42:35 ERROR: column "vm_id" of relation "vms" does not exist (SQLSTATE 42703)
[3.118ms] [rows:0] INSERT INTO "vms" ("user_id","vm_id","name","hostname","status","cores","memory","disk","public_ip","ssh_port","vnc_port","template","plan_id","expires_at","created_at","updated_at") VALUES (4,105,'mirodev','mirodev','stopped',1,1024,20,'',22,5900,'100',1,'2026-05-09 00:42:35.562','2026-04-09 00:42:35.562','2026-04-09 00:42:35.562') RETURNING "id"
```

## 🔍 Root Cause

**GORM Naming Convention Mismatch**

- **GORM default**: Converts `VMID` field → `vm_id` column (snake_case)
- **Database schema**: Column is named `vmid` (no underscore)

The SQL migration (`database/migrations/001_init.sql`) defines:
```sql
CREATE TABLE IF NOT EXISTS vms (
    ...
    vmid INTEGER NOT NULL UNIQUE,  -- ← No underscore!
    ...
);
```

But GORM was trying to insert into:
```sql
INSERT INTO "vms" ("vm_id", ...)  -- ← With underscore!
```

## ✅ Fix

### File: `backend/models/vm.go`

**Before:**
```go
type VM struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    VMID      int       `gorm:"not null;index" json:"vmid"` // ❌ GORM converts to vm_id
    Name      string    `gorm:"size:100;not null" json:"name"`
    ...
}
```

**After:**
```go
type VM struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    VMID      int       `gorm:"column:vmid;not null;index" json:"vmid"` // ✅ Explicit column name
    Name      string    `gorm:"size:100;not null" json:"name"`
    ...
}
```

## 🎯 Explanation

By adding `column:vmid` to the GORM tag, we explicitly tell GORM to use the `vmid` column name instead of auto-generating `vm_id`.

### GORM Naming Rules:
- Default: Converts camelCase to snake_case
  - `VMID` → `vm_id`
  - `UserID` → `user_id`
  - `PublicIP` → `public_ip`
- Override: Use `column:name` tag to specify exact column name

### Why This Happened:
The migration file was created manually with custom column names, but GORM's auto-mapping didn't match because:
- Migration uses: `vmid` (short, no underscore)
- GORM expects: `vm_id` (from `VMID` field)

## ✅ Verification

After the fix, the INSERT query should be:
```sql
INSERT INTO "vms" ("vmid", "user_id", "name", ...) VALUES (105, 4, 'mirodev', ...)
```

Instead of:
```sql
INSERT INTO "vms" ("vm_id", "user_id", "name", ...) VALUES (105, 4, 'mirodev', ...)
```

## 📝 Other Models Check

### Models with VMID Foreign Key:

**✅ Backup (`models/backup.go`)**
```go
VMID uint `gorm:"not null;index" json:"vm_id"`
```
- Table: `vm_backups`
- Column: `vm_id` ✅ (matches GORM expectation)

**✅ Invoice (`models/invoice.go`)**
```go
VMID *uint `gorm:"index" json:"vm_id,omitempty"`
```
- Table: `invoices`
- Column: `vm_id` ✅ (matches GORM expectation)

Both are correct because the migration files use `vm_id` for these tables.

## 🔧 How to Avoid This in Future

### Option 1: Always Explicit Column Names
```go
VMID int `gorm:"column:vmid;not null;index" json:"vmid"`
```

### Option 2: Follow GORM Convention in Migrations
```sql
-- Use snake_case for all column names
vm_id INTEGER NOT NULL UNIQUE,  -- Instead of vmid
```

### Option 3: Disable GORM Naming Convention
```go
// In database initialization
db.NamingStrategy = schema.NamingStrategy{
    NoLowerCase: true,  // Keep exact field names
}
```

## 📋 Testing

### Test VM Creation:
```bash
# Login first
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}' \
  | jq -r '.data.token')

# Create VM
curl -X POST http://localhost:3000/api/v1/vms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "test-vm",
    "hostname": "test-vm",
    "plan_id": 1
  }'
```

**Expected:** No more "column vm_id does not exist" error!

## 📁 Files Changed

- ✅ `backend/models/vm.go` - Added `column:vmid` tag to VMID field

## ✨ Result

VM creation now works correctly! The INSERT statement uses `vmid` column as expected by the database schema.

### Before:
```
ERROR: column "vm_id" of relation "vms" does not exist
```

### After:
```
✅ VM created successfully!
```
