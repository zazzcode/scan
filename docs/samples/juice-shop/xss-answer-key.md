# Juice Shop: source-reviewed XSS reference

[Sample reports](README.md) · [Pinned file hashes](../../../.zazz/evaluation/sample-reference-manifests.json) ·
[Baseline and remediation policy](../reference-policy.md)

Reviewed 2026-11-01 at `f189f521266f42c6bd05884468292000cee852ed` (v20.3.0).
Labels below come from code inspection, not scanner agreement or challenge names. They describe
supported construction/flow properties and configuration preconditions; no runtime exploitation
was performed. Independent human review remains pending. Unchanged baseline code does not need
another paid model assessment merely to rediscover these properties.

## Source-supported vulnerable constructions

| Case | Source evidence and reason | Scope / condition |
|---|---|---|
| `JS-XSS-01`, DOM search | `search-result.component.ts:138-146` reads `ActivatedRoute.snapshot.queryParams.q`, trims it, calls injected `DomSanitizer.bypassSecurityTrustHtml`, assigns `searchValue`; `search-result.component.html:11` binds that field to `[innerHTML]`. | Expected positive: trimming is not HTML protection, and the trust bypass disables Angular's normal sanitization at the linked sink. No API route is needed for this browser-source case. |
| `JS-XSS-02`, track order | `server.ts:645` registers `/rest/track-order/:id`; `routes/trackOrder.ts:15-25` keeps a truncated request ID when enabled and puts it in the fallback `result.data[0].orderId`; the track-order service fetches that endpoint; `track-result.component.ts:49` trusts this field as HTML, and its template `:9` renders `results.orderNo` through `[innerHtml]`. | Expected cross-tier positive on the enabled branch. Truncation limits length, not markup. The disabled branch's allow-list is a different path. JSON alone is not the sink. |
| `JS-XSS-03`, saved login IP | `routes/saveLoginIp.ts` takes `True-Client-IP`, stores `User.lastLoginIp` without sanitization on the enabled branch, and sanitizes only on the other branch. `routes/login.ts` signs user data into the token. `last-login-ip.component.ts:35-40` reads/decode the token and trusts `payload.data.lastLoginIp`; its template `:10` binds `[innerHTML]`. | Expected stored positive when the unsanitized header is saved and subsequently present in the signed login token. JWT encoding/signing preserves the field; it does not HTML-encode it. Socket-IP fallback is a separate source. |
| `JS-XSS-04`, administration email | `models/user.ts:56-67` stores email without the secure sanitizer when enabled; `administration.component.ts:68-80` obtains users and marks `user.email` trusted; its template `:26,113` renders the field with `[innerHTML]`. | Expected stored positive on that write branch. Table/iterator aliases must be linked; matching the word `email` alone is insufficient. |
| `JS-XSS-05`, feedback | `models/feedback.ts:42-51` chooses one-pass sanitization versus recursive `sanitizeSecure`; `administration.component.ts:92` trusts comments and its template `:60` renders them. `about.component.ts:116-126` builds trusted feedback HTML and passes it as gallery `args`; `about.component.html:51` renders `item?.args`. | Trust bypass and raw HTML construction are established. Sanitizer effectiveness is branch/version-sensitive: investigate `lib/insecurity.ts:58-67` and the actual `sanitize-html` dependency before claiming an exploitable payload on every write branch. Do not treat a function named sanitize as either automatically safe or automatically bypassable. |
| `JS-XSS-06`, product description | `models/product.ts:52-61` omits the secure sanitizer when enabled. `search-result.component.ts:113` trusts descriptions. Product data travels through the product component's `showDetail` (`product.component.ts:212-219`) into dialog `productData`; `product-details.component.html:18` renders its description with `[innerHTML]`. | Expected stored positive on the unsanitized write branch. The actual sink is in the details dialog, not the search template. This requires array, component input and dialog-data links. |
| `JS-XSS-07`, profile template construction | `routes/updateUserProfile.ts` writes the request username; `models/user.ts:43-53` has branch-specific setters. `routes/userProfile.ts:28,54-76,90-101` inserts the value into template source before Pug compilation; `views/userProfile.pug:43` contains the placeholder. | Unsafe template construction is established on its enabled branch. Exact transformation/setter behavior and branch conditions must remain attached. Escaping an ordinary Pug value after compilation is a different operation. |
| `JS-XSS-08`, subtitles | `routes/videoHandler.ts:52-72` reads/compiles a Pug file and inserts subtitle content into the generated script element before sending HTML; `getSubsFromFile` reads the configured subtitle file. | Unsafe raw file-to-HTML construction is established. Attacker control/write exposure of that configured file is **unresolved from this inspected route**, so a definitive externally exploitable stored-XSS claim remains Review. |

## Protected outputs and negative controls

| Case | Why it is not a positive XSS claim | Limits |
|---|---|---|
| `JS-XSS-N01`, ordinary interpolation | `product-details.component.html:17,58` renders product names and review messages through Angular `{{…}}`. Angular treats these as text rather than raw HTML. | Protection is for those text outputs; the same component has the raw/trusted description sink above. |
| `JS-XSS-N02`, default one-way HTML binding | A plain string passed to Angular `[innerHTML]` is sanitized unless a trusted wrapper is introduced. | A protected scanner record must prove a plain-value path and absence of an unresolved trusted-wrapper transformation; merely failing to find a bypass is insufficient. Two-way bindings and resource URL contexts are different contracts. |
| `JS-XSS-N03`, translated constants | Examples include `track-result.component.html:80`, which binds a translated constant through `[innerHtml]`. | No request/browser-controlled markup was identified in the selected expression. Do not infer taint from every `innerHtml` attribute. Translation source changes require separate review. |
| `JS-XSS-N04`, generated CAPTCHA SVG | `routes/imageCaptcha.ts:15-27` generates SVG from fixed factory options and returns `captcha.data`; `data-export.component.ts:58` trusts `data.image`, rendered by its template `:29`. | The trust bypass exists, but no attacker-controlled markup source is established in that generation path. It is not a positive request-to-XSS example merely because it uses the same trust API. |

The [Angular security contract](https://angular.dev/best-practices/security) explains default
interpolation/sanitization and trust bypasses. These reviewed distinctions are acceptance references;
the current extractor still leaves several of their aliases, model setters and producer fields
unresolved. An expected positive or negative in this reference MUST NOT be copied into a report as
a verdict without corresponding generated evidence.

Upstream source attribution: [included MIT license](UPSTREAM-LICENSE.txt).
