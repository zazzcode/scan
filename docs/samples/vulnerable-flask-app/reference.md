# Vulnerable Flask App (Erlik 2) source reference

[Hashes](source-hashes.json) pin the author's last commit through the retained fork. README lists
HTML injection, XSS, SSTI, SQL injection and several other types. Source review corrects imprecise
line references in the initial handoff: `/welcome` returns JSON and is a negative.

| Case | Source in `vulnerable-flask-app.py` | Expected and rationale |
| --- | --- | --- |
| VF-SQL-01 | `search_user`, line 20 | sqlite3 cursor executes percent-formatted path parameter in SQL. Construction finding plus injection candidate. |
| VF-XSS-01 | `welcome`, lines 29–32 | Concatenated input is passed to jsonify; no HTML sink. Negative. |
| VF-XSS-02 | `welcome2`, lines 34–37 | Path input returned as a string, Flask default HTML response. Raw reflected output even without literal HTML delimiters. |
| VF-XSS-03 | `hello_ssti`, lines 39–51 | Query value is interpolated into template source before render_template_string. Autoescaping template variables does not escape inserted source. XSS candidate; SSTI family separate. |
| VF-XSS-04 | `get_users`, exception return line 62 | Request hostname becomes raw Flask string response on exception path. Supported witness requires exception precondition. |
| VF-LOG-01 | `ImproperOutputNeutralizationforLogs`, lines 183–189 | Arbitrary `data` logged: log-injection observation, not established sensitive-data leak. No sensitive-data finding expected. |
| VF-LOG-02 | `search_user`, line 25 | Database results logged; row fields/sensitivity unresolved. Review at most, not a proven request-content leak. |
| VF-LOG-03 | `get_admin_mail`, line 109 | Literal internal admin email logged; no request-content transfer. Outside request-sensitive logging slice. |

Semgrep reports SQL, raw responses and command/traversal/configuration cases; Bearer reports SQL
and pickle. The two SQL locations in other apps are not independent vulnerabilities. Command/code
injection, traversal, deserialization, uploads, DoS, brute force and configuration remain out of scope.

The logging slice now records the literal internal admin address as a candidate. It proves the
emission, but email-shaped syntax alone does not prove a private or person-identifying value.
The initial Luna pass classified it as personal data; the final static policy retains Review until
sensitivity is established independently. This is an accuracy correction, not a target-name rule.

The [documented-family inventory](../python-documented-claims.md) records all advertised families,
additional local/upstream material reviewed, and current unsupported boundaries.
