# Python insecure app source reference

[Source hashes](source-hashes.json) pin this reference. The README names dependencies, hardcoded
secrets and SSTI; dependency/secret analysis is planned for future branches. SSTI is not itself an XSS verdict.

| Case | Source | Expected and rationale |
| --- | --- | --- |
| PI-XSS-01 | `app/main.py`, `content` then `Template(content).render()` | `name` is an HTTP query parameter; f-string places it into HTML and then Jinja template source. The route declares `response_class=HTMLResponse`. Raw rendering candidate; Jinja template compilation is not an escaping transformation for inserted template source. |
| PI-OUT-01 | `requests.get(config.PUBLIC_IP_SERVICE_URL)` | Configuration controls destination; no direct request-to-destination flow. Negative for request-influenced outbound finding. |
| PI-XSS-02 | `public_ip_response.text` in the same HTML | Remote response is an independent untrusted source, not a copy of the request URL. Exception fallback is a constant. |

Existing Semgrep and Bearer snapshots both report zero. This does not negate the source-backed
HTML/template flow. No runtime exploit or application execution is part of this reference.

The [documented-family inventory](../python-documented-claims.md) records all advertised families,
additional local/upstream material reviewed, and current unsupported boundaries.

The living [detection guide and answer key](../../../.zazz/evaluation/python-insecure-app-detection-guide.md)
expands the source traces, controls and future-family checklist beyond these original comparison cases.
