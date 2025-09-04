package shared

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed format/onlyoffice-docs-formats.json
var rawFormatsData []byte

type Format struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Actions map[string]string `json:"-"`
	Convert map[string]string `json:"-"`
	Mime    []string          `json:"mime"`
}

func (f Format) IsLossyEditable() bool {
	_, exists := f.Actions["lossy-edit"]
	return exists
}

func (f Format) IsEditable() bool {
	_, exists := f.Actions["edit"]
	return exists
}

func (f Format) IsViewable() bool {
	_, exists := f.Actions["view"]
	return exists
}

func (f Format) IsViewOnly() bool {
	_, exists := f.Actions["view"]
	return exists && len(f.Actions) == 1
}

func (f Format) IsFillable() bool {
	_, exists := f.Actions["fill"]
	return exists
}

func (f Format) IsAutoConvertable() bool {
	_, exists := f.Actions["auto-convert"]
	return exists
}

func (f Format) IsOpenXMLConvertable() bool {
	_, word := f.Convert["docx"]
	_, slide := f.Convert["pptx"]
	_, cell := f.Convert["xlsx"]
	return word || slide || cell
}

func (f Format) GetOpenXMLExtension() string {
	if f.Type == "cell" {
		return "xlsx"
	} else if f.Type == "slide" {
		return "pptx"
	} else {
		return "docx"
	}
}

type MapFormatManager struct {
	formats map[string]Format
}

func NewMapFormatManager() (FormatManager, error) {
	var manager MapFormatManager
	var rawFormats []struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Actions []string `json:"actions"`
		Convert []string `json:"convert"`
		Mime    []string `json:"mime"`
	}

	if err := json.Unmarshal(rawFormatsData, &rawFormats); err != nil {
		return manager, err
	}

	manager.formats = make(map[string]Format)
	for _, rawFormat := range rawFormats {
		actionsSet := make(map[string]string)
		for _, action := range rawFormat.Actions {
			actionsSet[action] = action
		}

		// Excludes unsuppored by the editor formats
		if _, exists := actionsSet["view"]; !exists {
			continue
		}

		convertSet := make(map[string]string)
		for _, conv := range rawFormat.Convert {
			convertSet[conv] = conv
		}

		manager.formats[rawFormat.Name] = Format{
			Name:    rawFormat.Name,
			Type:    rawFormat.Type,
			Actions: actionsSet,
			Convert: convertSet,
			Mime:    rawFormat.Mime,
		}
	}

	return manager, nil
}

type FormatManager interface {
	EscapeFileName(filename string) string
	GetFormatByName(name string) (Format, bool)
	GetAllFormats() map[string]Format
}

func (m MapFormatManager) EscapeFileName(filename string) string {
	f := strings.ReplaceAll(filename, "\\", ":")
	f = strings.ReplaceAll(f, "/", ":")
	return f
}

func (m MapFormatManager) GetFormatByName(name string) (Format, bool) {
	format, exists := m.formats[name]
	return format, exists
}

func (m MapFormatManager) GetAllFormats() map[string]Format {
	return m.formats
}
