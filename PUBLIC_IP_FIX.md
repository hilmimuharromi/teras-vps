# Fix: PostgreSQL inet Type Empty String Error

## ❌ Error

```
ERROR: invalid input syntax for type inet: "" (SQLSTATE 22P02)
[4.204ms] [rows:0] INSERT INTO "vms" (..., "public_ip", ...) VALUES (..., '', ...)
```

## 🔍 Root Cause

**PostgreSQL `inet` Type Constraint**

PostgreSQL's `inet` data type does NOT accept empty strings `""`. It only accepts:
- ✅ Valid IP addresses: `"192.168.1.100"`
- ✅ `NULL` values
- ❌ Empty strings: `""` ← **ERROR!**

### The Problem:

In Go, when you declare a string field:
```go
PublicIP string `gorm:"type:inet" json:"public_ip,omitempty"`
```

If you don't set a value, Go defaults it to **empty string `""`**, which PostgreSQL rejects.

## ✅ Fix

### Change `PublicIP` from `string` to `*string` (pointer type)

**File**: `backend/models/vm.go`

**Before**:
```go
PublicIP string `gorm:"type:inet" json:"public_ip,omitempty"`
// ❌ Empty string "" when not set → PostgreSQL error
```

**After**:
```go
PublicIP *string `gorm:"type:inet" json:"public_ip,omitempty"` // Nullable for inet type
// ✅ NULL when not set → PostgreSQL accepts it
```

### Why Pointer Type Works:

- **Non-pointer string**: Defaults to `""` (empty string)
- **Pointer string (`*string`)**: Defaults to `nil` (becomes `NULL` in database)

When GORM sees a `nil` pointer, it inserts `NULL` instead of empty string.

## 📊 SQL Comparison

### Before (Wrong):
```sql
INSERT INTO "vms" ("public_ip", ...) VALUES ('', ...)
-- ERROR: invalid input syntax for type inet: ""
```

### After (Correct):
```sql
INSERT INTO "vms" ("public_ip", ...) VALUES (NULL, ...)
-- ✅ Success! NULL is valid for inet type
```

## 🎯 When to Use Each Type

### Use `string` (non-pointer):
- When the field is **required**
- When you want to ensure it always has a value
- When default values are acceptable

```go
Name string `gorm:"not null" json:"name"` // Always has a value
```

### Use `*string` (pointer):
- When the field is **optional**
- When you want to store `NULL` in database
- When dealing with PostgreSQL special types (`inet`, `uuid`, etc.)

```go
PublicIP *string `gorm:"type:inet" json:"public_ip,omitempty"` // Can be NULL
```

## 💡 Usage Examples

### Creating VM Without IP:
```go
vm := models.VM{
    UserID:   1,
    VMID:     100,
    Name:     "my-vps",
    PublicIP: nil, // ✅ Will be NULL in database
}
db.Create(&vm)
```

### Creating VM With IP:
```go
ip := "192.168.1.100"
vm := models.VM{
    UserID:   1,
    VMID:     101,
    Name:     "my-vps",
    PublicIP: &ip, // ✅ Will be "192.168.1.100" in database
}
db.Create(&vm)
```

### Setting IP Later:
```go
// Get VM
var vm models.VM
db.First(&vm, 1)

// Set IP
newIP := "10.0.0.1"
vm.PublicIP = &newIP
db.Save(&vm)
```

### Clearing IP:
```go
vm.PublicIP = nil // Set back to NULL
db.Save(&vm)
```

## 🔍 Other Fields That Might Have This Issue

Check all PostgreSQL `inet`, `cidr`, `macaddr`, `uuid` columns:

### Common Fields Needing Pointer Types:
- ✅ `PublicIP *string` - IP addresses (fixed)
- ⚠️ Any `*string` for optional text fields
- ⚠️ Any `*time.Time` for optional dates

### Fields That Are OK as Non-Pointer:
- ✅ Required fields with `not null` constraint
- ✅ Fields with default values
- ✅ Fields that always get set in BeforeCreate hook

## ✅ Verification

### Test VM Creation:
```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}' | jq -r '.data.token')

# Create VM (without public IP)
curl -X POST http://localhost:3000/api/v1/vms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"test-vm","hostname":"test-vm","plan_id":1}'
```

**Expected:** ✅ VM created successfully without inet error!

### Check Database:
```sql
SELECT id, name, public_ip FROM vms ORDER BY id DESC LIMIT 1;
```

**Expected:**
```
 id |  name   | public_ip
----+---------+-----------
  1 | test-vm | NULL
(1 row)
```

## 📝 Related Issues

This is the **second** column issue fixed today:

1. **First Issue**: Column name mismatch `vm_id` vs `vmid`
   - **Fix**: Added `column:vmid` tag

2. **Second Issue**: Invalid inet value `""` vs `NULL`
   - **Fix**: Changed `string` to `*string`

Both issues stem from GORM's default behavior not matching PostgreSQL requirements.

## 📁 Files Changed

- ✅ `backend/models/vm.go` - Changed `PublicIP` from `string` to `*string`

## 🎓 Lessons Learned

### PostgreSQL Strict Type Checking:
PostgreSQL is very strict about data types:
- `inet` must be valid IP or NULL
- `uuid` must be valid UUID or NULL
- `json/jsonb` must be valid JSON or NULL
- `boolean` must be true/false/NULL

### GORM Default Values:
Go initializes types with defaults:
- `string` → `""` (empty string)
- `int` → `0`
- `bool` → `false`
- `*string` → `nil` (NULL in DB)
- `*int` → `nil` (NULL in DB)

### Best Practice:
**Use pointer types for optional fields** that map to PostgreSQL columns that:
- Don't accept empty strings
- Can be NULL
- Have special format requirements (IP, UUID, etc.)

## ✨ Result

VM creation now works correctly! PublicIP will be `NULL` when not set, avoiding the inet type error.

### Before:
```
ERROR: invalid input syntax for type inet: ""
```

### After:
```
✅ VM created successfully!
```

PublicIP field is now properly nullable, matching the database schema where `public_ip` is an optional field.
