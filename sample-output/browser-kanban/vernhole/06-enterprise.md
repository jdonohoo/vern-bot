

# ENTERPRISE GOVERNANCE REVIEW: Browser Extension Kanban Todo Application

## Document Classification: INTERNAL — CONFIDENTIAL
### Review ID: ERB-2026-0209-KANBAN-001
### Status: PENDING ARCHITECTURE REVIEW BOARD APPROVAL
### Distribution: All Stakeholders (see RACI Matrix, Section 4)

---

## 1. EXECUTIVE STAKEHOLDER BRIEFING

Let me be direct: I've reviewed the prior discovery plan, and while the technical analysis is *competent*, it is alarmingly deficient in governance, compliance, and enterprise process rigor. Before a single line of code is written — before we even *discuss* writing code — we need to address the following critical gaps.

The prior plan says "avoid over-abstracting; a small extension deserves a small architecture." Per the governance framework... **no**. A small extension that handles user data, runs in a browser security context, and could theoretically be distributed to millions of users deserves a *proportionate* architecture — which includes process, compliance, and audit trail from Day Zero.

Who are the stakeholders? That question wasn't even asked. We'll need to take this to the architecture review board.

---

## 2. GOVERNANCE GAP ANALYSIS

The prior discovery plan fails to address the following **mandatory** enterprise requirements:

| Gap Area | Severity | Status |
|---|---|---|
| RACI Matrix | CRITICAL | Missing entirely |
| Change Advisory Board (CAB) approval process | CRITICAL | Not referenced |
| Data Classification | HIGH | User task data unclassified |
| Privacy Impact Assessment (PIA) | HIGH | Not conducted |
| SOC 2 Compliance mapping | HIGH | Not assessed |
| GDPR/CCPA data handling procedures | HIGH | Not addressed |
| Vendor Risk Assessment (dependencies) | MEDIUM | Not performed |
| Business Continuity Plan (BCP) | MEDIUM | Not documented |
| Disaster Recovery Plan (DRP) | MEDIUM | Not documented |
| SLA Definitions | MEDIUM | Not defined |
| Rollback Procedures | MEDIUM | Mentioned but not formalized |
| Security Review signoff | HIGH | Not scheduled |
| Legal review of markdown rendering (IP/licensing) | MEDIUM | Not initiated |
| Accessibility compliance (WCAG 2.1 AA) | HIGH | Not mentioned at all |

Has legal signed off on this? I didn't think so.

---

## 3. FORMAL REQUIREMENTS ADDENDUM

The prior plan's requirements are a *starting point*. Below are the enterprise-grade additions required before we proceed to Architecture Review Board submission.

### 3.1 Data Classification

All task data created by users must be classified:

- **Title**: PII-Adjacent (could contain names, project references)
- **Description**: Potentially Sensitive (could contain links to internal systems, credentials if user is careless)
- **Due Date**: Business Confidential (reveals project timelines)
- **Estimate**: Business Confidential (reveals resource allocation)

**Implication**: Even with local-only storage, we must treat `chrome.storage.local` as a sensitive data store. Encryption at rest should be evaluated. At minimum, we need a data retention policy.

### 3.2 Accessibility Compliance (WCAG 2.1 AA)

The prior plan does not mention accessibility **once**. This is a compliance failure. Requirements:

- All color indicators (Green/Yellow/Red) MUST have non-color alternatives (icons, text labels, patterns)
- Keyboard navigation for all Kanban operations
- Screen reader compatibility (ARIA labels on all interactive elements)
- Minimum contrast ratios for both light and dark themes
- Focus management in the task detail panel
- Calendar picker must be keyboard-accessible

This is not optional. This is compliance.

### 3.3 Internationalization (i18n) Readiness

- Date format handling (MM/DD vs DD/MM vs YYYY-MM-DD)
- Text direction support (LTR/RTL)
- String externalization for future localization
- Time zone awareness documented and tested

### 3.4 Audit Trail

Every task state change must be logged:

```
{
  "taskId": "string",
  "action": "create|update|move|delete",
  "field": "string (if update)",
  "previousValue": "any",
  "newValue": "any",
  "timestamp": "ISO-8601",
  "source": "user|system"
}
```

"But it's just a personal todo app" — today it is. Tomorrow someone will ask for team sharing, and if we don't have an audit trail, we'll be retrofitting it under production pressure. Plan for enterprise scale from day one.

### 3.5 SLA Definitions

| Operation | Target | Measurement |
|---|---|---|
| Popup open to interactive | < 300ms | P95 |
| Task creation | < 100ms | P95 |
| Task move between columns | < 50ms | P95 |
| Theme switch | < 100ms (no flash of unstyled content) | P99 |
| Storage write confirmation | < 200ms | P95 |
| Markdown render | < 150ms for 5KB document | P95 |

---

## 4. RACI MATRIX

We'll need a RACI matrix for this. Here it is:

| Activity | Product Owner | Tech Lead | Security | Legal | QA | UX | Accessibility |
|---|---|---|---|---|---|---|---|
| Requirements signoff | **A** | R | C | C | I | C | C |
| Architecture design | C | **A/R** | C | I | I | C | C |
| Security review | I | R | **A** | C | I | I | I |
| Privacy impact assessment | C | R | R | **A** | I | I | I |
| Accessibility audit | I | R | I | I | R | C | **A** |
| Theme implementation | I | R | I | I | R | **A** | C |
| Markdown sanitization | I | R | **A** | C | R | I | I |
| Data model approval | C | **A/R** | C | I | C | I | I |
| Browser store submission | **A** | R | C | R | R | I | I |
| Change management | **A** | R | C | C | C | I | I |

**R** = Responsible, **A** = Accountable, **C** = Consulted, **I** = Informed

---

## 5. RISK ASSESSMENT AND MITIGATION PLAN

The prior plan identified 4 risks. That is insufficient. A proper enterprise risk register follows:

| ID | Risk | Likelihood | Impact | Severity | Mitigation | Owner |
|---|---|---|---|---|---|---|
| R-001 | XSS via markdown rendering | Medium | Critical | HIGH | DOMPurify + CSP headers + security review | Security Lead |
| R-002 | Data loss on extension update | Low | High | MEDIUM | Schema versioning + migration scripts + backup/export | Tech Lead |
| R-003 | Chrome storage quota exceeded | Low | Medium | LOW | Monitoring + archival of Done tasks + user warning at 80% | Tech Lead |
| R-004 | Accessibility lawsuit / complaint | Medium | High | HIGH | WCAG 2.1 AA audit pre-launch | Accessibility Lead |
| R-005 | Third-party dependency vulnerability | Medium | High | HIGH | Vendor risk assessment + dependency audit + lockfile | Security Lead |
| R-006 | GDPR right-to-erasure request | Low | Medium | MEDIUM | Export/delete-all functionality | Product Owner |
| R-007 | Browser API deprecation (Manifest V3 changes) | Medium | Medium | MEDIUM | Abstraction layer for storage/extension APIs | Tech Lead |
| R-008 | Scope creep (team features, sync, etc.) | High | Medium | HIGH | Strict change advisory board process | Product Owner |
| R-009 | Inconsistent urgency coloring across time zones | Medium | Low | LOW | All date comparisons in local time, documented | Tech Lead |
| R-010 | Dark/light theme flash on popup open | Medium | Low | LOW | Load theme preference synchronously before render | Tech Lead |
| R-011 | Clipboard API permission denied (share feature) | Medium | Medium | MEDIUM | Fallback: select-all text in modal | Tech Lead |
| R-012 | No rollback plan for browser store deployment | High | High | HIGH | Documented rollback procedure (see Section 7) | Release Mgr |

---

## 6. VENDOR RISK ASSESSMENT (DEPENDENCIES)

Every external dependency requires evaluation. Preliminary assessment for likely candidates:

| Dependency | Purpose | License | Maintenance Status | Risk Level | Recommendation |
|---|---|---|---|---|---|
| marked / markdown-it | Markdown parsing | MIT | Active | MEDIUM | Approved with sanitizer |
| DOMPurify | HTML sanitization | Apache 2.0 / MPL | Active | LOW | Approved — mandatory |
| Native date input | Calendar picker | N/A (browser built-in) | N/A | LOW | Preferred over third-party |
| Any CSS framework | Styling | Varies | Varies | MEDIUM | Evaluate; vanilla CSS preferred |
| Any drag-and-drop library | Column moves | Varies | Varies | MEDIUM | Defer to Phase 2; vendor review required |

**Recommendation**: Minimize dependencies. Each dependency is a supply chain risk vector. Is this SOC 2 compliant? Only if we audit every dependency.

---

## 7. CHANGE MANAGEMENT AND ROLLBACK PROCEDURES

### 7.1 Change Advisory Board (CAB) Process

All changes must follow the standard CAB workflow:

1. **RFC Submission** — Developer submits Request for Change with impact assessment
2. **CAB Review** — Weekly review meeting (Tuesdays, 2:00 PM)
3. **Approval / Rejection / Deferral** — Documented in change log
4. **Implementation Window** — Scheduled deployment window
5. **Post-Implementation Verification** — Smoke tests within 1 hour
6. **Post-Implementation Review** — 48-hour retrospective

### 7.2 Rollback Plan

What's the rollback plan? Glad you asked:

- **Extension Store**: Maintain previous version build artifacts for immediate re-submission
- **Data Migration**: All schema migrations must be forward AND backward compatible for at least one version
- **Feature Flags**: All Phase 2+ features behind toggles for kill-switch capability
- **Rollback Decision Authority**: Product Owner (business hours) / On-call Lead (off-hours)
- **Rollback SLA**: Decision within 30 minutes of incident detection, execution within 2 hours

---

## 8. ENVIRONMENT STRATEGY

The prior plan did not specify environments. This is unacceptable. Minimum three environments:

| Environment | Purpose | Storage | Distribution |
|---|---|---|---|
| **DEV** | Development and unit testing | Mock / ephemeral | Developer local install |
| **STAGING** | Integration testing, UAT, accessibility audit | Isolated chrome profile | Internal testers via sideload |
| **PRODUCTION** | End user distribution | Real `chrome.storage.local` | Chrome Web Store (+ Edge Add-ons) |

Automated promotion pipeline: DEV → STAGING → PROD with gate checks at each stage (tests pass, security scan clean, accessibility audit green, CAB approval documented).

---

## 9. PHASED IMPLEMENTATION PLAN (REVISED WITH GOVERNANCE GATES)

The prior plan's phasing is reasonable in *content* but lacks governance gates. Revised:

### Phase 0: Governance Foundation (Week 1-2)
- [ ] Complete RACI matrix — **DONE** (see Section 4)
- [ ] Complete risk register — **DONE** (see Section 5)
- [ ] Complete vendor risk assessment — **DONE** (see Section 6)
- [ ] Privacy Impact Assessment — **PENDING LEGAL**
- [ ] Data classification signoff — **PENDING SECURITY**
- [ ] Architecture Review Board submission and approval
- [ ] Accessibility requirements documented

### Phase 1: Core MVP (Week 3-5)
- Extension scaffold (Manifest V3)
- Kanban board with 4 columns
- Task CRUD (Title, Status, Due Date, Estimate)
- Urgency color indicators (with accessible alternatives)
- Theme toggle (light/dark)
- Local persistence with schema versioning
- **GATE**: Security review + QA signoff + Accessibility spot-check

### Phase 2: Enhanced Features (Week 6-8)
- Task detail panel with description editing
- Markdown editor + sanitized preview
- Share-as-markdown button
- Drag-and-drop between columns
- Audit trail logging
- **GATE**: Full accessibility audit + Security re-review + CAB approval

### Phase 3: Distribution (Week 9-10)
- "Open in Tab" expanded view
- Export/import functionality
- Browser store submission (Chrome, Edge)
- **GATE**: Legal review + Final security scan + Production readiness review

### Phase 4: Future Considerations (Backlog — requires separate RFC)
- Search/filter
- Recurring tasks
- Team sharing / sync (requires entirely new security model — let me schedule a meeting to discuss)
- Mobile companion

---

## 10. CRITICAL OBSERVATIONS ON THE PRIOR PLAN

### 10.1 What the Prior Plan Got Right
- Four-column Kanban model is appropriate
- Data model is reasonable for v1
- Urgency logic (end-of-day interpretation) is correct
- Local-first architecture is the right call
- Minimal permissions philosophy is sound

### 10.2 What the Prior Plan Got Wrong or Missed
1. **Zero governance process** — No CAB, no RACI, no change management
2. **No accessibility** — This is a compliance risk that could block distribution
3. **"Avoid over-abstracting"** — This philosophy leads to technical debt that compounds. Proportionate architecture is not over-architecture
4. **No environment strategy** — Shipping directly to production without staging is reckless
5. **No rollback plan formalized** — "Mentioned but not detailed" is not a plan
6. **Security review mentioned once, never scheduled** — Security is a process, not a checkbox
7. **No data classification** — How can we protect data we haven't classified?
8. **No SLA definitions** — How do we measure success without targets?
9. **Markdown XSS acknowledged but not mandated** — DOMPurify should be REQUIRED, not optional
10. **No audit trail** — We will regret this in Phase 4

### 10.3 Contradictions Resolved
- The prior plan wavers between "keep it simple" and comprehensive feature lists. **Resolution**: Strict phasing with governance gates. Every phase delivers value; every phase has approval checkpoints.
- Open decision on "array vs map storage" — **Resolution**: Use a task map (`Record<string, Task>`) with a separate column order array. The marginal complexity is negligible; the benefits for drag-and-drop and future features are substantial. This is a one-time decision, not ongoing complexity.

---

## 11. FINAL GOVERNANCE CHECKLIST

Before ANY development begins:

- [ ] Architecture Review Board approval
- [ ] Privacy Impact Assessment completed
- [ ] Data classification signed off
- [ ] Accessibility requirements documented and approved
- [ ] Vendor risk assessments for all dependencies
- [ ] Environment strategy confirmed
- [ ] CI/CD pipeline with gate checks designed
- [ ] Rollback procedures documented and tested
- [ ] SLAs defined and measurement instrumented
- [ ] RACI matrix distributed to all stakeholders
- [ ] Change Advisory Board cadence established
- [ ] Legal review of open-source license compatibility
- [ ] Security review scheduled (not just "planned")

---

## 12. RECOMMENDATION

**Conditional Approval** — This project may proceed to Phase 0 (Governance Foundation) immediately. Progression to Phase 1 (Core MVP) is contingent upon completion of all Phase 0 governance gates and Architecture Review Board signoff.

The technical plan is sound. The governance plan was nonexistent. Now it exists. We are ready to schedule the meeting to discuss the agenda for the Architecture Review Board submission meeting.

---

*Why did the Enterprise Architect refuse to deploy the todo app on a Friday? Because the Change Advisory Board only meets on Tuesdays, the rollback plan hadn't been approved, the security review was pending, the RACI matrix listed 14 stakeholders who needed to be informed, and by the time all sign-offs were collected... it was Friday again. The joke is: Why do enterprise architects love Kanban boards? Because even their jokes need to move through `New → Todo → In Progress → Done`... and that last column requires three approvals.*

-- Enterprise Vern *(per the governance framework)*
