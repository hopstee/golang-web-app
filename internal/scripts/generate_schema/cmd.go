package generateschemas

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mobile-backend-boilerplate/internal/kvstore"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// fieldExp = regexp.MustCompile(`@field\s+([\w\[\]\.]+):\s*([\w\[\]]+)(?:\s*=\s*"([^"]+)")?`)
	fieldExp = regexp.MustCompile(`@field\s+([\w\[\]\.]+):\s*([^\s=]+)(?:\s*=\s*"([^"]+)")?(?:\s+dependsOn=([^\s]+))?`)
	// fieldExp = regexp.MustCompile(`@field:\s*([\w\[\]\.]+):\s*([^\s=]+)(?:\s*=\s*"([^"]*)")?`)
)

var Command = &cobra.Command{
	Use:   "generate",
	Short: "Generate pages schema",
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString("version")

		root := "internal/view"

		entities := parseAll(root)

		for _, entity := range entities {
			refs := make([]string, 0)
			if len(entity.Blocks) > 0 {
				for _, b := range entity.Blocks {
					refs = append(refs, b)
					if child, ok := entities[b]; ok {
						entity.Children = append(entity.Children, child)
					}
				}
			}

			if entity.Layout != "" {
				refs = append(refs, entity.Layout)
				if layout, ok := entities[entity.Layout]; ok {
					entity.Children = append(entity.Children, layout)
				}
			}

			if entity.Parent != "" {
				refs = append(refs, entity.Parent)
				if parent, ok := entities[entity.Parent]; ok {
					parent.Children = append(parent.Children, entity)
				}
			}

			entity.Refs = refs
		}

		pages := make([]*kvstore.EntitySchema, 0)
		layouts := make([]*kvstore.EntitySchema, 0)
		blocks := make([]*kvstore.EntitySchema, 0)
		modules := make([]*kvstore.EntitySchema, 0)
		shared := make([]*kvstore.EntitySchema, 0)

		for _, entity := range entities {
			switch entity.Type {
			case "page":
				pages = append(pages, entity)
			case "layout":
				layouts = append(layouts, entity)
			case "block":
				blocks = append(blocks, entity)
			case "module":
				modules = append(modules, entity)
			case "shared":
				shared = append(shared, entity)
			default:
				pages = append(pages, entity)
			}
		}

		if err := saveSchemas("schema", version, entities); err != nil {
			log.Fatal("[error]: error saving schema:", err)
		}
		log.Println("[info]: schema saved successfully")

		if err := saveList("pages", version, pages); err != nil {
			log.Fatal("[error]: error saving pages list:", err)
		}
		log.Println("[info]: pages list saved successfully")

		if err := saveList("layouts", version, layouts); err != nil {
			log.Fatal("[error]: error saving layouts list:", err)
		}
		log.Println("[info]: layouts list saved successfully")

		if err := saveList("blocks", version, blocks); err != nil {
			log.Fatal("[error]: error saving blocks list:", err)
		}
		log.Println("[info]: blocks list saved successfully")

		if err := saveList("modules", version, modules); err != nil {
			log.Fatal("[error]: error saving modules list:", err)
		}
		log.Println("[info]: modules list saved successfully")

		if err := saveList("shared", version, shared); err != nil {
			log.Fatal("[error]: error saving shared list:", err)
		}
		log.Println("[info]: shared list saved successfully")
	},
}

func init() {
	Command.Flags().String("version", "v1", "Schema version")
}

func parseAll(dir string) map[string]*kvstore.EntitySchema {
	schemas := make(map[string]*kvstore.EntitySchema)
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".templ") {
			return nil
		}

		entity := parseFile(path)

		if entity.ID == "" {
			entity.ID = strings.TrimSuffix(info.Name(), ".templ")
		}

		if entity.Type == "" {
			log.Printf("[warning] file %s has no @type directive — skipping\n", path)
			return nil
		}

		schemas[entity.ID] = &entity
		return nil
	})
	return schemas
}

func parseFile(path string) kvstore.EntitySchema {
	file, err := os.Open(path)
	if err != nil {
		log.Println("[error] error open file:", path, err)
		return kvstore.EntitySchema{}
	}
	defer file.Close()

	entity := kvstore.EntitySchema{}
	entity.Blocks = make([]string, 0)
	entity.Refs = make([]string, 0)
	entity.Children = make([]*kvstore.EntitySchema, 0)

	rootSchema := &kvstore.Schema{Fields: []kvstore.Field{}}
	scanner := bufio.NewScanner(file)

	moduleMode := "standalone"

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.Contains(line, "@type:"):
			entity.Type = getTagValue(line, "@type:")
		case strings.Contains(line, "@id:"):
			entity.ID = getTagValue(line, "@id:")
		case strings.Contains(line, "@title:"):
			entity.Title = getTagValue(line, "@title:")
		case strings.Contains(line, "@layout:"):
			entity.Layout = getTagValue(line, "@layout:")
		case strings.Contains(line, "@parent:"):
			entity.Parent = getTagValue(line, "@parent:")
		case strings.Contains(line, "@blocks:"):
			entity.Blocks = parseList(getTagValue(line, "@blocks:"))
		case strings.Contains(line, "@shared:"):
			entity.Shared = parseList(getTagValue(line, "@shared:"))
		case strings.Contains(line, "@mode:"):
			moduleMode = getTagValue(line, "@mode:")
			entity.Mode = moduleMode
		case strings.Contains(line, "@field"):
			if m := fieldExp.FindStringSubmatch(line); len(m) > 0 {
				fullPath := m[1]
				fieldType := m[2]
				label := ""
				if len(m) > 3 {
					label = m[3]
				}
				depends := ""
				if len(m) > 4 {
					depends = m[4]
				}
				addNestedField(rootSchema, fullPath, fieldType, label, depends)
			}
		}
	}

	entity.Content = rootSchema.Fields
	return entity
}

func addNestedField(schema *kvstore.Schema, fullPath, fieldType, label, depends string) {
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
			baseType, opts := parseFieldType(fieldType)

			newField := kvstore.Field{
				Name:    part,
				Type:    guessType(isArray, isLeaf, baseType),
				Options: opts,
				Label:   label,
			}

			if depends != "" && isLeaf {
				newField.Depends = parseDepends(depends)
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

func parseFieldType(rawType string) (baseType string, options []string) {
	rawType = strings.TrimSpace(rawType)
	if strings.HasPrefix(rawType, "select[") && strings.HasSuffix(rawType, "]") {
		baseType = "select"
		optsRaw := strings.TrimSuffix(strings.TrimPrefix(rawType, "select["), "]")
		opts := strings.Split(optsRaw, "|")
		for _, opt := range opts {
			opt = strings.TrimSpace(opt)
			if opt != "" {
				options = append(options, opt)
			}
		}
		return baseType, options
	}
	return rawType, nil
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

func parseDepends(depends string) *kvstore.Dependency {
	parts := strings.SplitN(depends, ":", 2)
	if len(parts) != 2 {
		return &kvstore.Dependency{}
	}

	fieldParts := strings.Split(parts[0], ".")
	dependsOnField := fieldParts[len(fieldParts)-1]
	rawValues := strings.Trim(parts[1], "[]")
	values := strings.Split(rawValues, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}

	return &kvstore.Dependency{
		Field:  dependsOnField,
		Values: values,
	}
}

func saveSchemas(name, version string, schemas map[string]*kvstore.EntitySchema) error {
	outDir := "schemas"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}

	list := make([]*kvstore.EntitySchema, 0, len(schemas))
	for _, entity := range schemas {
		list = append(list, entity)
	}

	filename := fmt.Sprintf("%s/%s.%s.json", outDir, name, version)
	data, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func saveList(name, version string, list []*kvstore.EntitySchema) error {
	outDir := "schemas"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.%s.json", outDir, name, version)
	data, err := json.MarshalIndent(list, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
