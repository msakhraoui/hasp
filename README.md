🛡️ InDeception
AI-Assisted Deception & Autonomous Defense System
An autonomous cybersecurity infrastructure that lures attackers into a controlled environment, analyzes their behavior using a local LLM, and autonomously enforces defensive actions based on validated risk scoring.
PFE Project — IRS3 | ISIMA | 2025–2026
Status: 🚧 Active Development
🔍 Overview
InDeception transforms passive honeypots into an intelligent defense platform. Instead of merely logging attacks, it:
Deploys deceptive services (SSH, network) inside an isolated Ubuntu VM
Captures and cryptographically signs all telemetry at ingestion
Sends normalized events to a local LLM for behavioral analysis and risk scoring
Enforces time-limited firewall blocks via a deterministic Go policy engine
Provides real-time visibility through a live SOC dashboard
All detection, orchestration, and enforcement occur exclusively within the Ubuntu VM. The Fedora host serves only as an external LLM inference endpoint and attack simulation source.
graph LR
    A[Attacker] --> B[Honeypot / NIDS]
    B --> C[Go Reader + Ed25519 Sign]
    C --> D[Local LLM Analysis]
    D --> E[Policy Engine]
    E -->|Score ≥ Threshold| F[VM Firewall Block]
    E -->|Log Only| G[(TimescaleDB)]
    F --> G
    G --> H[Live Dashboard]
✅ Current Progress
Implemented
Cowrie SSH Honeypot: Fake terminal capturing commands, credentials (hashed), and file uploads
Snort3 NIDS: JSON-formatted alerts for SYN/NULL/Xmas/FIN scans and SSH brute-force
Go Cowrie Log Reader: Tail-f ingestion with JSON unmarshaling and event normalization
Ed25519 Signing Infrastructure: Key generation, signing at reader level, verification CLI
TimescaleDB Schema: Hypertables for signed events and firewall actions with expiration tracking
Project Skeleton: Modular Go structure (cmd/, internal/), sanitized configs, documentation
Planned / In Progress
Snort JSON log reader
Go orchestrator pipeline (event aggregation, LLM dispatch)
Ollama integration (structured JSON output, prompt engineering)
Policy engine (risk thresholds, allowlist, confidence validation)
Firewall enforcer (ipset/nftables with TTL inside Ubuntu VM)
Real-time WebSocket server + SOC dashboard
End-to-end testing and demo scenario
🔐 Security Principles
Log Integrity: All events signed with Ed25519 at ingestion; tamper-evident audit trail
Confidentiality: Local LLM only; no cloud exfiltration; credentials hashed before storage
Deterministic Enforcement: LLM advises; Go policy engine enforces with minimum confidence + allowlist checks
Isolation: Core security stack confined to Ubuntu VM; host LLM endpoint restricted to internal network
🚀 Getting Started
Prerequisites
Ubuntu 24.04 LTS Server VM
Go 1.22+, Snort 3, Cowrie, TimescaleDB
Fedora/Kali host with Ollama (internal network only)
Quick Start
git clone https://github.com/yourusername/indeception.git
cd indeception

⚠️ Note: Full orchestrator and dashboard are not yet implemented. See Planned / In Progress above.
📄 License
MIT License — see LICENSE for details.
<div align="center">
<sub>Built for PFE — IRS3 program, ISIMA | 2025–2026</sub>
</div>
