package parse

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/zazzcode/scan/internal/graph"
)

func TestCheckedRouteContext(t *testing.T) {
	for _, variant := range []struct {
		name, guard, model string
		unknown            bool
	}{
		{"vulnerable", `if(req.body.repeat!==undefined){if(req.body.handle.length===0){res.status(400).send('bad')}}next()`, `handle:{type:DataTypes.STRING}`, false},
		{"protected", `if(!acceptable(req.body.handle)){return res.status(400).send('bad')}next()`, `handle:{type:DataTypes.STRING,validate:{notEmpty:true}}`, false},
		{"unknown", `unknownGuard(req);next()`, `handle:{type:DataTypes.STRING}`, true},
	} {
		t.Run(variant.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, "model.ts", `import {Model,DataTypes} from 'sequelize';export class Account extends Model{};Account.init({`+variant.model+`},{sequelize:db})`)
			writeSource(t, root, "middleware.ts", `export const audit=()=> (req,res,next)=>next()`)
			writeSource(t, root, "server.ts", `import express from 'express';import * as rest from 'finale-rest';import {Account} from './model';import * as audit from './middleware';const api=express();function acceptable(value){return typeof value==='string'&&value.trim().length>0}api.post('/accounts',(req,res,next)=>{`+variant.guard+`});api.post('/accounts',audit.audit());const tables=[{name:'Account',model:Account}];for(const {name,model} of tables){rest.resource({model,endpoints:[`+"`/api/${name}s`"+`, '/unused']})}rest.resource({model:Account,endpoints:['/accounts','/accounts/:id']})`)
			check := graph.Check{ID: "account", Method: "POST", Path: "/accounts", Kind: "account_creation", Fields: []string{"handle"}, Validation: "handle must be a nonempty string", Protection: "reject invalid handles before persistence"}
			report, err := CheckedStrings(root, "server.ts", []graph.Check{check})
			if err != nil {
				t.Fatal(err)
			}
			if len(report.Uses) != 1 {
				t.Fatalf("uses=%+v", report.Uses)
			}
			use := report.Uses[0]
			if len(use.Completeness.Evidence["model_constraints"]) == 0 || len(use.Completeness.Evidence["persistence_mapping"]) == 0 || len(use.ControlEdges) == 0 {
				t.Fatalf("incomplete route/model/control context: %+v", use.Completeness)
			}
			found := false
			for _, g := range use.Completeness.Blockers {
				if g.Kind == "input_helper_unresolved" {
					found = true
				}
			}
			if found != variant.unknown {
				t.Fatalf("unknown guard blocker=%v want %v", found, variant.unknown)
			}
			encoded, _ := json.Marshal(use.Context)
			if variant.name == "protected" && (!strings.Contains(string(encoded), "acceptable") || !strings.Contains(string(encoded), "notEmpty")) {
				t.Fatal("custom guard/model protection omitted")
			}
		})
	}
}

func TestSQLClaimBindingAndUnrelatedGaps(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';import {Sequelize} from 'sequelize';const db=new Sequelize('sqlite::memory:');const api=express();api.post('/session',(req,res)=>{db.query(`+"`SELECT '${req.body.handle}'`"+`,{model:UnknownResult})})`)
	check := graph.Check{ID: "sql", Method: "POST", Path: "/session", Kind: "sql_query", Fields: []string{"handle"}, Validation: "valid handle", Protection: "handle cannot change SQL syntax"}
	report, err := CheckedStrings(root, "server.ts", []graph.Check{check})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Uses) != 1 {
		t.Fatalf("uses=%d", len(report.Uses))
	}
	c := report.Uses[0].Completeness
	binding, _ := json.Marshal(c.Evidence["sink_binding"])
	if !strings.Contains(string(binding), "import {Sequelize} from 'sequelize'") || !strings.Contains(string(binding), "db=new Sequelize") {
		t.Fatal("verified library import/receiver assignment omitted from judge evidence")
	}
	if len(c.Evidence["sink_binding"]) == 0 || len(c.Blockers) != 0 || len(c.Coverage) == 0 {
		t.Fatalf("claim completeness=%+v", c)
	}
}

func TestUnknownMiddlewareBlocksCheckedClaim(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';const api=express();api.use(unknownGuard());api.post('/x',(req,res)=>db.query(req.body.handle))`)
	c := graph.Check{ID: "x", Method: "POST", Path: "/x", Kind: "sql_query", Fields: []string{"handle"}}
	report, err := CheckedStrings(root, "server.ts", []graph.Check{c})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Uses) != 1 || len(report.Uses[0].Completeness.Blockers) == 0 {
		t.Fatal("unknown applicable guard was ignored")
	}
}

func TestMethodSpecificChecks(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';const api=express();api.get('/same',(req,res)=>db.query(req.query.handle));api.post('/same',(req,res)=>db.query(req.body.handle))`)
	report, err := CheckedStrings(root, "server.ts", []graph.Check{{ID: "post", Method: "POST", Path: "/same", Kind: "sql_query", Fields: []string{"handle"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Uses) != 1 || !strings.Contains(report.Uses[0].Sink.Text, "req.body") {
		t.Fatalf("method filtering failed: %+v", report.Uses)
	}
}

func TestRelevantHelperBounds(t *testing.T) {
	for _, tt := range []struct {
		name, call string
		limited    bool
	}{{"unrelated bookkeeping", "bookkeeping()", false}, {"checked-input chain", "level1(req.body.handle)", true}} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeSource(t, root, "model.ts", `export class Account{};`)
			writeSource(t, root, "server.ts", `import express from 'express';import * as rest from 'finale-rest';import {Account} from './model';const api=express();function bookkeeping(){level1('constant')}function level1(v){level2(v)}function level2(v){level3(v)}function level3(v){level4(v)}function level4(v){return v}api.post('/accounts',(req,res,next)=>{`+tt.call+`;next()});rest.resource({model:Account,endpoints:['/accounts','/accounts/:id']})`)
			c := graph.Check{ID: "bounded", Method: "POST", Path: "/accounts", Kind: "account_creation", Fields: []string{"handle"}}
			report, err := CheckedStrings(root, "server.ts", []graph.Check{c})
			if err != nil {
				t.Fatal(err)
			}
			if len(report.Uses) != 1 {
				t.Fatal("missing check")
			}
			limited := false
			for _, g := range report.Uses[0].Completeness.Blockers {
				if g.Kind == "helper_context_limit" {
					limited = true
				}
			}
			if limited != tt.limited {
				t.Fatalf("helper limit=%v want=%v", limited, tt.limited)
			}
		})
	}
}

func TestCheckedEvidenceGolden(t *testing.T) {
	const root = "testdata/string-flow"
	c := graph.Check{ID: "login-input", Method: "POST", Path: "/login", Kind: "sql_query", Fields: []string{"email"}, Validation: "email must be a valid string", Protection: "separate data from SQL structure"}
	report, err := CheckedStrings(root, "server.ts", []graph.Check{c})
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	const file = root + "/checked.graph.json"
	if *update {
		if err := os.WriteFile(file, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, want) {
		t.Fatalf("checked evidence differs; use -update and review diff")
	}
}

func TestShadowedValidationHelperRemainsUnknown(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "model.ts", `export class Account{}`)
	writeSource(t, root, "server.ts", `import express from 'express';import * as rest from 'finale-rest';import {Account} from './model';const api=express();function validate(v){return true}api.post('/accounts',(req,res,next)=>{const validate=unknown;validate(req);next()});rest.resource({model:Account,endpoints:['/accounts','/accounts/:id']})`)
	report, err := CheckedStrings(root, "server.ts", []graph.Check{{ID: "account", Method: "POST", Path: "/accounts", Kind: "account_creation", Fields: []string{"handle"}}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, g := range report.Uses[0].Completeness.Blockers {
		if g.Kind == "input_helper_unresolved" {
			found = true
		}
	}
	if !found {
		t.Fatal("shadowed helper was resolved to an unrelated definition")
	}
}

func TestUnreturnedSequelizeDoesNotBindSink(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';import {Sequelize} from 'sequelize';function make(){const unused=new Sequelize('sqlite::memory:');return fake}const db=make();const api=express();api.post('/x',(req,res)=>db.query(req.body.handle))`)
	report, err := CheckedStrings(root, "server.ts", []graph.Check{{ID: "x", Method: "POST", Path: "/x", Kind: "sql_query", Fields: []string{"handle"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Uses) != 1 || len(report.Uses[0].Completeness.Evidence["sink_binding"]) != 0 {
		t.Fatal("unreturned constructor was used as binding proof")
	}
}

func TestWildcardMiddlewareCannotBeSilentlyIgnored(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';const api=express();api.use('*',guard());api.post('/x',(req,res)=>db.query(req.body.handle))`)
	report, err := CheckedStrings(root, "server.ts", []graph.Check{{ID: "x", Method: "POST", Path: "/x", Kind: "sql_query", Fields: []string{"handle"}}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, g := range report.Uses[0].Completeness.Blockers {
		if g.Kind == "middleware_mount_unresolved" {
			found = true
		}
	}
	if !found {
		t.Fatal("unknown wildcard guard applicability was ignored")
	}
}

func TestConciseSQLHandlerHasCompleteLocalControl(t *testing.T) {
	root := t.TempDir()
	writeSource(t, root, "server.ts", `import express from 'express';import {Sequelize} from 'sequelize';const db=new Sequelize('sqlite::memory:');const api=express();api.post('/x',(req,res)=>db.query('SELECT $handle',{bind:{handle:req.body.handle}}))`)
	report, err := CheckedStrings(root, "server.ts", []graph.Check{{ID: "x", Method: "POST", Path: "/x", Kind: "sql_query", Fields: []string{"handle"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Uses) != 1 || len(report.Uses[0].Completeness.Blockers) != 0 {
		t.Fatalf("implicit-return handler acquired a false control gap: %+v", report.Uses)
	}
}
