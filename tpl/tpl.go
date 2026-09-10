package tpl

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// TemplateMgr manages the global template tree and provides methods for
// executing templates and looking up template names by id.
type TemplateMgr struct {
	root     string             // 模板目录路径
	tpl      *template.Template // 全局唯一的模板树
	tplIdMap map[string]string  // 模板id -> 相对路径
}

// NewTemplateMgr creates a new TemplateMgr with the given root directory and tplIdMap.
func NewTemplateMgr(root string, tplIdMap map[string]string) (*TemplateMgr, error) {
	tpl := template.New("root")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// 用相对路径作为模板名,比如 "user/list.html"
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath) // Windows 下把 \ 换成 /
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tpl.New(relPath).Parse(string(content))
		return err
	})
	if err != nil {
		return nil, err
	}

	if tplIdMap == nil {
		tplIdMap = make(map[string]string)
	}
	return &TemplateMgr{
		root:     root,
		tpl:      tpl,
		tplIdMap: tplIdMap,
	}, nil
}

// Template returns the global template tree.
func (m *TemplateMgr) Template() *template.Template {
	return m.tpl
}

// ExecuteTemplate applies the template associated with t that has the given
// name to the specified data object and writes the output to wr.
func (m *TemplateMgr) ExecuteTemplate(wr io.Writer, name string, data any) error {
	return m.tpl.ExecuteTemplate(wr, name, data)
}

// Lookup returns the template with the given name that is associated with t,
// or nil if there is no such template.
func (m *TemplateMgr) Lookup(name string) *template.Template {
	return m.tpl.Lookup(name)
}

// GetTemplateName returns the name of the template with the given tpl id.
func (m *TemplateMgr) GetTemplateName(tplId string) (string, error) {
	name, ok := m.tplIdMap[tplId]
	if !ok {
		return "", errors.New("template name not found")
	}
	return name, nil
}
