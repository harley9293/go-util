package random

import (
	"git.unlimityun.com/lib/go-util/io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const classNameSizeMin = 10    // Minimum length of class name.
const classNameSizeMax = 20    // Maximum length of class name.
const classMethodCountMin = 20 // Minimum number of class methods.
const classMethodCountMax = 50 // Maximum number of class methods.

const methodNameSizeMin = 5   // Minimum length of method name.
const methodNameSizeMax = 25  // Maximum length of method name.
const methodParamCountMin = 0 // Minimum number of method parameters.
const methodParamCountMax = 6 // Maximum number of method parameters.

type OCConfig struct {
	RootDir   string // Generate file root directory.
	FileCount int    // Number of files to generate.
	Seed      int64  // Random seed.
}

type ocType int

const (
	ocVoid = iota + 1
	ocInt
	ocFloat
	ocDouble
	ocChar
	ocNSObject
	ocNSString
	ocNSSet
	ocNSArray
	ocNSDictionary
)

type class struct {
	name        string
	methodCount int
	methods     []*method
}

type method struct {
	isStatic   bool
	retType    ocType
	name       string
	paramCount int
	params     []*param
}

type param struct {
	retType   ocType
	argName   string
	paramName string
}

// ObjectC Generate random Objective-C code.
func ObjectC(config *OCConfig) {
	if config.Seed == 0 {
		config.Seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(config.Seed)

	if io.PathExists(config.RootDir) == false {
		err := os.Mkdir(config.RootDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	generate(config)
}

func generate(config *OCConfig) string {
	var classList []*class
	for i := 0; i < config.FileCount; i++ {
		classList = append(classList, createFile(config.RootDir))
	}

	return createManager(config.RootDir, classList)
}

func createManager(rootDir string, cList []*class) string {
	name := String(10) + "Manager"

	// header
	content := ""
	content += "#import <UIKit/UIKit.h>\n\n"
	content += "@interface " + name + " : NSObject"
	content += "\n\n\n"
	content += "-(void) start;"
	content += "\n\n@end\n"
	ioutil.WriteFile(rootDir+"/"+name+".h", []byte(content), os.ModePerm)

	// mm
	content = ""
	content += "#import \"" + name + ".h\"\n"
	for _, c := range cList {
		content += "#import \"" + c.name + ".h\"\n"
	}
	content += "\n"
	content += "@implementation " + name + " : " + "NSObject\n"
	content += "\n\n\n"

	content += "-(void) start\n{\n"
	for i, c := range cList {
		content += "    " + c.name + " *a" + strconv.Itoa(i) + " = [[" + c.name + " alloc]init];\n"
		content += "    [a" + strconv.Itoa(i) + " start];\n"
	}
	content += "}\n"

	content += "\n\n@end\n"

	ioutil.WriteFile(rootDir+"/"+name+".mm", []byte(content), os.ModePerm)

	return name
}

func createFile(rootDir string) *class {
	c := createClass()
	hFile := createContent(true, c)
	mFile := createContent(false, c)

	ioutil.WriteFile(rootDir+"/"+c.name+".h", []byte(hFile), os.ModePerm)
	ioutil.WriteFile(rootDir+"/"+c.name+".mm", []byte(mFile), os.ModePerm)

	return c
}

func createClass() *class {
	c := new(class)
	c.name = String(UInt(classNameSizeMin, classNameSizeMax))
	c.methodCount = int(UInt(classMethodCountMin, classMethodCountMax))
	for i := 0; i < c.methodCount; i++ {
		c.methods = append(c.methods, createMethod())
	}
	return c
}

func createMethod() *method {
	m := new(method)
	m.name = String(UInt(methodNameSizeMin, methodNameSizeMax))
	m.isStatic = Bool() && Bool() && Bool()
	m.retType = ocType(UInt(ocVoid, ocNSDictionary))
	m.paramCount = int(UInt(methodParamCountMin, methodParamCountMax))
	for i := 0; i < m.paramCount; i++ {
		m.params = append(m.params, createParam(i))
	}
	return m
}

func createParam(index int) *param {
	p := new(param)
	p.argName = "arg" + strconv.Itoa(index)
	p.paramName = "param" + strconv.Itoa(index)
	p.retType = ocType(UInt(ocInt, ocNSDictionary))
	return p
}

func createContent(isHead bool, c *class) string {
	content := ""
	if isHead {
		content += "#import <UIKit/UIKit.h>\n"
	} else {
		content += "#import \"" + c.name + ".h\"\n"
	}
	content += "\n"
	if isHead {
		content += "@interface " + c.name + " : " + "NSObject\n"
	} else {
		content += "@implementation " + c.name + " : " + "NSObject\n"
	}
	content += "\n\n\n"

	content += "-(void) start"
	if isHead {
		content += ";"
	} else {
		content += "\n{\n"
		for i := 0; i < c.methodCount; i++ {
			m := c.methods[i]
			if m.isStatic {
				content += "    [" + c.name + " " + m.name
			} else {
				content += "    [self " + m.name
			}

			if m.paramCount > 0 {
				content += ":"
			}

			for j := 0; j < m.paramCount; j++ {
				if j > 0 {
					content += m.params[j].argName + ":"
				}
				content += getTypeDefaultValue(m.params[j].retType)
				if j < m.paramCount-1 {
					content += " "
				}
			}
			content += "];\n"
		}
		content += "}\n"
	}
	content += "\n\n"

	for i := 0; i < c.methodCount; i++ {
		m := c.methods[i]

		if m.isStatic {
			content += "+"
		} else {
			content += "-"
		}

		content += "(" + getTypeName(m.retType) + ")" + m.name
		if m.paramCount > 0 {
			content += ":"
		}

		for j := 0; j < m.paramCount; j++ {
			if j > 0 {
				content += m.params[j].argName + ":"
			}
			content += "(" + getTypeName(m.params[j].retType) + ")" + m.params[j].paramName
			if j < m.paramCount-1 {
				content += " "
			}
		}

		if isHead {
			content += ";\n\n"
		} else {
			content += "\n{\n"
			content += "    for (int i=0; i<100; i++) {\n"
			content += "        i++;\n"
			content += "    }\n"

			if m.retType != ocVoid {
				content += "    " + getTypeName(m.retType) + " xxxx = " + getTypeDefaultValue(m.retType) + ";\n"
				content += "    return xxxx;\n"
			}
			content += "}\n\n"
		}
	}

	content += "\n\n@end\n"
	return content
}

func getTypeName(t ocType) string {
	switch t {
	case ocVoid:
		return "void"
	case ocChar:
		return "char"
	case ocInt:
		return "int"
	case ocFloat:
		return "float"
	case ocDouble:
		return "double"
	case ocNSObject:
		return "NSObject*"
	case ocNSSet:
		return "NSSet*"
	case ocNSArray:
		return "NSArray*"
	case ocNSString:
		return "NSString*"
	case ocNSDictionary:
		return "NSDictionary*"
	default:
		return "void"
	}
}

func getTypeDefaultValue(t ocType) string {
	switch t {
	case ocChar:
		return "'A'"
	case ocInt:
		return "100"
	case ocFloat:
		return "101.11"
	case ocDouble:
		return "100.0011"
	case ocNSObject:
		return "[NSObject alloc]"
	case ocNSSet:
		return "[NSSet alloc]"
	case ocNSArray:
		return "[NSArray alloc]"
	case ocNSString:
		return "[NSString alloc]"
	case ocNSDictionary:
		return "[NSDictionary alloc]"
	default:
		return ""
	}
}
