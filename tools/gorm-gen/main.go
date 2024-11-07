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
	"regexp"
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
	
	// 解析表名和注释
	if err := parseTableInfo(sql, table); err != nil {
		return nil, fmt.Errorf("parse table info: %w", err)
	}
	
	fields := make([]Field, 0)
	lines := strings.Split(sql, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		switch {
		case shouldSkipLine(line):
			continue
		case strings.HasPrefix(line, ") ENGINE="):
			break
		case strings.HasPrefix(line, "PRIMARY KEY"):
			handlePrimaryKey(line, &fields)
		case strings.HasPrefix(line, "UNIQUE KEY"):
			handleUniqueKey(line, &fields)
		default:
			if field := parseField(line); field != nil {
				if field.Name == "deleted_at" {
					table.Model = "BaseModel"
				}
				if !SliceContain(skipFields, field.Name) {
					fields = append(fields, *field)
				}
			}
		}
	}

	// 生成所有字段的结构体标签
	for i := range fields {
		fields[i].Tag = generateStructTag(fields[i])
	}
	
	table.Fields = fields
	return table, nil
}

func parseTableInfo(sql string, table *Table) error {
	// 解析表名
	tableMatch := regexp.MustCompile(`CREATE TABLE\s+` + "`" + `?([^` + "`" + `\s]+)` + "`" + `?\s`).FindStringSubmatch(sql)
	if len(tableMatch) < 2 {
		return fmt.Errorf("failed to parse table name")
	}
	table.Name = tableMatch[1]
	
	// 解析表注释
	commentMatch := regexp.MustCompile(`COMMENT='([^']*)'`).FindStringSubmatch(sql)
	if len(commentMatch) >= 2 {
		table.Comment = commentMatch[1]
	}
	
	return nil
}

func shouldSkipLine(line string) bool {
	return strings.HasPrefix(line, "CREATE TABLE") ||
		   strings.HasPrefix(line, "KEY") ||
		   strings.HasPrefix(line, ")") ||
		   strings.HasPrefix(line, "CONSTRAINT") ||
		   line == ""
}

func handlePrimaryKey(line string, fields *[]Field) {
	pk := getPrimaryKey(line)
	for i, f := range *fields {
		if SliceContain(pk, f.Name) {
			(*fields)[i].Primary = true
		}
	}
}

func handleUniqueKey(line string, fields *[]Field) {
	idx := getIndex(line)
	for i, f := range *fields {
		if f.Name == idx {
			(*fields)[i].Unique = true
		}
	}
}

func parseField(line string) *Field {
	// 跳过非字段定义行
	if !strings.Contains(line, " ") {
		return nil
	}
	
	field := &Field{}
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	
	field.Name = strings.Trim(parts[0], "`")
	field.Type = getType(parts[1])
	field.Nullable = !strings.Contains(line, "NOT NULL")
	
	// 解析默认值
	if defaultMatch := regexp.MustCompile(`DEFAULT\s+([^ ,]+)`).FindStringSubmatch(line); len(defaultMatch) >= 2 {
		field.Default = pareDefaultValue(field.Type, defaultMatch[1])
	}
	
	// 解析注释
	if commentMatch := regexp.MustCompile(`COMMENT\s+'([^']*)'`).FindStringSubmatch(line); len(commentMatch) >= 2 {
		field.Comment = commentMatch[1]
	}
	
	return field
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
	// 预定义数据库类型到Go类型的映射
	typeMap := map[string]string{
		"bigint":    "int64",
		"int":       "int",
		"tinyint":   "int",
		"smallint":  "int",
		"tinyint(1)": "bool",
		"varchar":   "string", 
		"text":      "string",
		"char":      "string",
		"longtext":  "string",
		"enum":      "string",
		"decimal":   "float64",
		"double":    "float64",
		"float":     "float32",
		"timestamp": "jtime.JsonTime",
		"datetime":  "jtime.JsonTime",
		"json":      "json.RawMessage",
	}

	token = strings.ToLower(token)
	
	// 遍历映射表检查前缀匹配
	for dbType, goType := range typeMap {
		if strings.HasPrefix(token, dbType) {
			return goType
		}
	}

	return token
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
