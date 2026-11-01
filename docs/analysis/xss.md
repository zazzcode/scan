# Cross-site scripting analysis

[Analysis index](README.md) · [Result interpretation](../interpreting-results.md)

XSS is a source-to-rendering claim in a specific browser output context. Runtime input types and
validation do not replace context-appropriate encoding or HTML sanitization. The same value may
be protected in an escaped edit form and dangerous in a different raw template output.

## Implemented: reflected server output

The JS/TS adapter resolves verified Express receivers, direct handlers, named relative ESM imports
and literal CommonJS requires/exported functions. Shadowed or reassigned bindings are not trusted.
Request fields reaching `res.send` with string HTML, or `res.write`/`res.end` with an established HTML
content type, produce rendering claims. `res.json` and explicit non-HTML content types are negatives.
An unresolved content type stays Review. Response parameter names are taken from the handler,
rather than assumed to be `res`.

Static Express render configuration links render locals to EJS and Handlebars template outputs.
Raw EJS `<%- … %>` and Handlebars `{{{…}}}` output is distinguished from escaped `<%= … %>` and
`{{…}}`. Custom unresolved engines stay Review. Templates are bounded to the target root and their
source locations/hashes appear in the evidence. Template expressions not yet resolved do not become
protected claims. Pug raw/escaped output and verified file-read/replacement/compile/render flows are also modeled.
Unresolved template expressions and configured file origins stay Review.

These contracts follow the official [Express response API](https://expressjs.com/en/api.html#res.send),
[EJS tags](https://ejs.co/index.html), and [Handlebars escaping](https://handlebarsjs.com/guide/security.html).

## Output context and protection

The adapter distinguishes HTML text, quoted non-active attributes, URL attributes, script/style
content, active attributes and unresolved/mixed contexts. HTML escaping establishes protection only
for HTML text and quoted non-active attributes. It does not establish URL scheme policy, JavaScript
encoding, safe CSS or protection in an unquoted attribute.

Verified imports for `entities`/`he` HTML encoders, `escape-html`, `sanitize-html` and DOMPurify
sanitization supply bounded protection contracts for appropriate contexts. Unknown functions,
same-name fake libraries, shadowed imports and custom sanitizer options are not trusted. Every
modeled path must be protected; protection on one branch does not protect a raw sibling branch.

## Stored and browser analysis

Model-field writes/reads are linked for bounded Mongoose and Sequelize shapes. Storage identity,
field, originating request and rendering evidence are retained; unresolved origins remain Review.
The marked 0.3.5 sanitize-option bypass is a versioned unsafe-rendering contract, supported by
[upstream PR 592](https://github.com/markedjs/marked/pull/592). Unknown versions/configurations are
not treated as protected. Each source path must retain its own transformations and blockers.

Browser entries have their own identity, component/source location and sink context. Supported
bounded patterns include location/URL parameters, message/input events, verified DOM nodes and HTML,
text, URL and string-code sinks, imported jQuery, React JSX, Vue single-file script/template output,
and Angular injected sources/trust APIs with template bindings. An unlinked trust bypass is Review.
API response origins and collection/component/dialog aliases are not fully linked by this first pass.

Default Angular protection requires a demonstrated plain-value path, not merely absence of a known
bypass. Unknown trusted-wrapper transformations remain Review. Whole-object request render locals
can control engine options/layout; escaped field syntax alone cannot make that flow Protected.

## Outcomes and limits

- **Finding / orange:** a supported source path reaches verified raw output without appropriate
  protection. Other uncertain paths remain in the evidence.
- **Review / blue:** unknown source, context, origin, options, sanitizer or transformation.
- **Protected / green:** verified output binding, exact context and appropriate protection on every
  modeled path, without relevant unresolved options/transformations. This applies only to that output.
- **Red:** complete consistent `xss-use-v1` answers and source evidence establish a static XSS claim.
  Runtime confirmation and independent review remain separate.

The protocol asks output context, encoding adequacy, sanitizer effectiveness, path protection and
exposure. Provider confidence/probabilities are retained for audit but never establish a verdict.
Unknown middleware and unsupported paths remain visible. No output or finding is not a safety verdict.

HTML has separate Server/API and Browser/client headings. Hide protected results changes the view
only; canonical JSON/SARIF retain all records. Reusable source-reviewed sample references belong in
[the sample baseline documentation](../samples/reference-policy.md), outside detector/model inputs.
The shared facts use language-neutral source/context/protection/storage/entry fields. Python
security extraction remains unimplemented; it will supply the same facts through a separate adapter.
