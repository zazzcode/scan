package parse

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

type importBinding struct {
	module, name string
	declaration  *sitter.Node
}
type sourceModule struct {
	file      string
	src       []byte
	hash      string
	tree      *sitter.Tree
	funcs     map[string]*sitter.Node
	exports   map[string]*sitter.Node
	imports   map[string]importBinding
	express   map[string]bool
	receivers map[string]bool
}

type moduleSet struct {
	contextLimits []graph.Gap
	contextFields []string
	root          string
	parsers       *parsers
	modules       map[string]*sourceModule
	allowed       map[string]bool
}

func (s *moduleSet) close() {
	for _, m := range s.modules {
		m.tree.Close()
	}
	s.parsers.close()
}

func (s *moduleSet) load(file string) (*sourceModule, error) {
	file = filepath.ToSlash(filepath.Clean(file))
	if m := s.modules[file]; m != nil {
		return m, nil
	}
	if !s.allowed[file] {
		return nil, fmt.Errorf("module outside selected target sources: %s", file)
	}
	src, err := os.ReadFile(filepath.Join(s.root, filepath.FromSlash(file)))
	if err != nil {
		return nil, err
	}
	if len(src) > 2<<20 {
		return nil, fmt.Errorf("source exceeds 2 MiB limit: %s", file)
	}
	tree := s.parsers.forPath(file).Parse(src, nil)
	if tree == nil {
		return nil, fmt.Errorf("no syntax tree: %s", file)
	}
	m := &sourceModule{file: file, src: src, hash: fmt.Sprintf("%x", sha256.Sum256(src)), tree: tree,
		funcs: map[string]*sitter.Node{}, exports: map[string]*sitter.Node{}, imports: map[string]importBinding{}, express: map[string]bool{}, receivers: map[string]bool{}}
	s.modules[file] = m
	top := tree.RootNode()
	for i := uint(0); i < top.NamedChildCount(); i++ {
		n := top.NamedChild(i)
		if n.Kind() == "import_statement" {
			m.readImport(n)
			continue
		}
		exported := n.Kind() == "export_statement"
		if exported {
			n = n.ChildByFieldName("declaration")
			if n == nil {
				continue
			}
		}
		if n.Kind() == "function_declaration" {
			if name := n.ChildByFieldName("name"); name != nil {
				m.funcs[name.Utf8Text(src)] = n
				if exported {
					m.exports[name.Utf8Text(src)] = n
				}
			}
		}
		if n.Kind() == "lexical_declaration" || n.Kind() == "variable_declaration" {
			for j := uint(0); j < n.NamedChildCount(); j++ {
				d := n.NamedChild(j)
				name := d.ChildByFieldName("name")
				v := d.ChildByFieldName("value")
				if name == nil || v == nil || name.Kind() != "identifier" {
					continue
				}
				key := name.Utf8Text(src)
				if isFunction(v) {
					m.funcs[key] = v
					if exported {
						m.exports[key] = v
					}
				}
				if v.Kind() == "call_expression" {
					fn := v.ChildByFieldName("function")
					if fn != nil && m.express[fn.Utf8Text(src)] {
						m.receivers[key] = true
					}
				}
			}
		}
	}
	return m, nil
}

func (m *sourceModule) readImport(n *sitter.Node) {
	path := n.ChildByFieldName("source")
	module, ok := stringLiteral(path, m.src)
	if !ok {
		return
	}
	for i := uint(0); i < n.NamedChildCount(); i++ {
		clause := n.NamedChild(i)
		if clause.Kind() != "import_clause" {
			continue
		}
		for j := uint(0); j < clause.NamedChildCount(); j++ {
			c := clause.NamedChild(j)
			switch c.Kind() {
			case "identifier":
				m.imports[c.Utf8Text(m.src)] = importBinding{module, "default", n}
				if module == "express" {
					m.express[c.Utf8Text(m.src)] = true
				}
			case "namespace_import":
				if c.NamedChildCount() > 0 {
					m.imports[c.NamedChild(0).Utf8Text(m.src)] = importBinding{module, "*", n}
					if module == "express" {
						m.express[c.NamedChild(0).Utf8Text(m.src)] = true
					}
				}
			case "named_imports":
				for k := uint(0); k < c.NamedChildCount(); k++ {
					spec := c.NamedChild(k)
					name := spec.ChildByFieldName("name")
					alias := spec.ChildByFieldName("alias")
					if name == nil {
						continue
					}
					if alias == nil {
						alias = name
					}
					m.imports[alias.Utf8Text(m.src)] = importBinding{module, name.Utf8Text(m.src), n}
				}
			}
		}
	}
}

func (s *moduleSet) resolve(m *sourceModule, name string) (*sourceModule, *sitter.Node) {
	if parts := strings.Split(name, "."); len(parts) == 2 {
		if imp, ok := m.imports[parts[0]]; ok && imp.name == "*" {
			if target := s.importModule(m, imp.module); target != nil {
				return target, target.exports[parts[1]]
			}
		}
	}

	if n := m.funcs[name]; n != nil {
		return m, n
	}
	imp, ok := m.imports[name]
	if !ok || !strings.HasPrefix(imp.module, ".") {
		return nil, nil
	}
	base := filepath.ToSlash(filepath.Join(filepath.Dir(m.file), imp.module))
	for _, p := range []string{base, base + ".ts", base + ".tsx", base + ".js", base + ".mjs", base + "/index.ts", base + "/index.js"} {
		if !s.allowed[p] {
			continue
		}
		target, err := s.load(p)
		if err != nil {
			return nil, nil
		}
		return target, target.exports[imp.name]
	}
	return nil, nil
}

func (s *moduleSet) importModule(m *sourceModule, path string) *sourceModule {
	if !strings.HasPrefix(path, ".") {
		return nil
	}
	base := filepath.ToSlash(filepath.Join(filepath.Dir(m.file), path))
	for _, p := range []string{base, base + ".ts", base + ".tsx", base + ".js", base + ".mjs", base + "/index.ts", base + "/index.js"} {
		if s.allowed[p] {
			target, err := s.load(p)
			if err == nil {
				return target
			}
			return nil
		}
	}
	return nil
}
