package parse

import "testing"

func TestBrowserSourceSinkControls(t *testing.T) {
	tests := []struct{ name, source, protection, context string }{
		{"DOM HTML", "const el=document.getElementById('result');el.innerHTML=location.search", "raw", "html_text"},
		{"DOM text", "const el=document.querySelector('#result');el.textContent=location.hash", "encoded", "html_text"},
		{"URL parameters", "const p=new URLSearchParams(location.search);document.body.innerHTML=p.get('q')", "raw", "html_text"},
		{"document write", "document.write(location.hash)", "raw", "html_text"},
		{"HTML fragment", "const r=document.createRange();r.createContextualFragment(location.hash)", "raw", "html_text"},
		{"adjacent HTML", "document.body.insertAdjacentHTML('beforeend',location.hash)", "raw", "html_text"},
		{"URL sink", "const link=document.querySelector('a');link.href=location.hash", "raw", "url"},
		{"string timer", "setTimeout(location.hash,1)", "raw", "script"},
		{"function constructor", "new Function(location.hash)", "raw", "script"},
		{"message event", "const el=document.querySelector('#x');window.addEventListener('message',(message)=>{el.innerHTML=message.data})", "raw", "html_text"},
		{"sanitized HTML", "import DOMPurify from 'dompurify';document.body.innerHTML=DOMPurify.sanitize(location.hash)", "sanitized", "html_text"},
		{"fake document", "const document={write:x=>x};document.write(location.hash)", "", ""},
		{"shadowed location", "function run(location){document.body.innerHTML=location.hash}", "", ""},
		{"fake DOM object", "const el={innerHTML:''};el.innerHTML=location.hash", "", ""},
		{"fixed script creation", "const script=document.createElement('script');script.src='https://example.test/app.js';document.body.appendChild(script)", "", ""},
		{"imported jQuery", "import $ from 'jquery';const el=$('#result');el.html(location.hash)", "raw", "html_text"},
		{"helper-contained DOM sink", "function render(value){document.body.innerHTML=value};render(location.hash)", "raw", "html_text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, "browser.js", tt.source)
			r, err := Strings(root, "")
			if err != nil {
				t.Fatal(err)
			}
			uses := xssUses(r)
			if tt.protection == "" {
				if len(uses) > 0 {
					t.Fatalf("unexpected XSS: %+v", uses)
				}
				return
			}
			if len(uses) != 1 || uses[0].ClientID == "" || uses[0].RouteID != "" || uses[0].XSS.Protection != tt.protection || uses[0].XSS.OutputContext != tt.context {
				t.Fatalf("want client %s/%s: %+v", tt.protection, tt.context, uses)
			}
			if len(r.ClientEntries) != 1 {
				t.Fatal("browser source has no dedicated entry")
			}
		})
	}
}

func TestReactAndVueRenderingControls(t *testing.T) {
	for _, tt := range []struct{ name, file, source, protection string }{
		{"React raw HTML", "browser.jsx", "import React from 'react';const value=location.hash;const View=()=> <div dangerouslySetInnerHTML={{__html:value}} />", "raw"},
		{"React escaped text", "browser.jsx", "import React from 'react';const value=location.hash;const View=()=> <p>{value}</p>", "encoded"},
		{"Vue raw HTML", "view.vue", `<script setup>const value=location.hash</script><template><div v-html="value"></div></template>`, "raw"},
		{"Vue escaped text", "view.vue", `<script setup>const value=location.hash</script><template><div>{{value}}</div></template>`, "encoded"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, tt.file, tt.source)
			r, err := Strings(root, "")
			if err != nil {
				t.Fatal(err)
			}
			uses := xssUses(r)
			if len(uses) != 1 || uses[0].XSS.Protection != tt.protection || len(r.ClientEntries) != 1 {
				t.Fatalf("want %s; got %+v", tt.protection, uses)
			}
		})
	}
}

func TestAngularTrustAndTemplateBindings(t *testing.T) {
	for _, tt := range []struct{ name, expression, template, protection string }{
		{"trust bypass", "this.sanitizer.bypassSecurityTrustHtml(q)", `<div [innerHTML]="value"></div>`, "raw"},
		{"default sanitization", "q", `<div [innerHTML]="value"></div>`, "sanitized"},
		{"missing renderer", "this.sanitizer.bypassSecurityTrustHtml(q)", `<div>constant</div>`, "unknown"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, "component.ts", "import {Component,inject} from '@angular/core';import {ActivatedRoute} from '@angular/router';import {DomSanitizer} from '@angular/platform-browser';@Component({templateUrl:'./page.html'}) export class Page {private route=inject(ActivatedRoute);private sanitizer=inject(DomSanitizer);run(){const q=this.route.snapshot.queryParams.q;this.value="+tt.expression+"}}")
			writeSource(t, root, "page.html", tt.template)
			r, err := Strings(root, "")
			if err != nil {
				t.Fatal(err)
			}
			uses := xssUses(r)
			if len(uses) != 1 || uses[0].XSS.Protection != tt.protection {
				t.Fatalf("want %s; got %+v", tt.protection, uses)
			}
			if tt.protection != "unknown" && uses[0].Sink.Span.File != "page.html" {
				t.Fatal("template witness missing")
			}
		})
	}
}
