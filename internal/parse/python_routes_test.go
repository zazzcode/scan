package parse

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func pythonFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, text := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestPythonRoutes(t *testing.T) {
	root := pythonFixture(t, map[string]string{"app.py": `from flask import Flask, Blueprint
from fastapi import FastAPI
from flask_restx import Api, Resource
app = Flask(__name__)
bp = Blueprint('bp', __name__, url_prefix='/v1')
app.register_blueprint(bp)
a = FastAPI()
api = Api(app)
ns = api.namespace('api')
@app.route('/hello/<name>', methods=['GET', 'POST'])
def hello(name):
    return name
@bp.get('/status')
def status():
    return 'ok'
@a.post('/items/{item}')
async def item(item: str):
    return item
@ns.route('/entry')
class Entry(Resource):
    def get(self):
        return {}
    def put(self):
        return {}
`})
	g, err := Routes(root)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"GET /hello/<name>": "flask", "POST /hello/<name>": "flask", "GET /v1/status": "flask", "POST /items/{item}": "fastapi", "GET /api/entry": "flask_restx", "PUT /api/entry": "flask_restx"}
	if len(g.Routes) != len(want) {
		t.Fatalf("routes: %#v", g.Routes)
	}
	for _, r := range g.Routes {
		if want[r.Method+" "+r.PathPattern] != r.Framework {
			t.Errorf("unexpected route %#v", r)
		}
	}
}

func TestPythonRouteGolden(t *testing.T) {
	g, err := Routes("testdata/python/routes")
	if err != nil {
		t.Fatal(err)
	}
	g.Root = "python/routes"
	encoded := mustJSON(t, g)
	const golden = "testdata/python/routes/routes.graph.json"
	if *update {
		if err := os.WriteFile(golden, encoded, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, want) {
		t.Fatal("Python route graph differs; run -update and review")
	}
}
