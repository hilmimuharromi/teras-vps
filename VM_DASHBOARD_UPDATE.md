# VM Dashboard Update - SSH Info & Real-time Graphs

## ✅ What Was Updated

The My VMs page (`/dashboard/vms`) has been redesigned to focus on what matters most:
1. **SSH Connection Info** - Quick access to connect to your VM
2. **Real-time CPU & Memory Graphs** - Monitor VM performance
3. **Clean, Focused UI** - Remove clutter, show what's important

## 🎨 New Features

### 1. SSH Connection Info Box

For running VMs with a public IP, shows:
```
💻 ssh root@192.168.1.100 -p 22     📋 Copy
```

- **One-click copy** - Click "Copy" button to copy SSH command
- **Smart display** - Only shown when VM is running AND has public IP
- **Visual feedback** - Shows "✅ Copied!" temporarily when copied

### 2. Real-time Performance Graphs

Two beautiful line charts for each running VM:

**CPU Usage Graph:**
- Real-time updates every 5 seconds
- Shows last 20 data points (~100 seconds)
- Blue gradient fill
- Y-axis: 0-100%
- Current value displayed at top

**Memory Usage Graph:**
- Real-time updates every 5 seconds
- Shows last 20 data points
- Purple gradient fill
- Y-axis: Memory in MB
- Current value displayed at top

### 3. Clean VM Card Layout

Each VM card now shows:
```
┌─────────────────────────────────────────────┐
│ 🟢 RUNNING    my-vps-server                │
│ hostname: my-vps                            │
│                                             │
│ 💻 ssh root@192.168.1.100 -p 22  📋 Copy  │
│                                             │
│ ┌──────────────┐  ┌──────────────┐         │
│ │ 🖥️ CPU 45%   │  │ 🧠 RAM 512MB │         │
│ │ [graph]      │  │ [graph]      │         │
│ └──────────────┘  └──────────────┘         │
│                                             │
│ 🖥️ 2 Cores  🧠 2.0 GB RAM  💾 40GB Disk   │
│                                      ⏹️ Stop│
└─────────────────────────────────────────────┘
```

## 🔧 Technical Details

### Chart.js Integration

- **CDN**: Loaded from jsDelivr (Chart.js 4.4.0)
- **Lazy loading**: Charts only initialized for running VMs
- **Auto-refresh**: Fetches stats every 5 seconds via `/api/v1/vms/{id}/stats`
- **Smooth updates**: Updates without animation for real-time feel
- **Data retention**: Keeps last 20 data points, removes old ones

### SSH Command Generation

```javascript
const sshCommand = vm.public_ip 
  ? `ssh root@${vm.public_ip} -p ${vm.ssh_port || 22}`
  : 'SSH not available (VM stopped)';
```

### Stats Processing

```javascript
// Proxmox returns CPU as 0-1 scale, convert to percentage
const cpuUsage = Math.round((stats.cpu / 100) * 100);

// Proxmox returns memory in bytes, convert to MB
const memoryUsed = Math.round(stats.mem / 1024 / 1024);
```

## 📊 API Integration

### Endpoint: `/api/v1/vms/{id}/stats`

**Request:**
```http
GET /api/v1/vms/1/stats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "stats": {
      "cpu": 0.45,           // 0-1 scale (45%)
      "mem": 536870912,      // bytes (512 MB)
      "maxmem": 2147483648,  // max bytes (2 GB)
      "disk": 10737418240,   // disk used
      "maxdisk": 42949672960 // max disk
    }
  }
}
```

## 🎯 User Experience Improvements

### Before:
- Basic VM info (name, hostname, specs)
- Status badge
- Start/Stop buttons
- View button (goes to detail page)
- No performance data
- No SSH info

### After:
- ✅ SSH command with copy button
- ✅ Real-time CPU graph
- ✅ Real-time Memory graph
- ✅ Auto-refresh every 5 seconds
- ✅ Clean, focused layout
- ✅ Visual status indicators (🟢🔴🟡)
- ✅ Smart display of RAM (auto-converts MB to GB)

## 🎨 Design Choices

### Color Scheme:
- **CPU**: Blue (#3B82F6) - Represents processing/computation
- **Memory**: Purple (#A855F7) - Represents storage/RAM
- **SSH**: Green (#22C55E) - Represents connection/success
- **Background**: Dark slate (#0F172A) - Reduces eye strain

### Graph Styling:
- **No grid lines on X-axis** - Cleaner look
- **Light grid on Y-axis** - Helps read values
- **Gradient fill** - Modern, polished appearance
- **No points** - Smooth line without dots
- **Curved tension** - Organic, natural look

## 💡 Usage Examples

### SSH to Your VM:
1. Find your running VM in the list
2. Look for the dark SSH box below the VM name
3. Click "📋 Copy" button
4. Paste in your terminal: `ssh root@192.168.1.100 -p 22`
5. Press Enter and connect!

### Monitor Performance:
1. Look at the two graphs below SSH info
2. **CPU graph** shows processing usage in real-time
3. **Memory graph** shows RAM usage in MB
4. Graphs update automatically every 5 seconds
5. Watch for spikes or high usage

### Start/Stop VM:
1. Click "▶️ Start" or "⏹️ Stop" button (top right)
2. Page reloads to update charts
3. SSH box appears when VM is running and has IP

## 🔍 Troubleshooting

### SSH Box Not Showing:
- **VM must be running** - Stopped VMs don't show SSH
- **VM must have public IP** - Check VM details
- **IP might be empty** - Proxmox might not have assigned IP yet

### Graphs Not Loading:
- **Check internet connection** - Chart.js loaded from CDN
- **Check browser console** - Look for Chart.js errors
- **VM must be running** - Stopped VMs don't show graphs
- **Proxmox must be accessible** - Stats come from Proxmox API

### Graphs Not Updating:
- **Check Network tab** - Look for `/api/v1/vms/{id}/stats` requests
- **Check auth token** - Expired token will fail
- **Check Proxmox** - VM must be running to get stats
- **Check backend logs** - Look for errors in terminal

### Copy Button Not Working:
- **Browser support** - Requires `navigator.clipboard` API
- **HTTPS required** - Clipboard API needs secure context
- **Check console** - Look for permission errors

## 📁 Files Modified

- ✅ `frontend/src/pages/dashboard/vms/index.astro` - Complete redesign

## 🚀 Future Enhancements

Potential improvements:
- [ ] Add disk usage graph
- [ ] Add network I/O graphs
- [ ] Add uptime counter
- [ ] Add VM screenshot/snapshot
- [ ] Add quick terminal (web-based SSH)
- [ ] Add alert thresholds (notify when CPU > 90%)
- [ ] Add historical data (24h, 7d, 30d graphs)
- [ ] Add export stats button

## ✨ Result

A clean, focused VM dashboard that shows exactly what you need:
- **How to connect** (SSH command)
- **How it's performing** (CPU & RAM graphs)
- **How to control it** (Start/Stop buttons)

No clutter, just the essentials! 🎯
