package generateschemas

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"mobile-backend-boilerplate/internal/kvstore"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// fieldExp = regexp.MustCompile(`@field\s+([\w\[\]\.]+):\s*([\w\[\]]+)(?:\s*=\s*"([^"]+)")?`)
	fieldExp = regexp.MustCompile(`@field\s+([\w\[\]\.]+):\s*([^\s=]+)(?:\s*=\s*"([^"]+)")?`)
	// fieldExp = regexp.MustCompile(`@field:\s*([\w\[\]\.]+):\s*([^\s=]+)(?:\s*=\s*"([^"]*)")?`)
)

var Command = &cobra.Command{
	Use:   "generate",
	Short: "Generate pages schema",
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString("version")

		layoutSchemas := parseDir("internal/view/layouts", "layout")
		blocksSchemas := parseDir("internal/view/blocks", "block")
		pagesSchemas := parseDir("internal/view/pages", "page")
		modulesSchemas := parseDir("internal/view/modules", "module")

		for _, page := range pagesSchemas {
			if page.Layout != "" {
				if layout, ok := layoutSchemas[page.Layout]; ok {
					page.SEO = layout.Content
				}
			}

			for _, blockId := range page.Blocks {
				if block, ok := blocksSchemas[blockId]; ok {
					page.Children = append(page.Children, block)
				}
			}
		}

		for _, block := range blocksSchemas {
			if block.Parent != "" {
				if parent, ok := blocksSchemas[block.Parent]; ok {
					parent.Children = append(parent.Children, block)
				}
			}
		}

		pages := make([]*kvstore.EntitySchema, 0, len(pagesSchemas))
		for _, p := range pagesSchemas {
			pages = append(pages, p)
		}

		output := fmt.Sprintf("internal/schemas/pages.%s.json", version)
		data, _ := json.MarshalIndent(pages, "", "	")
		os.WriteFile(output, data, 0644)
		fmt.Println("Pages schemas generated:", output)

		modules := make([]*kvstore.EntitySchema, 0, len(modulesSchemas))
		for _, m := range modulesSchemas {
			modules = append(modules, m)
		}

		output = fmt.Sprintf("internal/schemas/modules.%s.json", version)
		data, _ = json.MarshalIndent(modules, "", "	")
		os.WriteFile(output, data, 0644)
		fmt.Println("Modules schemas generated:", output)
	},
}

func init() {
	Command.Flags().String("version", "v1", "Schema version")
}

func parseDir(dir, typ string) map[string]*kvstore.EntitySchema {
	schemas := make(map[string]*kvstore.EntitySchema)
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".templ") {
			return nil
		}

		entity := parseFile(path)
		entity.Type = typ

		if entity.ID == "" {
			entity.ID = strings.TrimSuffix(info.Name(), ".templ")
		}

		schemas[entity.ID] = &entity
		return nil
	})
	return schemas
}

func parseFile(path string) kvstore.EntitySchema {
	fmt.Println("File: ", path)
	file, err := os.Open(path)
	if err != nil {
		return kvstore.EntitySchema{}
	}
	defer file.Close()

	entity := kvstore.EntitySchema{}
	scanner := bufio.NewScanner(file)
	rootSchema := &kvstore.Schema{Fields: []kvstore.Field{}}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.Contains(line, "@page:"):
			entity.ID = getTagValue(line, "@page:")
		case strings.Contains(line, "@layout:"):
			entity.Layout = getTagValue(line, "@layout:")
		case strings.Contains(line, "@parent:"):
			entity.Parent = getTagValue(line, "@parent:")
		case strings.Contains(line, "@title:"):
			entity.Title = getTagValue(line, "@title:")
		case strings.Contains(line, "@blocks:"):
			entity.Blocks = parseList(getTagValue(line, "@blocks:"))
		case strings.Contains(line, "@module:"):
			entity.ID = getTagValue(line, "@module:")
		case strings.Contains(line, "@field"):
			if m := fieldExp.FindStringSubmatch(line); len(m) > 0 {
				fullPath := m[1]
				fieldType := m[2]
				label := ""
				if len(m) > 3 {
					label = m[3]
				}
				addNestedField(rootSchema, fullPath, fieldType, label)
			}
		}
	}

	entity.Content = rootSchema.Fields
	fmt.Println(entity)
	return entity
}

func addNestedField(schema *kvstore.Schema, fullPath, fieldType, label string) {
	parts := strings.Split(fullPath, ".")
	current := schema

	for i, rawPart := range parts {
		isArray := strings.HasSuffix(rawPart, "[]")
		part := strings.TrimSuffix(rawPart, "[]")
		isLeaf := i == len(parts)-1

		var existing *kvstore.Field
		for j := range current.Fields {
			if current.Fields[j].Name == part {
				existing = &current.Fields[j]
				break
			}
		}

		if existing == nil {
			newField := kvstore.Field{
				Name:  part,
				Type:  guessType(isArray, isLeaf, fieldType),
				Label: label,
			}
			if !isLeaf {
				newField.Schema = &kvstore.Schema{}
			}
			current.Fields = append(current.Fields, newField)
			existing = &current.Fields[len(current.Fields)-1]
		}

		if !isLeaf {
			if existing.Schema == nil {
				existing.Schema = &kvstore.Schema{}
			}
			current = existing.Schema
		}
	}
}

func getTagValue(line, tag string) string {
	i := strings.Index(line, tag)
	if i == -1 {
		return ""
	}

	val := strings.TrimSpace(line[i+len(tag):])
	// если есть комментарий, то избавляемся от него
	val = strings.Split(val, "#")[0]
	return strings.TrimSpace(val)
}

func parseList(value string) []string {
	items := strings.Split(value, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func guessType(isArray, isLeaf bool, explicit string) string {
	if isLeaf {
		return explicit
	}
	if isArray {
		return "list[object]"
	}
	return "object"
}
