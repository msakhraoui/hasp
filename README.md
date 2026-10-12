# InDeception

### LLM-assisted SSH deception with a deterministic Go policy engine. Research prototype.

A honeypot that doesn't just record attacks — it scores attacker sessions with a local LLM and blocks the source at the firewall. The LLM advises; deterministic Go code decides.

**Status:** Active development. Not production software. Run only in an isolated lab you own.

---

## Why

Classic tools like fail2ban block on simple counters (N failed logins) — trivially evaded by slow, distributed, or novel behavior. InDeception tests whether a small local LLM can score full session behavior well enough to drive automated blocking, under three hard constraints:

- **No cloud** — all inference local, no telemetry leaves the host
- **Deterministic enforcement** — the LLM never touches the firewall directly
- **Auditable** — every event is signed (Ed25519) at ingestion; every block is logged with the score and threshold that caused it

The honeypot (Cowrie, port 22) is the sensor. The defended asset is the real SSH service on a nonstandard port on the same host — a source-IP block protects both.

---

## Architecture

```
attacker
   │
   ▼
 :22 ── iptables redirect ──▶ Cowrie :2222 ──┐      (Ubuntu VM)
                                             │
 Snort 3 (NIDS alerts) ──────────────────────┤
                                             ▼
                      Go log readers — tail, normalize, Ed25519 sign
                                             │
                                             ▼
                      Orchestrator — aggregates events into sessions
                            │                        │
                            ▼                        ▼
                Ollama (local LLM)           Go policy engine
                JSON risk score              threshold · allowlist · min confidence
                            │                        │
                            └──────────┬─────────────┘
                                       ▼
                        nftables / ipset — time-limited (TTL) block
                                       │
                                       ▼
             TimescaleDB (signed events + actions) · SOC dashboard (optional)

Fedora host: attack simulation + LLM inference endpoint (internal network only)
```

---

## Research questions

1. Can a small local LLM score SSH session risk well enough to trigger automatic blocks — and at what false-block rate on simulated benign sessions?
2. Does the LLM's score separate attackers from simulated benign sessions better than a rule-based score and a supervised classical baseline?
3. Can an attacker who controls honeypot session content manipulate the LLM's risk score — and does the deterministic policy layer hold?
4. How long from an attacker's first action to a block, and is that fast enough to matter under the threat model?
5. How do accuracy and latency scale with model size (~1B / 3B / 8B) on ordinary hardware?

---

## Threat model

- **Attacker:** remote, no valid credentials, targets SSH services. Adaptive — assumed to know a deception platform may exist, and able to shape session content (including prompt injection against the LLM) to avoid or manipulate scoring.
- **Defended asset:** the real SSH service (nonstandard port). The honeypot is expendable.
- **Trusted:** the VM host, the policy engine binary, the signing keys.
- **Out of scope:** kernel exploits, VM escape, encrypted C2 analysis.

The prompt-injection surface is deliberate: if an attacker can manipulate the score, the policy engine's minimum-confidence check is the last line. Testing that is part of the evaluation.

---

## Evaluation plan (planned)

- **Attack sessions:** port scans, SSH brute force, scripted post-auth command sequences, manual operator sessions — generated from a Kali VM on the internal network using nmap, hydra, and custom scripts.
- **Simulated benign sessions:** scripted admin work (file edits, package updates, log inspection), authored blind — written without reference to the attack scripts — with randomized timing and command ordering. A honeypot has no real users; this is a stated limitation.
- **Prompt injection sessions:** attacker-authored session content crafted to manipulate the LLM score. Measures score drift and policy engine resilience.
- **Splits:** calibration set for threshold selection; held-out set for all reported metrics.
- **Models:** three sizes via Ollama, pinned tags and quantization (e.g. `llama3.1:8b-instruct-q4_K_M`), temperature 0, JSON-schema structured output. Prompt template versioned in the repo.
- **Metrics:** per-session block precision/recall, false-block rate, miss rate, time-to-block (median, P95), LLM latency. Class balance stated explicitly.
- **Baselines:** (1) simple rule-based score, (2) fail2ban-equivalent counter policy, (3) supervised classical ML on session features.

### Results (placeholder)

| Model | Block precision | False-block rate | Miss rate | Time-to-block P95 |
|---|---|---|---|---|
| — | — | — | — | — |

---

## Tech stack

| Layer | Technology |
| :--- | :--- |
| **Core language** | Go 1.22+ |
| **Honeypot** | Cowrie (Python-based) |
| **NIDS** | Snort 3 (Lua rules + JSON output) |
| **AI analysis** | Ollama (local LLM inference) |
| **Database** | TimescaleDB (PostgreSQL extension) |
| **Dashboard** | HTML/CSS/JavaScript + WebSocket (optional) |
| **Infrastructure** | Ubuntu Server 24.04 VM, iptables/nftables, ipset |
| **Signing** | Ed25519 (Go stdlib `crypto/ed25519`) |

---

## Current status

### Implemented
- Cowrie SSH honeypot: fake terminal capturing commands, credential attempts (hashed), file uploads
- Snort 3 NIDS: JSON alerts for SYN/NULL/Xmas/FIN scans and SSH brute force
- Go log readers (`cmd/reader/`): tail-based ingestion, JSON unmarshaling, event normalization

### In progress
- Ed25519 signing at reader level + verification CLI
- Orchestrator: session aggregation, LLM dispatch (structured JSON output)
- Policy engine: thresholds, allowlist, confidence validation
- Firewall enforcer: nftables/ipset with TTL

### Planned
- Evaluation harness: attack + benign scenarios, metrics runner, model comparison
- TimescaleDB hypertables for signed events and actions
- Optional WebSocket SOC dashboard
- End-to-end demo scenario

---

## Repository layout

```
cmd/
  reader/          log tailer — cowrie/snort, normalize, Ed25519 sign, forward
  orchestrator/    session aggregation, LLM dispatch
  policy/          deterministic policy engine
  verifier/        signature verification CLI
internal/
  event/           normalized event types
  sign/            Ed25519 helpers
  llm/             Ollama client, JSON schema, prompt templates
  firewall/        nftables/ipset drivers
eval/              attack + benign scenarios, metrics runner
configs/           thresholds, model pins, deployment settings
deploy/            nftables rules, Cowrie/Snort configs
```

---

## Security design

- **Integrity:** Ed25519 signature per event at ingestion. Tamper-evident, not tamper-proof — the signer lives on the same host.
- **Confidentiality:** local inference only; no telemetry leaves the host; credentials hashed before storage.
- **Determinism:** LLM output is advisory. Blocks require threshold + minimum confidence + allowlist pass. Every action is logged with its full inputs.
- **Containment:** security stack confined to the Ubuntu VM; LLM endpoint restricted to the internal network.
- **Collateral:** IP blocks affect everyone behind a shared address (NAT). Allowlist and short TTLs limit the impact.

---

## Limitations

- Simulated benign traffic is self-authored; results won't transfer to real user distributions.
- Single host, single attacker source; no distributed or slow-burn attacks yet.
- LLM latency dominates time-to-block on the slow path.
- Quantized models; absolute numbers are specific to the pinned versions.
- Signing is tamper-evident only — the signer and the signed data live on the same host.

---

## Ethics & scope

Isolated lab environment only. Do not deploy on networks you do not own or operate. Blocking real source IPs has real consequences for the people behind them. This is a research prototype, not a production security tool.

---

## License

GPL-3.0 — see [LICENSE](LICENSE).

## Contact

Malek SK · malekskwork@gmail.com · feedback and threat-model reviews welcome via GitHub Discussions.
