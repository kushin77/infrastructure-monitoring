# Copilot Instructions — infrastructure-monitoring

<!-- BEGIN shared-governance:scope-boundary-block (Phase 12) -->
# Repository Scope & Boundary Rules (GOVERNANCE FIX)

**PURPOSE**: Prevent agent auto-redirect across repositories without explicit user authorization. This section must appear in EVERY tenant repo's `.github/copilot-instructions.md`.

---

## Primary Scope

This repository is `/home/akushnir/infrastructure-monitoring` only. Examples:
- shared-services: `/home/akushnir/shared-services`
- gohighlevel: `/home/akushnir/gohighlevel`
- elevatediq: `/home/akushnir/elevatediq`
- infrastructure-monitoring: `/home/akushnir/infrastructure-monitoring`
- shared-governance: `/home/akushnir/shared-governance`

**All work, by default, is scoped to this repository ONLY.**

---

## NO AUTO-REDIRECT RULE

**Do NOT assume work in other repositories** unless EXPLICITLY scoped by:
1. The current GitHub issue body (states which repo to work in)
2. The user's direct request (user says "work in shared-governance")
3. An interrepo dependency documented in `docs/cross-repo-boundaries.md` (with explicit authorization)

**Standing rules from user memory do NOT override repo-scoped instructions.**

If you are uncertain whether a task belongs in this repo or another, ask the user:
> **"Is this work for infrastructure-monitoring or another repo? I want to confirm scope before proceeding."**

---

## Explicit Scope Required

| Scenario | Action |
|----------|--------|
| User says "complete all issues" | Assume infrastructure-monitoring issues only — do NOT redirect to other repos |
| User says "fix shared-governance issues" | Work ONLY in shared-governance repo (do NOT also work in this one) |
| Issue mentions cross-repo work | Read issue body for EXPLICIT scope directive |
| Standing rule from user memory applies | VERIFY it applies to THIS repo first before using it |
| Unclear whether cross-repo work is authorized | ASK the user for clarification |

---

## Governance Principle

**Hierarchy of Authority** (highest to lowest):
1. Repo-scoped instructions (`.github/copilot-instructions.md`)
2. User's current explicit request ("work in X")
3. Issue body scope (if documented)
4. User memory standing rules (lowest priority — must be repo-specific)

---

## Cross-Repo Work (Rare Exception)

Cross-repo work is ONLY allowed when:
1. The issue explicitly mentions it in the body
2. The cross-repo dependency is documented in `docs/cross-repo-boundaries.md`
3. The user explicitly authorizes it in their request

**Example**: "Phase 51 deployment in shared-services requires Vault setup in gohighlevel. Do Phase 51 first, then wire Vault."

In this case, read the issue body to understand the EXPLICIT scope, then proceed.

---

## If Scope Violation Detected

If you notice you've been working in the wrong repo:
1. Stop immediately
2. Document what was done in which repo
3. Report the scope violation to the user
4. Ask for re-scoping authority before continuing

This prevents the pattern of silent scope-creep that led to the June 29, 2026 incident (Phases 119-186 completed in wrong repo).

---

<!-- END shared-governance:scope-boundary-block (Phase 12) -->
