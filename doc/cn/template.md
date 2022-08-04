# 前言
`text/template`是Go语言标准库，实现数据驱动模板以生成文本输出，可以理解为一组文字按照特定格式动态嵌入另一组文字中。

还有个处理html文字的模板（`html/template`），html/template包实现了数据驱动的模板，用于生成可对抗代码注入的安全HTML输出。它提供了和`text/template`包相同的接口，Go语言中输出HTML的场景都应使用`text/template`包。

简单区分了两个长得几乎一模一样的包之后，本章节主要开始对 `text/template` 包的介绍。
# 介绍
## 模板标签
```
{{.}}
```
模板语法都包含在`{{`和`}}`中间，其中`{{.}}`中的点表示当前对象。

## 注释
```
{{/* a comment */}}
```
使用`“{{/*”和“*/}}”`来包含注释内容

## go语言示例

当我们传入一个结构体对象时，我们可以根据.来访问结构体的对应字段。例如：

```
package main

import (
	"os"
	"text/template"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	name := "world"

	// 建立一个叫test的模板，内容是 “hello,{{.}}”
	tmpl, err := template.New("test").Parse("hello,{{.}}")
	CheckErr(err)

	// 将string与模板合成，变量name的内容会替换掉 {{.}}
	// 合成结果放到os.Stdout输出
	err = tmpl.Execute(os.Stdout, name)

	CheckErr(err)
}

// Stdin、Stdout 和 Stderr 是指向标准输入、标准输出和标准错误文件描述符的打开文件。
// 请注意，Go 运行时会写入死机和崩溃的标准错误;关闭 Stderr 可能会导致这些消息转到其他位置，可能是稍后打开的文件。

```

输出：
```
hello, world
```
模板的输入文本是任何格式的`UTF-8`编码文本。`{{ 和 }}` 包裹的内容统称为 `action`，分为两种类型：

* 数据求值（data evaluations）
* 控制结构（control structures）

**action 求值的结果会直接复制到模板**中，控制结构和我们写 Go 程序差不多，也是条件语句、循环语句、变量、函数调用等等…

将模板成功解析（Parse）后，可以安全地在并发环境中使用，如果输出到同一个 `io.Writer` 数据可能会重叠（因为不能保证并发执行的先后顺序）。

**这里 {{ 和 }} 中间的句号（.）代表传入的数据，数据不同渲染不同**，可以代表 go 语言中的任何类型，如结构体、哈希等。

本篇文章我们重点关注json文件中template定义和格式，如何对template在go语言中的解析以及定义通过上述`hello，world`例子有个大概的概念即可。

# 变量
在golang渲染template的时候，可以接受一个`interface{}`类型的变量，我们在模板文件中可以读取变量内的值并渲染到模板里。

有两个常用的传入参数的类型。一个是struct，在模板内可以读取该struct域的内容来进行渲染。还有一个是`map[string]interface{}`，在模板内可以使用key来进行渲染。

我一般使用第二种，效率可能会差一点儿，但是用着方便。

模板内内嵌的语法支持，全部需要加`{{}}`来标记。

在模板文件内， `. `代表了当前变量，即在非循环体内，`.`就代表了传入的那个变量。假设我们在go代码中定义了一个结构体`type Article struct`：

```
type c struct {
	... ...
}

type Article struct {
    a int
    b string
    c
}
```
在模板内可以通过如下形式调用结构体。
```
{{.a}}

{{.b}}
```
* `{{.a}}`表示输出struct对象中字段`“a”`的值

* `{{.b}}`表示输出struct对象中字段`“b”`的值

当`“c”`是匿名字段（没有名字的函数）时，可以访问其内部字段或方法，比如
```
"Com"：{{.c.Com}} 
```
如果“Com”是一个方法并返回一个`Struct`对象，同样也可以访问其字段或方法：`{{.c.Com.Field1}}`
```
{{.Method1 "参数值1" "参数值2"}}
```
调用方法“Method1”，将后面的参数值依次传递给此方法，并输出其返回值。
```
{{$c}}
```
此标签用于输出在模板中定义的名称为`“c”`的变量。当`$c`本身是一个`Struct`对象时，可访问其字段：`{{$c.Com}}`

在模板中定义变量：变量名称用字母和数字组成，并带上`“$”`前缀，采用符号`“:=”`进行赋值。

比如：`{{$x := "OK"}}` 或 `{{$x := pipeline}}`

# 管道函数

用法1：
```
{{FuncName1}}
```
此标签将调用名称为`“FuncName1”`的模板函数（等同于执行`“FuncName1()”`，不传递任何参数）并输出其返回值。

用法2：
```
{{FuncName1 "参数值1" "参数值2"}}
```
此标签将调用`“FuncName1("参数值1", "参数值2")”`，并输出其返回值

用法3：
```
{{.c|FuncName1}}
```
此标签将调用名称为`“FuncName1”`的模板函数（等同于执行`“FuncName1(this.c)”`，将竖线`“|”`左边的`“.c”`变量值作为函数参数传送）并输出其返回值。

# 判断
golang的模板也支持`if`的条件判断，当前支持最简单的`bool`类型和字符串类型的判断。
```
{{if .condition}}
{{end}}
```
当`.condition`为`bool`类型的时候，则为`true`表示执行，当`.condition`为`string`类型的时候，则非空表示执行。

当然也支持`else` ， `else if`嵌套
```
{{if .condition1}}
{{else if .contition2}}
{{end}}
```
假设我们需要逻辑判断，比如与或、大小不等于等判断的时候，我们需要一些内置的模板函数来做这些工作，目前常用的一些内置模板函数有：

## not 非
```
{{if not .condition}}
{{end}}
```
## and 与
```
{{if and .condition1 .condition2}}
{{end}}
```
## or 或
```
{{if or .condition1 .condition2}}
{{end}}
```
## eq 等于
```
{{if eq .var1 .var2}}
{{end}}
```
## ne 不等于
```
{{if ne .var1 .var2}}
{{end}}
```
## lt 小于 (less than)
```
{{if lt .var1 .var2}}
{{end}}
```
## le 小于等于
```
{{if le .var1 .var2}}
{{end}}
```
## gt 大于
```
{{if gt .var1 .var2}}
{{end}}
```
## ge 大于等于
```
{{if ge .var1 .var2}}
{{end}}
```
# 遍历
golang的`template`支持`range`循环来遍历`map、slice`内的内容，语法为：
```
{{range $i, $v := .slice}}
{{end}}
```
在这个`range`循环内，我们可以通过`i、v`来访问遍历的值，还有一种遍历方式为：
```
{{range .slice}}
{{end}}
```
这种方式无法访问到`index`或者`key`的值，需要通过.来访问对应的`value`
```
{{range .slice}}
{{.field}}
{{end}}
```
当然这里使用了.来访问遍历的值，那么我们想要在其中访问外部的变量怎么办？(比如渲染模板传入的变量)，在这里，我们需要使用`$.`来访问外部的变量
```
{{range .slice}}
{{$.c}}
{{end}}
```
# 预定义的模板全局函数
## and
```
{{and x y}}
```
表示：如果x为真，返回y，否则返回x。

等同于Golang中的：`x && y`

## call
```
{{call .X.Y 1 2}}
```
表示：
```
dot.X.Y(1, 2)
```
call后面的第一个参数的结果必须是一个函数（即这是一个函数类型的值），其余参数作为该函数的参数。

该函数必须返回一个或两个结果值，其中第二个结果值是error类型。

如果传递的参数与函数定义的不匹配或返回的error值不为nil，则停止执行。

## html

转义文本中的html标签，如将`“<”`转义为`“&lt;”`，`“>”`转义为`“&gt;”`等

## index
```
{{index x 1 2 3}}
```
返回`index`后面的第一个参数的某个索引对应的元素值，其余的参数为索引值

表示：`x[1][2][3]`

`x`必须是一个`map、slice`或数组

## js

返回用`JavaScript`的`escape`处理后的文本

## len

返回参数的长度值（`int`类型）

## not

返回单一参数的布尔否定值。

## or
```
{{or x y}}
```
表示：`if x then x else y`。等同于Golang中的：`x || y`

如果x为真返回x，否则返回y。

## print

`fmt.Sprint`的别名

## printf

`fmt.Sprintf`的别名

 
## println

`fmt.Sprintln`的别名

## urlquery

返回适合在`URL`查询中嵌入到形参中的文本转义值。（类似于`PHP`的`urlencode`）