# Node.js Goof XSS answer key

[Sample](README.md) · [Source and hash manifest](../../../.zazz/evaluation/nodejs-goof/xss-cases.json)

This answer key is independent of scanner output. It is based on Goof revision
`add14ba59e98240d9e00a235dd7d42cd61ae9912`, its README, route/model/template source and upstream
`exploits/marked-exploit.sh`. Expected labels describe these specific flows, not the whole app.
Source inspection and upstream exploit documentation support the positive labels; runtime/browser
confirmation has not been performed in this evaluation. Labels are development cases, not an
independently reviewed held-out benchmark. They were documented after the first comparator run.

## Expected vulnerable flows

| ID | Expected property | Source → sink and reason |
|---|---|---|
| `GOOF-XSS-01` | Reflected server XSS | `GET /login` (`app.js:54`) → `req.query.redirectPage` (`routes/index.js:80-85`) → `res.render('admin', …)` → `views/admin.ejs:17`: raw `<%- redirectPage %>` inside a double-quoted `value` attribute. Untrusted quotes/markup can leave the attribute; no encoder or sanitizer is on this path. The upstream README's Open redirect section explicitly documents this XSS separately from the redirect behavior. |
| `GOOF-XSS-02` | Stored server XSS through Markdown, create path | `POST /create` (`app.js:60`) → `req.body.content` (`routes/index.js:168`) → `parse` (`:143-162`, on the non-image branch) → `Todo.content` save (`:185-188`) → `GET /` reads todos (`:35-46`) → `views/index.ejs:20` emits raw `marked(new String(todo.content))`. The upstream documented Markdown bypass has no reminder token ` in ` and no image prefix, so `parse` leaves it unchanged. |
| `GOOF-XSS-03` | Stored server XSS through Markdown, update path | `POST /update/:id` (`app.js:63`) → `req.body.content` assigned to `todo.content` and saved (`routes/index.js:231-239`) → the same `GET /` Markdown render. This write path bypasses `parse` entirely. It shares the rendering weakness with case 02; keep write paths separate without claiming two unrelated sink vulnerabilities. |

For cases 02/03, `app.js:21,75-76` binds the `marked` package and enables `sanitize: true`.
`package.json` and the lockfile pin version `0.3.5`. The upstream README lists this dependency's XSS
weakness and `exploits/marked-exploit.sh:16-20` documents a sanitizer bypass using browser entity/URL
interpretation. Thus these are expected positives supported by upstream reproduction guidance,
not a claim that every Markdown link or every use of `marked` is exploitable. A plain `javascript:`
link and one encoded variant are explicitly labeled failed attempts at `:10-14`; sanitizer option
presence must not turn the known bypass into a Protected answer.

The exact source excerpts and upstream demonstrations remain in the pinned target. Do not execute
exploit scripts as part of static scanning or submit these expected answers to a model.

## Expected negative controls and qualified cases

| ID | Specific claim that should not be an XSS finding | Evidence and limits |
|---|---|---|
| `GOOF-XSS-N01` | Todo content rendered as HTML text by the edit template | `GET /edit/:id` → `routes/index.js:216-228` → `views/edit.ejs:23` uses escaped `<%= todo.content %>`. The same stored value is vulnerable at a different raw Markdown sink; protection belongs to this output only. |
| `GOOF-XSS-N02` | Todo content in the edit form's quoted attribute | `views/edit.ejs:20` uses `<%= todo.content %>` within a double-quoted `value`. Escaped output prevents markup/quote breakout for this context. It is not a URL or script context. |
| `GOOF-XSS-N03` | **Review:** escaped Handlebars fields with request-controlled rendering options | `routes/index.js:102-120` → `views/account.hbs` uses escaped `{{firstname}}` and other field interpolations. Field validation and trimming are not the HTML protection; escaped output is. The complete request object also controls render locals/options; the upstream README documents layout injection on this route. Field syntax is escaped, but the full rendering flow MUST remain Review until those options are resolved. |
| `GOOF-XSS-N04` | Fixed Google Analytics script creation | `public/js/ga.js:3-5` creates a script with literal tag name `script`, assigns a URL from two fixed Google Analytics hosts selected by a protocol equality, and inserts the node. There is no untrusted HTML string or attacker-chosen script destination in this source. Fixed script creation alone is not an XSS flow. HTTP transport/privacy concerns are different claims. |

## Boundaries and verification

The checked-in client JavaScript contains the fixed analytics loader; this inventory establishes
no positive browser-source-to-DOM-sink case. The positive cases above execute in a browser but are
**server-rendered XSS**, not browser DOM-source cases. Use Juice Shop for Angular/DOM positives.
The Dust helper expression in `views/about_new.dust:4` and the upstream Dust exploit script concern
server template evaluation; do not relabel that code-execution demonstration as proven DOM XSS.

Before reusing labels after a target update, check the manifest's file hashes, re-trace each source,
write/read edge, template binding and context, and revise the key when behavior changes. Keep scanner
results and comparison dispositions in the evaluation results document. Neither a scanner finding
nor its absence changes an expected label without source or independent reproduction evidence.

Upstream examples/source are attributed to Snyk Labs contributors under the
[Apache-2.0 license](UPSTREAM-LICENSE.txt).

## Other implemented-family reference cases

`POST /login` at `routes/index.js:51-57` validates username email syntax, then forwards the password
as a MongoDB selector value. The upstream README's NoSQL section documents an operator-object
password bypass against this pinned Mongoose application. This is a source-supported NoSQL test
case, distinct from the XSS cases. The `findById` paths at :204/:232 cast scalar route IDs and are
not automatically equivalent operator-injection positives; keep their schema/casting evidence and
access-control concerns separate. The mounted `/users` TypeORM routes use object operations rather
than request-built raw SQL. Their field validation/persistence needs review; no SQL string injection
is established by the MySQL dependency or the name `repo.find` alone.

`GET /about_new` at :312 logs `JSON.stringify(req.query)`. Whole-query logging is established; its
`device` input is not proven sensitive. This is not a confirmed PII/credential leak without additional
data-class evidence. Likewise the account-details route's documented render-options/layout issue
must not be hidden by its fields' escaped Handlebars syntax.
