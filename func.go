package interfaces

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Func represents an interface function.
type Func struct {
	Name       string `json:"name,omitempty"` // name of the function
	Ins        []Type `json:"ins,omitempty"`  // input parameters
	Outs       []Type `json:"outs,omitempty"` // output parameters
	IsVariadic bool   // whether the function is variadic
}

var variadic = strings.NewReplacer("[]", "...")

// String gives Go code representation of the function.
func (f Func) String() string {
	var buf bytes.Buffer
	if len(f.Ins) == 0 {
		fmt.Fprintf(&buf, "%s()", f.Name)
	} else {
		fmt.Fprintf(&buf, "%s(%s", f.Name, f.in(0))
		for i := range f.Ins[1:] {
			fmt.Fprintf(&buf, ", %s", f.in(i+1))
		}
		buf.WriteString(")")
	}
	if len(f.Outs) == 1 {
		fmt.Fprintf(&buf, " %s", f.Outs[0])
	} else if len(f.Outs) > 1 {
		fmt.Fprintf(&buf, " (%s", f.Outs[0])
		for _, typ := range f.Outs[1:] {
			fmt.Fprintf(&buf, ", %s", typ)
		}
		buf.WriteString(")")
	}
	return buf.String()
}

func (f Func) in(i int) string {
	if typ := f.Ins[i]; i == len(f.Ins)-1 && f.IsVariadic {
		return variadic.Replace(typ.String())
	} else {
		return typ.String()
	}
}

// Deps gives a list of packages the function depends on. E.g. if the function
// represents Serve(net.Listener, http.Handler) error, calling Deps() will
// return []string{"http", "net"}.
//
// The packages are sorted by name.
func (f Func) Deps() []string {
	pkgs := make(map[string]struct{}, 0)
	for _, in := range f.Ins {
		pkgs[in.ImportPath] = struct{}{}
		for _, dep := range in.Deps {
			pkgs[dep] = struct{}{}
		}
	}
	for _, out := range f.Outs {
		pkgs[out.ImportPath] = struct{}{}
		for _, dep := range out.Deps {
			pkgs[dep] = struct{}{}
		}
	}
	delete(pkgs, "")
	if len(pkgs) == 0 {
		return nil
	}
	deps := make([]string, 0, len(pkgs))
	for pkg := range pkgs {
		deps = append(deps, pkg)
	}
	sort.Strings(deps)
	return deps
}

type funcs []Func

func (f funcs) Len() int           { return len(f) }
func (f funcs) Less(i, j int) bool { return f[i].Name < f[j].Name }
func (f funcs) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
