# CLI 命令行实用程序开发基础 实验报告
## 16340299 赵博然
# 设计说明
参考以下C代码, 和CSDN博客. 原谅我英语水平太差, `usage()`函数是照着博客写的.

https://www.ibm.com/developerworks/cn/linux/shell/clutil/selpg.c

https://blog.csdn.net/qq_33454112/article/details/78283471

参见`selpg.go`.

引用如下包.
```
import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    flag "github.com/spf13/pflag"
)
```
常量`INT_MAX`, 用来检查数字的合法性.
```
const INT_MAX = int(^uint(0) >> 1)
```
结构体`selpg_args`. 与C代码是一样的.
```
type selpg_args struct
{
    start_page int
    end_page int
    in_filename string
    page_len int
    page_type bool
    print_dest string
}
```
函数.
```
func usage() //命令的说明. 当命令格式错误时输出.
func main() //程序入口. 在初始化selpg_args, pflag后执行下两个函数.
func process_args(ac int, av []string, psa *selpg_args) //解析命令. 翻译自C代码.
func process_input(sa *selpg_args) //计算输出. 包括是否使用-dDestination, 以及使用-f或-lnumber等不同方式和格式的输出. 在每计算完一行后使用下个函数输出.
func print(sa *selpg_args, line string, stdin io.WriteCloser) //输出. 包括是否使用-dDestination的两种方式的输出.
func atoi(s string, i int) int //相当于C语言的atoi()函数, 将s字符串i位置后的子串解析为整数.
```
# 测试结果
测试内容见每次测试的终端命令.

---
## 测试1
新建测试文件`test`.
```
good morning
zaoshanghao
ohayougozaimasu
hello
nihao
konnichiwa
i am hungry
woele
onakagasukimashita
i am full
wochibaole
itadakimashita
```
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 test 
hello
nihao
konnichiwa
i am hungry
woele
onakagasukimashita
```
---
## 测试2
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 <test
hello
nihao
konnichiwa
i am hungry
woele
onakagasukimashita
```
---
## 测试3
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 test >outputtest
[yaroglek@centos-new 文档]$ cat outputtest 
hello
nihao
konnichiwa
i am hungry
woele
onakagasukimashita
```
---
## 测试4
```
[yaroglek@centos-new 文档]$ selpg -s0 -e3 -l3 test
selpg: invalid start page 0
Usage of selpg:

selpg is a tool to select pages from what you want.

Usage:

	selpg -s=Number -e=Number [options] [filename]

The arguments are:

	-s=Number	Start from Page <number>.
	-e=Number	End to Page <number>.
	-l=Number	[options]Specify the number of line per page. Default is 72.
	-f		[options]Specify that the pages are sperated by \f.
	[filename]	[options]Read input from the file.

If no file specified, selpg will read input from stdin. Control-D to end.


```
---
## 测试5
```
[yaroglek@centos-new 文档]$ selpg -s0 -e3 -l3 test 2>errortest
Usage of selpg:

selpg is a tool to select pages from what you want.

Usage:

	selpg -s=Number -e=Number [options] [filename]

The arguments are:

	-s=Number	Start from Page <number>.
	-e=Number	End to Page <number>.
	-l=Number	[options]Specify the number of line per page. Default is 72.
	-f		[options]Specify that the pages are sperated by \f.
	[filename]	[options]Read input from the file.

If no file specified, selpg will read input from stdin. Control-D to end.

[yaroglek@centos-new 文档]$ cat errortest 
selpg: invalid start page 0
```
---
## 测试6
```
[yaroglek@centos-new 文档]$ selpg -s0 -e3 -l3 test >outputtest2 2>errortest2
[yaroglek@centos-new 文档]$ cat outputtest2
Usage of selpg:

selpg is a tool to select pages from what you want.

Usage:

	selpg -s=Number -e=Number [options] [filename]

The arguments are:

	-s=Number	Start from Page <number>.
	-e=Number	End to Page <number>.
	-l=Number	[options]Specify the number of line per page. Default is 72.
	-f		[options]Specify that the pages are sperated by \f.
	[filename]	[options]Read input from the file.

If no file specified, selpg will read input from stdin. Control-D to end.

[yaroglek@centos-new 文档]$ cat errortest2
selpg: invalid start page 0
```
---
## 测试7
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 test >/dev/null
[yaroglek@centos-new 文档]$ cat /dev/null 

```
---
## 测试8
导入hello命令. hello.go的代码如下.
```
package main

import (
    "fmt"
    "os"
)

func main() {
    for i, a := range os.Args[1:] {
        fmt.Printf("Argument %d is %s\n", i+1, a)
    }

}
```
```
[yaroglek@centos-new 文档]$ hello mynameissuzuki woshilingmu suzukitoiimasu whoareyou nishishui donatadesuka iammiura wojiaosanpu miuradesu | selpg -s1 -e2 -l3
Argument 1 is mynameissuzuki
Argument 2 is woshilingmu
Argument 3 is suzukitoiimasu
Argument 4 is whoareyou
Argument 5 is nishishui
Argument 6 is donatadesuka
```
---
## 测试9
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 test | wc
      6       8      60
```
---
## 测试10
```
[yaroglek@centos-new 文档]$ selpg -s0 -e3 -l3 test 2>errortest3 | wc
     18      72     469
[yaroglek@centos-new 文档]$ cat errortest3
selpg: invalid start page 0
```
---
## 测试11
新建测试文件`ftest`.
```
[yaroglek@centos-new 文档]$ echo -e good morning'\n'morning'\f'good afternoon'\n'afternoon'\f'good evening'\n'evening'\f'good night'\n'night >ftest
[yaroglek@centos-new 文档]$ cat ftest
good morning
morning
       good afternoon
afternoon
         good evening
evening
       good night
night
```
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -f ftest 
good afternoon
afternoon

good evening
evening

```
---
## 测试12
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 -dlp1 test 
     1	hello
     2	nihao
     3	konnichiwa
     4	i am hungry
     5	woele
     6	onakagasukimashita
```
---
## 测试13
在输入`ten`换行后键入`Control-D`.
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 -dlp1
one
two
three
four 
five
six
seven
eight
nine
ten
     1	four
     2	five
     3	six
     4	seven
     5	eight
     6	nine
```
## 测试14
```
[yaroglek@centos-new 文档]$ selpg -s2 -e3 -l3 test >outputtest3 2>errortest4 &
[1] 25657
[1]+  完成                  selpg -s2 -e3 -l3 test > outputtest3 2> errortest4
[yaroglek@centos-new 文档]$ cat outputtest3
hello
nihao
konnichiwa
i am hungry
woele
onakagasukimashita
[yaroglek@centos-new 文档]$ cat errortest4
```
---
## 测试15
在基本确认`selpg`命令无明显bug后, 将`selpg`进行应用. `hamlet.txt`是一部小说, 篇幅较长, 因此使用一般的阅读器如`Vim`, `记事本`等并不容易读取后面内容. 使用`selpg`命令读取`hamlet.txt`的第15页, 规定每页72行.
```
[yaroglek@centos-new 文档]$ selpg -s15 -e15 -dlp1 hamlet.txt 
     1	  Ham. Hic et ubique? Then we'll shift our ground.
     2	    Come hither, gentlemen,
     3	    And lay your hands again upon my sword.
     4	    Never to speak of this that you have heard:
     5	    Swear by my sword.
     6	  Ghost. [beneath] Swear by his sword.
     7	  Ham. Well said, old mole! Canst work i' th' earth so fast?
     8	    A worthy pioner! Once more remove, good friends."
     9	  Hor. O day and night, but this is wondrous strange!
    10	  Ham. And therefore as a stranger give it welcome.
    11	    There are more things in heaven and earth, Horatio,
    12	    Than are dreamt of in your philosophy.
    13	    But come!
    14	    Here, as before, never, so help you mercy,
    15	    How strange or odd soe'er I bear myself
    16	    (As I perchance hereafter shall think meet
    17	    To put an antic disposition on),  
    18	    That you, at such times seeing me, never shall,
    19	    With arms encumb'red thus, or this head-shake,
    20	    Or by pronouncing of some doubtful phrase,
    21	    As 'Well, well, we know,' or 'We could, an if we would,'
    22	    Or 'If we list to speak,' or 'There be, an if they might,'
    23	    Or such ambiguous giving out, to note
    24	    That you know aught of me- this is not to do,
    25	    So grace and mercy at your most need help you,
    26	    Swear.
    27	  Ghost. [beneath] Swear.
    28	                                                   [They swear.]
    29	  Ham. Rest, rest, perturbed spirit! So, gentlemen,
    30	    With all my love I do commend me to you;
    31	    And what so poor a man as Hamlet is
    32	    May do t' express his love and friending to you,
    33	    God willing, shall not lack. Let us go in together;
    34	    And still your fingers on your lips, I pray.
    35	    The time is out of joint. O cursed spite
    36	    That ever I was born to set it right!
    37	    Nay, come, let's go together.  
    38	                                                         Exeunt.
    39	
    40	
    41	
    42	
    43	Act II. Scene I.
    44	Elsinore. A room in the house of Polonius.
    45	
    46	Enter Polonius and Reynaldo.
    47	
    48	  Pol. Give him this money and these notes, Reynaldo.
    49	  Rey. I will, my lord.
    50	  Pol. You shall do marvell's wisely, good Reynaldo,
    51	    Before You visit him, to make inquire
    52	    Of his behaviour.
    53	  Rey. My lord, I did intend it.
    54	  Pol. Marry, well said, very well said. Look you, sir,
    55	    Enquire me first what Danskers are in Paris;
    56	    And how, and who, what means, and where they keep,
    57	    What company, at what expense; and finding
    58	    By this encompassment and drift of question
    59	    That they do know my son, come you more nearer
    60	    Than your particular demands will touch it.
    61	    Take you, as 'twere, some distant knowledge of him;
    62	    As thus, 'I know his father and his friends,
    63	    And in part him.' Do you mark this, Reynaldo?
    64	  Rey. Ay, very well, my lord.  
    65	  Pol. 'And in part him, but,' you may say, 'not well.
    66	    But if't be he I mean, he's very wild
    67	    Addicted so and so'; and there put on him
    68	    What forgeries you please; marry, none so rank
    69	    As may dishonour him- take heed of that;
    70	    But, sir, such wanton, wild, and usual slips
    71	    As are companions noted and most known
    72	    To youth and liberty.
```
---
