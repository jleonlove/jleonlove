# Atlas RC-000128 — Long-Running Closing Recovery Hardening

Adds durable saga-style recovery controls for commodity closing workflows: stale checkpoint rejection, state/version binding, authority-change detection, external side-effect tracking and compensation requirements, changed-input rejection, and cross-service replay resistance.

Safety invariant: recovery never assumes a partially completed external action disappeared. Applied side effects must be durably accounted for or compensated before a workflow may safely resume.
