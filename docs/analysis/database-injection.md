# Database injection analysis

[Analysis index](README.md) · [Input validation](input-validation.md) · [Calibration](calibration.md)

## Three independent questions

Scan separates **business validation**, **construction practice**, and **injection protection**.
An executable query string can deserve refactoring even if a numeric conversion prevents injection.
A parameterized query can prevent injection while accepting values that violate business policy.
A query category label by itself answers neither question.

## Evidence before judgment

Analysis MUST identify the request field, its transfers into the actual operation slot, the bound
library/API, relevant guard and middleware behavior, and a supported source-to-operation path.
Names such as `query` or `find` alone MUST NOT establish a database API. An array's `find` is not a
NoSQL sink. Unresolved applicable checks, dynamic bindings, or missing helper bodies remain blockers.

The bounded automatic adapters currently include verified Sequelize raw-query usage and MarsDB
collection construction/imports. MarsDB has Mongo-style matching semantics; this does not make it
a MongoDB server or establish support for every MongoDB driver. Unsupported bindings remain Review.

## SQL rules

- Values interpolated or concatenated into executable SQL are a construction-practice finding.
- Actual query argument roles matter. Binding parameters or replacements preserve data separately
  from SQL structure; `model` or `plain` options do not protect an already-built SQL string.
- Interpolated ORM search values are not automatically executable SQL. Verify the library slot.
- A string merely constructed or logged without a database call is not an executed SQL sink.
- A vulnerability verdict requires request-controlled syntax on a supported path, independently of
  the construction-practice answer. Missing field validation alone does not prove SQL injection.

## NoSQL rules

- Analyze selector keys/values, update data, and executable expressions as different slots.
- An unconstrained request object in a selector value can introduce operators when the verified
  matcher supports them. Do not assume every structured document or update value is unsafe.
- MarsDB `$where` string values are executable JavaScript predicates. Untrusted expression text
  can change behavior; truncating text or stripping newlines does not exclude expression operators.
- A positively constrained scalar selector value can be protected from operator injection.
- JSON/body parsing is transport handling. It does not establish a per-field scalar schema.

## Conditional values and conversions

Analysis MUST evaluate each output alternative separately. A conversion's guarantee applies only
to that output. For example, these actual calibration-control forms have different constraints:

```typescript
// No raw-text alternative: injection protection can hold; construction practice remains.
const id = req.query.mode === 'numeric' ? Number(req.query.id) : 0;
store.find({ $where: 'this.id == ' + id });

// The unconverted alternative still permits request text in the expression.
const id = req.query.mode === 'numeric' ? Number(req.query.id) : req.query.id;
store.find({ $where: 'this.id == ' + id });
```

The scanner verifies the unshadowed native `Number` binding and records bounded per-alternative
facts, including nested ternaries. It does not infer native semantics for a locally redefined
`Number`. Numeric output may include NaN/Infinity or conversion may throw; this does not establish
business validity. Unmodeled alternatives MUST remain explicit rather than inheriting a sibling's
constraint. The current numeric fact adapter recognizes numeric literals and native Number calls;
it does not claim arbitrary arithmetic/alias/heap reasoning.

Protection requires all relevant feasible alternatives to be protected. A vulnerability can be
scoped to one supported unsafe alternative, with its configuration precondition stated. Supplied
source proving a branch impossible excludes it; an unknown necessary path fact remains Inconclusive.
Do not require runtime exploitation to judge a source-supported path, or assume every configuration
is vulnerable merely because one unsafe path exists.

## Outcomes

- **Finding / soft red:** an injection-specific vulnerable judgment, required evidence and API
  contract, no claim blockers, and consistent path/coverage answers support a static vulnerability.
- **Finding / soft orange:** established executable-string construction practice to remediate.
- **Review / soft blue:** the problem or sufficient protection is Inconclusive.
- **Protected / soft green:** protection supported for the named injection claim.

Poor practice can remain orange alongside a protected injection facet. A definitive injection
finding upgrades the same operation record rather than duplicating its construction concern.
Input-validation status does not inherit the injection verdict. Raw answers stay in JSON.
These bounded judgments still require calibration: consistent model answers can be wrong.

## Library contracts

[Sequelize raw queries](https://sequelize.org/docs/v6/core-concepts/raw-queries/) describes query
arguments and value binding. [MarsDB DocumentMatcher](https://github.com/c58/marsdb/blob/master/lib/DocumentMatcher.js)
is the primary implementation reference for executable predicates and operator matching.

## Python adapter

[Python analysis](../python-support.md) verifies sqlite3 connection/cursor provenance and keeps
SQL text separate from the parameter operand. Percent/f-string/format/concatenation construction
remains a practice concern; bound values do not protect dynamic identifiers. Unknown pre-use calls
on the SQL value block an injection claim; transformations of unrelated bound values do not turn
literal SQL into executable interpolation. Python NoSQL is unimplemented and out of the current corpus.

Verified Python connection/cursor bindings also include Psycopg 2/3, PyMySQL and MySQL Connector.
Driver-specific placeholder contracts are supplied to the judge; a parameter operand does not
repair SQL text interpolated earlier. SQLite storage correlation is not extended to those drivers.
Contract references: [Psycopg 3](https://www.psycopg.org/psycopg3/docs/basic/params.html),
[Psycopg 2](https://www.psycopg.org/docs/usage.html#passing-parameters-to-sql-queries),
[PyMySQL](https://pymysql.readthedocs.io/en/latest/modules/cursors.html), and
[MySQL Connector](https://dev.mysql.com/doc/connector-python/en/connector-python-api-mysqlcursor-execute.html).

Generated Python handlers can now retain source-derived registry selection, captured request values,
local callback arguments and helper-returned SQLite connection/cursor provenance. A callback call by
itself is not value flow: the report includes the captured input and actual cursor/SQL-slot witnesses.
Unresolved or mutated registry/callback bindings remain gaps; schema setup and constant SELECT text
are not request-derived SQL findings. See the [bounded Python subset](../python-support.md#finite-generated-http-routes).
