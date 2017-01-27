package guruweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/godoc"
	"guruweb/output"
	"guruweb/web"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type index struct {
	identifiers map[string]token.Position
	files       []string
}

var defaultIndex = &index{}

func InitIndex() error {
	if defaultConfig.verbose {
		output.Trace("init index")
	}
	defaultIndex.identifiers = make(map[string]token.Position)
	defaultIndex.files = make([]string, 0)
	prog, err := defaultConfig.loader.Load()
	if err != nil {
		return err
	}

	for pkg, info := range prog.AllPackages {
		defaultConfig.loadedPkgs = append(defaultConfig.loadedPkgs, pkg.Path())
		impPath := pkg.Path() + "."
		for _, astFile := range info.Files {
			file := prog.Fset.File(astFile.Pos())
			filename := file.Name()
			defaultIndex.files = append(defaultIndex.files, filename)
			for _, decl := range astFile.Decls {
				position := file.Position(decl.Pos())
				if fun, ok := decl.(*ast.FuncDecl); ok {
					defaultIndex.identifiers[impPath+fun.Name.String()] = position
					continue
				}

				if gen, ok := decl.(*ast.GenDecl); ok {
					if gen.Tok == token.TYPE {
						typeSpec, ok := gen.Specs[0].(*ast.TypeSpec)
						if ok {
							_, ok := typeSpec.Type.(*ast.StructType)
							if ok {
								defaultIndex.identifiers[impPath+typeSpec.Name.String()] = position
								continue
							}
							_, ok = typeSpec.Type.(*ast.InterfaceType)
							if ok {
								defaultIndex.identifiers[impPath+typeSpec.Name.String()] = position
							}
						}
					}
				}
			}
		}
	}
	sort.Sort(sort.StringSlice(defaultIndex.files))

	return nil
}

// search return matched files, struct/interface type
func (index *index) search(query string) (map[string]token.Position, []string) {
	findIdents := make(map[string]token.Position)
	findFiles := make([]string, 0)
	for ident, pos := range index.identifiers {
		if strings.Index(ident, query) != -1 {
			findIdents[ident] = pos
		}
	}
	for _, file := range index.files {
		if strings.Index(file, query) != -1 {
			findFiles = append(findFiles, file)
		}
	}

	return findIdents, findFiles
}

// isForbiddenPath validate whether the path in analysis files.
func (index *index) isForbiddenPath(path string) bool {
	i := sort.SearchStrings(defaultIndex.files, path)
	return i >= len(defaultIndex.files) || defaultIndex.files[i] != path
}

func ServeIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeContent(w, req, "", time.Time{}, strings.NewReader(static.Files["index.html"]))
}

func ServeRecommendSearch(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeRecommendSearch: %s", req.URL)
	}
	w.Header().Set("Content-Type", "application/json")
	query := req.FormValue("query")
	if query == "" {
		response, err :=
			json.Marshal(map[string]interface{}{"idents": map[string]token.Position{}, "files": []string{}})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(response)
	} else {
		idents, files := defaultIndex.search(query)
		response, err :=
			json.Marshal(map[string]interface{}{"idents": idents, "files": files})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(response)
	}
}

func ServeRecommendPkgs(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeRecommendPkgs: %s", req.URL)
	}
	w.Header().Set("Content-Type", "application/json")
	query := req.FormValue("query")
	recommendPkgs := make([]string, 0)
	if query != "" {
		recommendPkgs = make([]string, 0)
		for _, pkg := range defaultConfig.allPkgs {
			if strings.Index(pkg, query) != -1 {
				recommendPkgs = append(recommendPkgs, pkg)
			}
		}
	}
	response, err :=
		json.Marshal(recommendPkgs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(response)
}

func ServeQuery(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeQuery: %s", req.URL)
	}
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 500)
		}
	}()
	req.ParseForm()
	mode := req.FormValue("mode")
	pos := req.FormValue("pos")
	accept := req.Header.Get("Accept")
	if mode == "" || pos == "" {
		http.Error(w, "invalid params", 400)
		return
	}
	json := false
	if accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json = true
	}
	response, err := queryGuru(mode, pos, json)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(response)
}

func ServeFile(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeFile: %s", req.URL)
	}
	path := req.FormValue("path")
	if defaultIndex.isForbiddenPath(path) {
		http.Error(w, "Forbidden", 403)
		return
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(req.RemoteAddr, err)
		http.NotFound(w, req)
		return
	}

	var sel godoc.Selection
	s, err := parseSelection(req.FormValue("s"))
	if err == nil {
		offsets := s.byteOffsetsIn(content)
		sel = godoc.RangeSelection(offsets)
	}

	var buf bytes.Buffer
	godoc.FormatText(&buf, content, -1, true, "", sel)

	buf.WriteTo(w)
}

func ServeConfig(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeConfig: %s", req.URL)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if req.Method == "POST" {
		scope := req.FormValue("scope")
		verbose := req.FormValue("verbose")
		v := false
		if verbose == "true" {
			v = true
		}
		if scope != "" {
			InitConfig(strings.Split(scope, ","), v)
			InitIndex()
		} else {
			http.Error(w, "scope should not be empty.", 400)
			return
		}
	}
	scopeJson, err := json.Marshal(defaultConfig)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(scopeJson)
}

func ServeStatic(w http.ResponseWriter, req *http.Request) {
	if defaultConfig.verbose {
		output.Trace("ServeStatic: %s", req.URL)
	}
	name := req.URL.Path
	data, ok := static.Files[name]
	if !ok {
		http.NotFound(w, req)
		return
	}
	http.ServeContent(w, req, name, time.Time{}, strings.NewReader(data))
}
