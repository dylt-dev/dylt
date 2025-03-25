package template

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"

	"github.com/dylt-dev/dylt/lib"
)

func GetServiceTemplate (svcName string) (*Template, error) {
	fsSvcFiles, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	if err != nil { return nil, err }
	fsSvc, err := fs.Sub(fsSvcFiles, svcName)
	if err != nil { return nil, err }
	tmpl := New(svcName)
	_, err = tmpl.AddContentFS(fsSvc)
	if err != nil { return nil, err }

	return tmpl, nil
}

type Template struct {
	*textTemplate.Template
	Name string
}

func New (name string) *Template {
	return &Template{
		Template: textTemplate.New(name),
		Name: name,
	}
}

// Thin wrapper around AddComponentsFS()
func (t *Template) AddComponentsFolder (componentsPath string) (*Template, error) {
	fsComponents := os.DirFS(componentsPath)
	tmpl, err := t.AddComponentsFS(fsComponents)
	if err != nil { return nil, err }

	return tmpl, nil
}


// Templates in a components folder are intended to be reused across multiple pages,
// eg a navigation bar. Components are no different from any other template, except their
// template name is set to their filename with no path elements at all. Eg 'navbar.tmpl'
// will be named 'navbar'. This makes it easy to reference components from other templates.
// No path knowledge is necessary. All it needs is the name.
func (t *Template) AddComponentsFS (fsComponents fs.FS) (*Template, error) {
	err := fs.WalkDir(fsComponents, ".", func (path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			// Get template name + create child template
			name := getComponentTemplateName(path) 
			tmplChild := t.New(name)
			// Get path to file contents + call Parse() on contents
			buf, err := fs.ReadFile(fsComponents, path)
			if err != nil { return err }
			_, err = tmplChild.Parse(string(buf))
			if err != nil { return err }
		}
		return nil
	})
	if err != nil { return nil, err }

	return t, nil
}

// Thin wrapper around AddContentFS()
func (t *Template) AddContentFolder (contentPath string) (*Template, error) {
	fsContent := os.DirFS(contentPath)
	tmpl, err := t.AddContentFS(fsContent)
	if err != nil { return nil, err }

	return tmpl, nil
}

// Templates in a content folder represent typical Web content. Content folder templates
// are named with the full path relative to their content folder. Eg given a content folder at
// /www that contains a file with a full path of /www/admin/home.tmpl, template.AddContentFolder("/www")
// creates a template with a child template at "/admin/home.tmpl"
func (t *Template) AddContentFS (fsContent fs.FS) (*Template, error) {
	err := fs.WalkDir(fsContent, ".", func (path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			name := getContentTemplateName(path) 
			tmplChild := t.New(name)
			buf, err := fs.ReadFile(fsContent, path)
			if err != nil { return err }
			_, err = tmplChild.Parse(string(buf))
			if err != nil { return err }
		}
		return nil
	})
	if err != nil { return nil, err }

	return t, nil
}
// Deprecated: lib/template Templates are meant to be containers for child templates,
// which get explicitly invoked by ExecuteTemplate(). The additional ability to Execute()
// a lib/template Template is not useful and only leads to confusion. Call ExecuteTemplate()
// with the name of a child template instead.
func (t *Template) Execute (w io.Writer, data any) error {
	return t.Template.Execute(w, data)
}

func (t *Template) GetRunScriptPath () string {
	filename := "run.sh"
	path := filepath.Join(string(filepath.Separator), filename)

	return path
} 

func (t *Template) GetUnitFilePath () string {
	filename := fmt.Sprintf("%s.service", t.Name)
	path := filepath.Join(string(filepath.Separator), filename)

	return path
}

func (t *Template) WriteRunScript (w io.Writer, svc *lib.ServiceSpec) error {
	path := t.GetRunScriptPath()
	err := t.ExecuteTemplate(w, path, svc.Data)
	if err != nil { return err }

	return nil
}

func (t *Template) WriteUnitFile (w io.Writer, svc *lib.ServiceSpec) error {
	path := t.GetUnitFilePath()
	err := t.ExecuteTemplate(w, path, svc.Data)
	if err != nil { return err }

	return nil
}

func getComponentTemplateName (path string) string {
	basename := filepath.Base(path)
	filename := strings.TrimSuffix(basename, ".tmpl")
	templateName := filepath.Join("/", filename)

	return templateName
}

func getContentTemplateName (path string) string {
	filename := strings.TrimSuffix(path, ".tmpl")
	templateName := filepath.Join("/", filename)

	return templateName
}