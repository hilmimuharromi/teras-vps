# Teras VPS Design System

## 🎯 Design Principles

* Simple & Professional
* Fast & Lightweight
* Developer Friendly
* Dark Mode First
* Infrastructure Focused

---

# 🎨 Color System

## Primary

* Primary: `#2563EB` (Blue)
* Primary Hover: `#1D4ED8`
* Primary Light: `#DBEAFE`

## Neutral

* Background: `#0B0F19`
* Card: `#111827`
* Border: `#1F2937`
* Text Primary: `#E5E7EB`
* Text Secondary: `#9CA3AF`

## Status Colors

* Success: `#10B981`
* Warning: `#F59E0B`
* Danger: `#EF4444`
* Info: `#3B82F6`

---

# 🔤 Typography

## Font

Primary Font:

* Inter
* System UI fallback

```css
font-family: Inter, system-ui, -apple-system, sans-serif;
```

## Font Size

| Type      | Size |
| --------- | ---- |
| Heading 1 | 28px |
| Heading 2 | 22px |
| Heading 3 | 18px |
| Body      | 14px |
| Small     | 12px |

---

# 📦 Spacing

Base unit: 4px

| Token | Size |
| ----- | ---- |
| xs    | 4px  |
| sm    | 8px  |
| md    | 12px |
| lg    | 16px |
| xl    | 24px |
| 2xl   | 32px |

---

# 🔘 Buttons

## Primary Button

* Background: Primary
* Text: White
* Radius: 8px
* Height: 36px

Example:

```html
<button class="btn-primary">Create VPS</button>
```

---

## Secondary Button

* Background: Transparent
* Border: 1px solid border
* Text: Primary

---

## Danger Button

* Background: Danger
* Text: White

---

# 🧱 Card

Card digunakan untuk:

* VPS Instance
* Billing
* Usage
* Monitoring

Style:

* Background: Card
* Border: Border
* Radius: 12px
* Padding: 16px

Example:

```html
<div class="card">
  <h3>VPS Instance</h3>
</div>
```

---

# 📊 Table

Digunakan untuk:

* VPS List
* Billing
* Logs

Style:

* Header: Bold
* Border Bottom
* Hover Row

---

# 🧭 Navigation

## Sidebar

Sidebar kiri:

Menu:

* Dashboard
* VPS
* Snapshot
* Networking
* Billing
* Settings

Style:

* Width: 240px
* Background: Background
* Border Right

---

# 🧠 Status Badge

## Running

* Background: Success
* Text: White

## Stopped

* Background: Gray

## Error

* Background: Danger

Example:

```html
<span class="badge-running">Running</span>
```

---

# 🖥️ Layout

## Main Layout

```
Sidebar | Topbar
        | Content
```

---

# 📊 Dashboard Widgets

Widgets:

* Total VPS
* CPU Usage
* Memory Usage
* Bandwidth
* Billing

---

# 🎛️ Form

Input Style:

* Background: Card
* Border: Border
* Radius: 8px
* Height: 36px

---

# 🔔 Notification

Position:

* Top Right

Type:

* Success
* Warning
* Error

---

# 🌙 Dark Mode (Default)

Dark mode default:

* Background dark
* Card dark
* Text light

---

# 🧩 Component Naming

Gunakan:

```
vps-card
vps-table
vps-sidebar
vps-button
```

---

# 🚀 Future Components

* VPS Create Wizard
* Metrics Chart
* SSH Console
* Logs Viewer
* Billing Usage

---

# 🎯 UI Inspiration

* DigitalOcean
* Hetzner Cloud
* Vercel
* Railway
