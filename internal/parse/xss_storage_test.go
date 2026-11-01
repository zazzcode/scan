package parse

import (
	"strings"
	"testing"
)

func TestStoredTemplatePathsAndControls(t *testing.T) {
	tests := []struct {
		name, write, template, version, protection string
		resolved                                   bool
	}{
		{"raw stored field", "new Note({body:req.body.text}).save()", `<% notes.forEach(function(note){ %><%- note.body %><% }); %>`, "0.3.5", "raw", true},
		{"escaped stored field", "new Note({body:req.body.text}).save()", `<% notes.forEach(function(note){ %><%= note.body %><% }); %>`, "0.3.5", "encoded", true},
		{"old Markdown sanitizer", "new Note({body:req.body.text}).save()", `<% notes.forEach(function(note){ %><%- marked(new String(note.body)) %><% }); %>`, "0.3.5", "ineffective", true},
		{"unknown Markdown version", "new Note({body:req.body.text}).save()", `<% notes.forEach(function(note){ %><%- marked(new String(note.body)) %><% }); %>`, "^0.3.5", "unknown", true},
		{"mismatched field", "new Note({other:req.body.text}).save()", `<% notes.forEach(function(note){ %><%- note.body %><% }); %>`, "0.3.5", "raw", false},
		{"unknown write transform", "new Note({body:opaque(req.body.text)}).save()", `<% notes.forEach(function(note){ %><%- note.body %><% }); %>`, "0.3.5", "unknown", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, "app.js", "const express=require('express');const mongoose=require('mongoose');const marked=require('marked');const Note=mongoose.model('Note');const app=express();app.set('view engine','ejs');app.locals.marked=marked;app.post('/write',(req,res)=>{"+tt.write+"});app.get('/read',(req,res)=>{Note.find({}).exec(function(err,notes){res.render('page',{notes})})})")
			writeSource(t, root, "views/page.ejs", tt.template)
			writeSource(t, root, "package.json", `{"dependencies":{"marked":"`+tt.version+`"}}`)
			r, err := Strings(root, "app.js")
			if err != nil {
				t.Fatal(err)
			}
			uses := xssUses(r)
			if len(uses) != 1 || uses[0].XSS.Type != "stored" || uses[0].XSS.OriginResolved != tt.resolved || uses[0].XSS.Protection != tt.protection {
				t.Fatalf("want stored resolved=%v protection=%s; got %+v", tt.resolved, tt.protection, uses)
			}
			if tt.resolved {
				if len(uses[0].XSS.RelatedEntries) == 0 || !strings.Contains(uses[0].XSS.RelatedEntries[0], "/write") {
					t.Fatal("write entry link missing")
				}
				if !strings.Contains(uses[0].Dependencies[0].Source.Text, "req.body.text") {
					t.Fatal("original request source missing")
				}
			}
		})
	}
}
