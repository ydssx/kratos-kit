//go:build tools

package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/ydssx/kratos-kit/common/conf"

	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/inflection"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed model.tmpl
var modelTmp string

var (
	outputDir         = "models/"
	defaultConfigPath = "configs/config.local.yaml"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", defaultConfigPath, "path to config file")
	flag.Parse()

	bootStrap := conf.Bootstrap{}
	conf.MustLoad(&bootStrap, configFile)

	for _, v := range bootStrap.Data.Database.Source {
		dsn, _ := mysql.ParseDSN(v)
		db, err := gorm.Open(gmysql.Open(v))
		if err != nil {
			// 处理连接错误
			log.Fatalf("failed to connect database: %v", err)
		}

		tables, _ := db.Migrator().GetTables()
		for _, tableName := range tables {
			var createSQL string
			if err := db.Raw("SHOW CREATE TABLE "+tableName).Row().Scan(&tableName, &createSQL); err != nil {
				// 处理错误
				log.Fatalf("failed to get sql: %v", err)
			}

			generate(dsn.DBName, createSQL, outputDir)
		}
	}
}

func generate(dbName, createSQL, outPath string) {
	table, err := ParseSQL(createSQL)
	if err != nil {
		log.Print("failed to parse sql:", err)
		return
	}

	funcMap := template.FuncMap{"Title": strings.Title, "Lower": toLowerFirst, "CamelCase": UnderscoreToCamelCase}
	// 解析模板
	tmpl, err := template.New("model").Funcs(funcMap).Parse(modelTmp)
	if err != nil {
		fmt.Println("failed to parse template:", err)
		return
	}
	// 将模型转换为模板需要的数据
	data := map[string]interface{}{
		"TableName":    table.Name,
		"TableComment": table.Comment,
		"Fields":       table.Fields,
		"Name":         GetSingularTableName(table.Name),
		"PrimaryKey":   findPrimaryKey(*table),
		"model":        table.Model,
		"dbName":       dbName,
	}
	// 将模板应用到数据上，生成代码
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Println("failed to generate code:", err)
		return
	}

	// 格式化生成的代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("failed to format code:", err)
		return
	}

	// 将生成的代码写入文件
	filename := filepath.Join(outPath, strings.ToLower(table.Name)+".go")
	if FileExist(filename) {
		msg := color.YellowString("file %s already exists, skipped.", filename)
		fmt.Println(msg)
		return
	}
	if err := ioutil.WriteFile(filename, formattedCode, 0o644); err != nil {
		fmt.Println("failed to write code to file:", err)
		return
	}
	s := color.BlueString("[table %s]", table.Name)
	fmt.Printf("%s:code generation succeeded!\n", s)
}

type Table struct {
	Name    string
	Comment string
	Fields  []Field
	Model   string
}

type Field struct {
	Name     string
	Type     string
	Primary  bool
	Unique   bool
	Nullable bool
	Default  interface{}
	Comment  string
	Tag      string
}

var skipFields = []string{"id", "created_at", "updated_at", "deleted_at"}

func ParseSQL(sql string) (*Table, error) {
	table := &Table{Model: "BaseModelNoDelete"}
	fields := []Field{}
	// Extract table name
	tableNameStart := strings.Index(sql, "CREATE TABLE ") + 13
	tableNameEnd := strings.Index(sql[tableNameStart:], " ")
	table.Name = sql[tableNameStart : tableNameStart+tableNameEnd]
	table.Name = strings.ReplaceAll(table.Name, "`", "")

	// Extract table comment
	if strings.Contains(sql, "COMMENT='") {
		commentStart := strings.Index(sql, "COMMENT='") + 9
		commentEnd := strings.Index(sql[commentStart:], "'")
		table.Comment = sql[commentStart : commentStart+commentEnd]
	}

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CREATE TABLE") || strings.HasPrefix(line, "KEY") || strings.HasPrefix(line, ")") || strings.HasPrefix(line, "CONSTRAINT") {
			continue
		} else if strings.HasPrefix(line, ") ENGINE=") {
			break
		} else if strings.HasPrefix(line, "PRIMARY KEY") {
			pk := getPrimaryKey(line)
			for i, f := range fields {
				if SliceContain(pk, f.Name) {
					fields[i].Primary = true
				}
			}
		} else if strings.HasPrefix(line, "UNIQUE KEY") {
			idx := getIndex(line)
			for i, f := range fields {
				if f.Name == idx {
					fields[i].Unique = true
				}
			}
		} else {
			field := getField(line)
			if field.Name == "deleted_at" {
				table.Model = "BaseModel"
			}
			if SliceContain(skipFields, field.Name) {
				continue
			}

			// field.Tag = generateStructTag(field)
			fields = append(fields, field)
		}
	}
	for i, field := range fields {
		fields[i].Tag = generateStructTag(field)
	}
	table.Fields = fields
	return table, nil
}

func getTableName(line string) string {
	tokens := strings.Split(line, " ")
	return strings.TrimSuffix(tokens[2], " (")
}

func getTableComment(line string) string {
	comment := ""
	if strings.Contains(line, "COMMENT") {
		start := strings.Index(line, "COMMENT '") + 9
		end := strings.Index(line[start:], "'") + start
		comment = line[start:end]
	}
	return comment
}

func getPrimaryKey(line string) []string {
	start := strings.Index(line, "(") + 1
	end := strings.Index(line, ")")
	x := line[start:end]
	x = strings.ReplaceAll(x, "`", "")
	return strings.Split(x, ",")
}

func getIndex(line string) string {
	start := strings.Index(line, "(") + 1
	end := strings.Index(line, ")")
	return line[start:end]
}

func getField(line string) Field {
	field := Field{}
	tokens := strings.Split(line, " ")
	field.Name = strings.TrimSuffix(tokens[0], ",")
	field.Name = strings.ReplaceAll(field.Name, "`", "")
	field.Type = getType(tokens[1])
	if strings.Contains(line, "NOT NULL") {
		field.Nullable = false
	} else {
		field.Nullable = true
	}
	if strings.Contains(line, "DEFAULT") {
		start := strings.Index(line, "DEFAULT ") + 8
		fval := strings.TrimRight(strings.Split(line[start:], " ")[0], ",")
		field.Default = pareDefaultValue(field.Type, fval)
	}
	if strings.Contains(line, "COMMENT") {
		start := strings.Index(line, "COMMENT '") + 9
		end := strings.Index(line[start:], "'") + start
		field.Comment = line[start:end]
	}
	return field
}

// 生成模型tag
func generateStructTag(field Field) (tag string) {
	// fieldStr := fmt.Sprintf("%s %s", field.Name, field.Type)
	tags := []string{fmt.Sprintf("column:%s", field.Name)}
	if field.Primary {
		tags = append(tags, "primaryKey")
	}
	if field.Unique {
		tags = append(tags, "unique")
	}
	if !field.Nullable {
		tags = append(tags, "not null")
	}
	if field.Default != "" && field.Default != nil {
		tags = append(tags, fmt.Sprintf("default:%v", field.Default))
	}
	tag += fmt.Sprintf("`json:\"%s\" gorm:\"%s\"`", field.Name, strings.Join(tags, ";"))
	return tag
}

func findPrimaryKey(table Table) string {
	for _, v := range table.Fields {
		if v.Primary {
			return v.Name
		}
	}
	return ""
}

func getType(token string) string {
	token = strings.ToLower(token)
	switch {
	case strings.HasPrefix(token, "bigint"):
		return "int64"
	case strings.HasPrefix(token, "int"), strings.HasPrefix(token, "tinyint"), strings.HasPrefix(token, "smallint"):
		return "int"
	case strings.HasPrefix(token, "tinyint(1)"):
		return "bool"
	case strings.HasPrefix(token, "varchar"), strings.HasPrefix(token, "text"), strings.HasPrefix(token, "char"), strings.HasPrefix(token, "longtext"), strings.HasPrefix(token, "enum"):
		return "string"
	case strings.HasPrefix(token, "decimal"), strings.HasPrefix(token, "double"):
		return "float64"
	case strings.HasPrefix(token, "float"):
		return "float32"
	case strings.HasPrefix(token, "timestamp"), strings.HasPrefix(token, "datetime"):
		return "time.Time"
	case strings.HasPrefix(token, "json"):
		return "json.RawMessage"
	default:
		return token
	}
}

func pareDefaultValue(ftype, fval string) (v interface{}) {
	if strings.ToLower(fval) == "null" {
		return fval
	}
	switch ftype {
	case "int64", "int", "int32":
		v, _ = strconv.ParseInt(strings.Trim(fval, "'"), 10, 64)
	case "float64", "float32":
		vf, _ := strconv.ParseFloat(strings.Trim(fval, "'"), 64)
		v = math.Round(vf*100) / 100
	default:
		return fval
	}
	return
}

func GetSingularTableName(tableName string) string {
	// 1. 先将表名转换为驼峰式命名
	tableName = strings.ReplaceAll(tableName, "_", " ")
	tableName = strings.Title(tableName)
	tableName = strings.ReplaceAll(tableName, " ", "")

	// 2. 使用 inflection 库将驼峰式命名的表名转换为单数形式
	singularName := inflection.Singular(tableName)

	return singularName
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

// 下划线转驼峰
func UnderscoreToCamelCase(s string) string {
	var (
		b  strings.Builder
		up bool
	)

	for _, c := range s {
		if c == '_' {
			up = true
			continue
		}

		if up {
			b.WriteRune(unicode.ToUpper(c))
			up = false
		} else {
			b.WriteRune(c)
		}
	}

	return b.String()
}

func SliceContain(s []string, elem string) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
}

func DirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}
