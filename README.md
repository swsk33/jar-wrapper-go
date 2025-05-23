# jar包装器-Go

## 1，介绍
一个简单的命令行程序，包含一个Go语言代码模板，通过Go的`embed`标准库将`jar`文件嵌入至`exe`中，实现将`jar`打包为`exe`。

该包装器可以将`jar`文件打包成一个单独的`exe`文件，最终生成的`exe`文件不依赖于原始的`jar`文件，双击即可运行，目前支持生成`i386`（即`32`位）和`amd64`（即`64`位）的`exe`可执行文件。

## 2，安装和配置

### (1) 环境配置

该包装器使用Go语言实现，因此构建`exe`时需要借助Go语言的编译器，要求`1.23.0`及其以上版本。

首先安装Go语言开发工具，在[官网下载](https://go.dev/dl/)最新版即可，注意根据你的操作系统架构选择正确的安装包下载并安装，安装完成即可，安装后正常情况下`go`命令会被自动地配置到你的环境变量`Path`中去，可以打开命令行执行`go version`进行确认。

除此之外，**如果你想自定义生成的`exe`文件图标**，还需要安装`go-winres`命令，配置完成Go开发环境之后执行下列命令安装`go-winres`：

```bash
go install github.com/tc-hib/go-winres@latest
```

> 若安装速度很慢，可以先尝试配置[七牛云镜像](https://goproxy.cn/)后再试。

### (2) 下载该工具并配置环境变量

在仓库页面右侧**发行版/Releases**处即可下载该程序的压缩包，根据你自己的系统架构选择下载`i386`还是`amd64`的，下载后将其中的内容解压至一个目录中，并将该目录添加至`Path`环境变量。

然后打开终端，测试命令`jar2exe-go`是否可用：

```bash
jar2exe-go version
```

如果能输出版本信息说明配置成功。

### (3) 命令自动补全脚本

下载解压后，除了命令程序本身`exe`文件、`wrapper`包装器代码模板之外，还有下列脚本文件，是用于命令自动补全的，将该脚本配置到对应的终端后，使用`jar2exe-go`命令时即可按下Tab键自动补全命令：

- `jar2exe-go.fish` 用于Fish Shell的自动补全脚本，在Windows中通常是在Msys2环境下运行Fish Shell，将该文件放到`你的Msys2安装目录\etc\fish\completions`目录下即可，请勿修改文件名
- `jar2exe-go-completion.bash` 用于Bash Shell的自动补全脚本，在Windows中可以在Git Bash或者Msys2环境中运行Bash Shell，这里分别说明：
	- 使用Git Bash时，将`jar2exe-go-completion.bash`文件的扩展名改成`sh`，然后放到`你的Git安装目录\etc\profile.d`目录下
	- 使用Msys2时，直接把`jar2exe-go-completion.bash`文件放在`你的Msys2安装目录\etc\bash_completion.d`目录下

将文件放在对应位置后，配置就完成了，重启终端，后续在使用`jar2exe-go`命令时，即可使用Tab键补全命令。

## 3，打包教程

使用`jar2exe-go`工具将`jar`打包成`exe`是非常简单的，执行下列命令即可查看帮助：

```bash
jar2exe-go -h
```

下面就来详细讲解一下将一个`jar`文件打包成`exe`的详细步骤。

### (1) 准备`jar`文件并生成配置

准备好你要打包的`jar`文件，例如我这里有一个`demo.jar`位于当前目录下：

![image-20231013194113572](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013194113572.png)

这时首先需要生成一个配置文件，该配置文件中存放的是一些和`jar`运行相关的配置，执行下列命令：

```bash
jar2exe-go init-config
```

此时就会在**当前路径**下生成一个名为`config.yaml`的配置文件：

![image-20231013194254181](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013194254181.png)

### (2) 按需修改配置文件

生成的`config.yaml`配置文件为`YAML`格式，若对这个语法不太熟悉可以参考菜鸟教程：[传送门](https://www.runoob.com/w3cnote/yaml-intro.html)

初始时配置文件内容如下：

```yaml
# 配置文件

# 运行配置
run:
  # Java的运行路径，默认安装了Java运行环境的电脑直接填"java"即可，便携式jre需要在此指定，若指定为便携jre，则必须是相对路径，相对于生成的exe可执行文件的路径
  # 例如指定为"java"时，那么程序会使用该命令执行jar文件：java -jar path/to/xxx.jar
  # 如果指定为便携的jre，例如"jre/bin"，那么程序会使用该命令执行jar文件：jre/bin/java -jar path/to/xxx.jar
  java-path: "java"
  # 没有检测到Java运行环境时的提示内容
  error-message: "未能找到可用的Java运行环境（jre）！请先安装Java 8运行环境！"
  # 前置运行参数，即双击exe后自动加上的命令行参数，这个参数会先于命令行运行exe时加上的参数，该配置项为字符串数组类型
  pre-parameters:
# 日志相关
log:
  # 是否把程序的标准输出和标准错误重定向到本地文件，建议控制台应用程序不要开启此项
  write-to-file: false
  # 日志文件输出路径，为相对路径，为相对于可执行文件的位置，若上面变量write-to-file为false，则此配置项无效
  path: "error.log"
# 构建器自动生成和修改的选项，请勿手动修改！
build:
  # 是否为窗体应用程序
  win-app: false
  # 是否使用内嵌的JRE
  # 若开启此项，run.java-path将无效
  use-embed-jre: false
```

在配置文件中有着比较详细的注释，大家可以根据注释进行修改。

可能大多数情况下不需要修改配置文件，保持默认值即可。

### (3) 执行构建命令

现在通过下列命令将`demo.jar`打包为`demo.exe`并保存在当前目录下：

```bash
jar2exe-go build demo.jar demo.exe
```

> 可以指定相对路径或者绝对路径，若路径中带有空格，则需要用英文双引号`"`包围路径。

显示构建成功即可：

![image-20241005234121458](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20241005234121458.png)

这时，我们就成功地生成了`exe`文件了！

![image-20231013194926147](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013194926147.png)

双击`exe`即可运行，或者通过`cmd`或者其它命令行终端调用也可以，并且可以正常地接受命令行参数。

该`exe`文件可以单独地存在，不依赖于原有的`jar`文件，但是这并不意味着它脱离了Java运行环境后还能够运行，当然如果你希望`exe`文件能够脱离Java运行环境也可以正常运行的话，也是可以做到的，参考下面内嵌JRE部分。

## 4，常见使用示例

### (1) 构建时指定配置文件位置

上述教程中，构建`exe`时并没有指定配置文件的位置，这是因为**默认情况下构建程序会使用当前路径下的`config.yaml`作为配置文件**。

如果说你的配置文件位于其它位置，或者文件名和默认情况不符，则可以加上`-c`参数手动指定配置文件：

```bash
# 配置文件位于当前目录下，名为cfg.yaml
jar2exe-go build demo.jar demo.exe -c cfg.yaml

# 配置文件位于当前目录下的cfg文件夹中，名为config.yaml
jar2exe-go build demo.jar demo.exe -c ./cfg/config.yaml

# 配置文件位于上一级目录下，名为config.yaml
jar2exe-go build demo.jar demo.exe -c ../config.yaml
```

同样地，如果说指定的配置文件路径包含空格，则需要使用英文双引号进行包围。

### (2) 自定义输出`exe`的图标

如果说需要自定义输出的`exe`文件的图标，需要确保你的电脑上已经安装了`go-winres`命令，这方面在上面安装和配置部分已经讲解。

我们需要准备一个图片文件作为生成的`exe`的图标，该图片要满足下列要求：

- `png`格式
- 分辨率不大于`256 x 256`

例如我这里准备了一个`gopher.png`作为图标文件：

![image-20231013200601997](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013200601997.png)

那么构建时，使用`-i`参数指定图片文件路径即可：

```bash
# gopher.png位于当前路径下
jar2exe-go build demo.jar demo.exe -i gopher.png
```

这样，我们就能够得到一个带有图标的`exe`文件：

![image-20231013200800250](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013200800250.png)

### (3) 前置运行参数

假设我们的`jar`文件可以传递命令行参数，并且我们想每次运行这个程序的时侯就自动地添加一些命令行参数，而不是手动地使用命令行传参，那么我们可以通过修改**配置文件**中的`run.pre-parameters`配置项实现。

添加了前置的运行参数后，**每次通过`exe`运行程序时，会先自动地传递我们指定的前置运行参数给我们的Java程序**。

例如修改配置文件如下：

```yaml
# 运行配置
run:
  # 省略其它部分...
  pre-parameters:
  - "a"
  - "b"
  - "c"
# 省略其它部分...
```

可见这个配置项是**字符串数组**类型的，我们添加了`3`个前置运行参数，这样运行`exe`时，这三个参数就会按照顺序自动地传递给Java程序。

在此构建程序，并直接点击`exe`运行，可见能够成功地传递这些前置参数：

![image-20231013201353003](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013201353003.png)

当然，如果我们通过命令行调用程序并传参，程序仍然可以正常接收，只不过命令行传入的参数位于前置参数之后：

```bash
./demo.exe 1 2 3
```

结果：

![image-20231013201458211](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013201458211.png)

### (4) 日志重定向

默认情况下，`jar`程序向终端输出的消息会打印到终端上，但是有些时候我们需要日志重定向至文件而非终端，修改**配置文件**中的`log.write-to-file`为`true`即可开启日志重定向功能。

还可以定义`log.path`配置项表示配置文件保存位置，可以使用相对路径或者绝对路径，相对路径是相对于生成的`exe`的文件路径。

这个功能在GUI窗体应用程序中非常实用，因为窗体应用程序通常不会显示命令行窗口，这样将输出重定向到日志就可以更好地进行问题排查。

需要注意的是，如果开启了日志重定向功能，那么`jar`向控制台输出的内容则不会再显示在终端里，因此不推荐命令行程序开启该功能。

### (5) 打包GUI窗体应用程序

假设你的`jar`是GUI窗体应用程序，那么按照上述方式进行打包，启动窗体应用程序的时候还会出现一个命令行窗口，这样非常的不美观。

因此你可以在构建时加上`-g`参数，这样启动程序时就不会显示命令行窗口了！

```bash
jar2exe-go build demo-gui.jar demo-gui.exe -g
```

除了窗口应用程序，一些后台应用程序也可以使用`-g`隐藏控制台窗口。

### (6) 使得`exe`在没有安装Java运行环境的电脑上也能运行

虽然我们可以把`jar`打包成`exe`文件，但是**这并不代表着缺失Java运行环境时我们的程序仍然可以正常运行**。

如果说你想要发布你的应用程序，但是又不希望用户去安装Java运行环境，也是有很多办法实现的。

#### ① 便携JRE

如果希望你的程序拿来就可以运行，通常的做法就是把Java运行环境目录复制过来和你的程序放在一起，然后打包给用户。

这样，我们需要修改配置文件中的`run.java-path`配置指定你的JRE路径，也就是`java`命令所在的路径。

假设现在我把Java运行环境拷贝了过来，并且命名为`jre`放在当前目录下：

![image-20231013203145657](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231013203145657.png)

并且`java.exe`位于`./jre/bin`里面。

这时修改配置文件：

```yaml
run:
  # 省略其它...
  java-path: "jre/bin"
# 省略其它...
```

可见这里也是指定相对路径，相对于你的`exe`文件的路径。

构建后，保证`exe`文件和这个`jre`文件夹位于同级目录，`exe`即可调用这个`jre`中的`java`来运行程序，而不再依赖本机的Java运行环境。

如果你默认用户电脑已经安装了Java运行环境，则保持`run.java-path`为默认值`java`即可。

> JDK 9及其以上版本支持根据`jar`程序生成一个精简版的JRE，具体操作方式可以查看：[传送门](https://juejin.cn/post/6989214620279898126)

除此之外，在程序运行但是找不到`java`运行环境的时候，会弹出一个错误提示框来告知用户找不到JRE，这个错误提示消息可以在配置文件中`run.error-message`配置项进行自定义。

#### ② 内嵌JRE

虽然便携JRE可以保证即使用户没有安装Java运行环境，也能够使得`exe`正常运行，但是这也有一个缺点：依赖一个外部的JRE文件夹，如果用户移动了`exe`位置或者是JRE文件夹，那么也会导致`exe`找不到Java运行环境而运行失败。

因此我们还可以采用另一种方式：**将整个JRE也内嵌到`exe`中去**，这样这个`exe`文件可以完全地单独存在，并且本地未安装Java运行环境时也能够正常地运行。

在进行构建时，传入`--embed-jre`参数就可以开启内嵌JRE功能，然后使用`--embed-jre-path`参数指定你要内嵌的JRE文件夹。

例如我这里已经把JRE文件夹复制到当前目录下了：

![image-20231014222754866](https://swsk33-note.oss-cn-shanghai.aliyuncs.com/image-20231014222754866.png)

那么`java`命令位于：`./jre17/bin/java.exe`

那么执行下列构建命令：

```bash
jar2exe-go build demo.jar demo.exe --embed-jre --embed-jre-path ./jre17
```

这样，生成的`exe`中就是包含了你指定的JRE的，这个`exe`可以完全地单独存在，不依赖于任何外部的JRE就可以正常运行。

需要注意的是，内嵌JRE功能虽然很好地提升了生成的`exe`文件的可移植性，但是这样也会大幅增加`exe`文件的大小，并且可能会导致`exe`文件启动速度变慢，需要根据实际情况选择是否使用内嵌JRE功能。

除此之外，内嵌的JRE架构（`32`位还是`64`位的JRE）最好是保证和`exe`架构一致，否则可能会在运行时出现意想不到的错误。

#### ③ 自动内嵌JRE

在上述使用内嵌JRE功能时，我们仍然需要手动复制一个JRE文件夹过来，或者是使用`jlink`命令生成JRE文件夹，这样也有点麻烦。

因此，构建器还提供了自动内嵌JRE功能，只需要加上一个参数，即可在构建`exe`时自动地根据`jar`文件依赖情况，生成一个精简版JRE，并且自动地内嵌至`exe`文件中去，无需像上面一样手动指定JRE路径。

在构建时加上`--auto-embed-jre`参数即可一键生成并嵌入JRE到`exe`中去：

```bash
jar2exe-go build demo.jar demo.exe --auto-embed-jre
```

这样，得到的`exe`也是包含了JRE的，可以完全单独运行。

使用自动内嵌JRE功能有以下几点需要注意：

- 要求你的电脑已经安装并正确配置了JDK 9及其以上版本的JDK（`jdeps`和`jlink`命令必须可用）
- 内嵌的JRE版本和架构（`32`位还是`64`位）也取决于你本地安装的JDK版本和架构

因此如果你使用的是JDK 8，或者你需要在`64`位电脑上生成一个完全的`32`位程序并内嵌`32`位JRE，那么自动内嵌JRE功能是不适合的。

> 请注意，对于使用Spring Boot框架开发的Java程序`jar`文件，使用自动内嵌JRE功能可能无法打包，因为`deps`命令本身无法分析Spring Boot打包的`fatJar`文件，请手动下载JRE运行环境并手动内嵌。

### (7) 生成指定架构的可执行文件

默认情况下，生成的`exe`文件是`64`位的，这代表着在`32`位电脑上无法运行。

不过我们可以在构建时通过`-a`参数指定生成的`exe`架构，可以指定为如下值：

- `i386` 生成`32`位`exe`文件
- `amd64` 生成`64`位`exe`文件

```bash
# 生成32位exe程序
jar2exe-go build demo-gui.jar demo-gui.exe -a i386
```

一般来说，`32`位`exe`可以在`32`位和`64`位操作系统上运行，而`64`位`exe`只能够在`64`位操作系统上运行。

### (8) 减小构建文件体积

在进行构建时，构建器默认去除了`exe`文件的调试信息，已经减小了体积，不过如果想进一步地减小`exe`体积，可以借助UPX工具实现，参考这篇文章的第`3`部分：[传送门](https://juejin.cn/post/6989149609478553630#heading-2)

## 5，源代码目录说明

该目录下有下列Go模块的源代码目录：

- `wrapper` 打包器代码，这个代码实现了将`jar`打包为`exe`，以及运行`jar`的逻辑
- `builder` 构建器代码，能够读取命令行参数，然后调用`go build`命令构建`wrapper`目录的源代码完成`jar`到`exe`的构建

## 6，大致原理

该程序的原理其实很简单，主要是利用了Go的`embed`标准库，实现将`jar`文件以及配置文件作为资源文件，内嵌至最终的可执行文件中去。

在用户运行`exe`时，又会将`jar`和配置文件释放至临时目录中，并且读取配置文件并调用`java`命令运行`jar`文件。

具体实现可以参考：

- Go语言内嵌资源文件：[传送门](https://juejin.cn/post/7288963080855568421)
- Go语言自定义`exe`图标：[传送门](https://juejin.cn/post/7281118263966351372)
