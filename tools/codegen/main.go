//go:build tools

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/emicklei/proto"
	"github.com/fatih/color"
)

var (
	//go:embed template/service/service.tmpl
	serviceFile string

	//go:embed template/biz/biz.tmpl
	bizFile string
)

func main() {
	protoFiles := []string{}
	err := filepath.Walk("./api", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".proto" {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, protoFile := range protoFiles {
		srvName := strings.TrimSuffix(filepath.Base(protoFile), ".proto")
		gen(srvName, protoFile)
	}
}

// gen generates a new Go application from the given app name, proto file, and port.
// It creates the application directory structure, parses the proto file for service info,
// generates files from templates using that info, runs `make api` and `wire`,
// and executes the compiled application.
func gen(appName, protoFile string) {
	internalDir := "internal"
	serverDir := internalDir + "/server"
	serviceDir := internalDir + "/service"
	bizDir := internalDir + "/biz"
	dataDir := internalDir + "/data"
	paths := []string{internalDir, serverDir, bizDir, serviceDir, dataDir}
	for _, v := range paths {
		err := os.MkdirAll(v, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
	serviceInfo, supportHttp := parseProto(protoFile)
	data := map[string]interface{}{
		"appName":     appName,
		"serviceInfo": serviceInfo,
		"module":      parseGoModule(),
		"supportHttp": supportHttp,
	}

	if len(serviceInfo.RpcMeths) == 0 {
		return
	}

	mkFile(data, serviceDir+"/"+appName+".go", serviceFile)
	mkFile(data, bizDir+"/"+appName+".go", bizFile)
}

// mkFile generates a new file with the provided data and template text.
// It first parses the template, then executes it with the provided data.
// If the output file has a ".go" extension, it formats the generated code using go/format.
// If the file already exists, it writes the generated code using WriteDecl.
// Otherwise, it writes the generated code to a new file.
// It prints a success message when the file is generated successfully.
func mkFile(data map[string]interface{}, outFile string, text string) error {
	// Define custom template function
	funcs := template.FuncMap{"Title": strings.Title}

	// Parse the template
	tmpl, err := template.New("tmp").Funcs(funcs).Parse(text)
	if err != nil {
		return err
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	// Get the bytes of the generated code
	codes := buf.Bytes()

	// log.Print(string(codes))

	// Format the code if the file extension is ".go"
	if strings.HasSuffix(outFile, ".go") {
		formattedCodes, err := format.Source(codes)
		if err != nil {
			return err
		}
		codes = formattedCodes
	}

	// Check if the file exists
	if fileExists(outFile) {
		WriteDecl(outFile, string(codes))
		return nil
	}

	// Write the code to the file
	err = os.WriteFile(outFile, codes, 0o644)
	if err != nil {
		return err
	}

	// Print success message
	color.Green("generate file [%s] succeed.\n", outFile)

	return nil
}

type ServiceInfo struct {
	Name     string
	Pkgs     []string
	RpcMeths []MethInfo
	PkgPath  string
	PkgName  string
}

type MethInfo struct {
	MethName       string // rpc方法名
	Param          string // rpc方法入参
	Return         string // rpc方法返回值
	Comment        string // rpc方法注释
	StreamsRequest bool   // 是否为流式请求
	StreamsReturns bool   // 是否为流式返回
}

// parseProto parses a Protocol Buffers (protobuf) file and extracts information about the
// services and methods defined in the file. It returns a ServiceInfo struct containing
// the extracted information, as well as a boolean indicating whether the service supports
// HTTP/JSON transcoding.
//
// The function opens the specified protobuf file, parses its contents using the proto
// package, and walks through the parsed definition to extract the relevant information.
// It converts the request and return types as needed, and populates the ServiceInfo
// struct with the extracted data.
func parseProto(protoFile string) (info ServiceInfo, supportHttp bool) {
	// Open the proto file
	reader, err := os.Open(protoFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// Create a parser and parse the proto file
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Walk through the proto definition and extract information
	proto.Walk(definition,
		proto.WithPackage(func(p *proto.Package) {
			info.PkgName = strings.ReplaceAll(p.Name, ".", "")
		}),
		proto.WithService(func(s *proto.Service) {
			info.Name = s.Name
		}),
		proto.WithRPC(func(r *proto.RPC) {
			// Convert the request type and return type
			req, pkg := convertRequest(info.PkgName, r.RequestType)
			returnType := convertReturnType(info.PkgName, r.ReturnsType)

			// Create a MethInfo struct and populate it with the extracted information
			x := MethInfo{
				MethName:       r.Name,
				Param:          req,
				Return:         returnType,
				StreamsRequest: r.StreamsRequest,
				StreamsReturns: r.StreamsReturns,
			}

			// Add the comment to the MethInfo struct if it exists
			if r.Comment != nil {
				if !strings.HasPrefix(r.Comment.Lines[0], r.Name) {
					r.Comment.Lines[0] = r.Name + strings.TrimPrefix(r.Comment.Lines[0], "//")
				}
				x.Comment = strings.Join(r.Comment.Lines, "\n//")

			}

			// Add the MethInfo struct to the RpcMeths slice in the ServiceInfo struct
			info.RpcMeths = append(info.RpcMeths, x)

			// Add the package to the Pkgs slice in the ServiceInfo struct if it doesn't already exist
			if !slices.Contains(info.Pkgs, pkg) && pkg != "" {
				info.Pkgs = append(info.Pkgs, pkg)
			}
		}),
		proto.WithOption(func(o *proto.Option) {
			// Check if the constant source exists
			if o.Constant.Source == "" {
				return
			}

			// Split the constant source and set the package path in the ServiceInfo struct
			x := strings.Split(o.Constant.Source, ";")
			if len(x) > 0 {
				info.PkgPath = x[0]
			} else {
				info.PkgPath = o.Constant.Source
			}
		}),
	)

	return
}

func convertRequest(pkgName, reqType string) (string, string) {
	switch reqType {
	case "google.protobuf.Empty":
		return "emptypb.Empty", "google.golang.org/protobuf/types/known/emptypb"
	case "google.protobuf.Timestamp":
		return "timestamppb.Timestamp", "google.golang.org/protobuf/types/known/timestamppb"
	case "google.protobuf.Duration":
		return "durationpb.Duration", "google.golang.org/protobuf/types/known/durationpb"
	default:
		return pkgName + "." + reqType, ""
	}
}

func convertReturnType(pkgName, returnType string) string {
	switch returnType {
	case "google.protobuf.Empty":
		return "emptypb.Empty"
	case "google.protobuf.Any":
		return "anypb.Any"
	default:
		return pkgName + "." + returnType
	}
}

// fileExists 检查给定的文件名是否存在。
// 如果文件存在返回 true,否则返回 false。
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func parseGoModule() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	x := strings.Split(string(content), "\n")[0]
	return strings.TrimSpace(strings.Split(x, " ")[1])
}

// WriteDecl writes a new function declaration to the specified file.
// If the file does not exist, it will be created. If the file already exists,
// the existing content will be overwritten.
//
// filename: the path to the file to write the declaration to
// decl: the function declaration to write, as a string
func WriteDecl(filename, decl string) {
	// 解析文件
	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.Contains(decl, "package ") {
		decl = "package main\n" + decl
	}

	// 将新函数的源代码解析为语法树
	funcAST, err := decorator.ParseFile(fset, "", decl, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var funcs []*dst.FuncDecl
	for _, v := range funcAST.Decls {
		if f, ok := v.(*dst.FuncDecl); ok {
			funcs = append(funcs, f)
		}
	}

	for _, newFunc := range funcs {

		index, _ := isFunctionExists(file, newFunc.Name.Name)
		if index < 0 {
			file.Decls = append(file.Decls, newFunc)
			fmt.Print(color.GreenString("New function ["))
			color.New(color.FgCyan, color.Bold).Print(newFunc.Name.Name)
			color.Green("] is added to %s.\n", filename)
		} else {
			file.Decls[index].Decorations().After = newFunc.Decs.After
		}
	}
	if err := reWrite(filename, file); err != nil {
		log.Fatal(err)
	}

	return
}

// 检查函数名是否存在
func isFunctionExists(file *dst.File, functionName string) (index int, exist bool) {
	for i, decl := range file.Decls {
		if fn, ok := decl.(*dst.FuncDecl); ok && fn.Name.Name == functionName {
			return i, true
		}
	}
	return -1, false
}

// reWrite writes the contents of the provided dst.File to the specified filename.
// If the file already exists, it will be truncated and overwritten.
// Returns an error if the file cannot be opened or written to.
func reWrite(filename string, file *dst.File) error {
	outputFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return err
	}
	defer outputFile.Close()

	err = decorator.Fprint(outputFile, file)
	if err != nil {
		fmt.Println("Failed to write file:", err)
		return err
	}

	return nil
}
