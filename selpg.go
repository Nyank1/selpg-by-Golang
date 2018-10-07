package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	flag "github.com/spf13/pflag"
)

const INT_MAX = int(^uint(0) >> 1)

type selpg_args struct
{
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type bool
	print_dest string
}

var progname string

func usage() {
	fmt.Printf("Usage of %s:\n\n", progname)
	fmt.Printf("%s is a tool to select pages from what you want.\n\n", progname)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tselpg -sNumber -eNumber [options] [filename]\n\n")
	fmt.Printf("The arguments are:\n\n")
	fmt.Printf("\t-sNumber\tStart from Page <number>.\n")
	fmt.Printf("\t-eNumber\tEnd to Page <number>.\n")
	fmt.Printf("\t-lNumber\t[options]Specify the number of line per page. Default is 72.\n")
	fmt.Printf("\t-f\t\t[options]Specify that the pages are sperated by \\f.\n")
	fmt.Printf("\t-tDestination\tSend output to the printer of Destination.\n")
	fmt.Printf("\t[filename]\t[options]Read input from the file.\n\n")
	fmt.Printf("If no file specified, %s will read input from stdin. Control-D to end.\n\n", progname)
}

func main() {
	av := os.Args
	ac := len(av)
	progname = av[0]
	var sa selpg_args
	flag.IntVarP(&sa.start_page, "startpage", "s", -1, "Start page.")
	flag.IntVarP(&sa.end_page, "endpage", "e", -1, "End page.")
	flag.IntVarP(&sa.page_len, "length", "l", 72, "Line number per page.")
	flag.BoolVarP(&sa.page_type, "form", "f", false, "Determine form-feed-delimited")
	flag.StringVarP(&sa.print_dest, "destination", "d", "", "specify the printer")
	flag.Usage = usage
	flag.Parse()
	sa.start_page = -1
	sa.end_page = -1
	sa.in_filename = ""
	sa.page_len = 72
	sa.page_type = false
	sa.print_dest = ""
	process_args(ac, av, &sa)
	process_input(&sa)
}

func process_args(ac int, av []string, psa *selpg_args) {
	var s1 string
	var s2 string
	var argno int
	var i int
	if ac < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	s1 = av[1]
	if s1[0] != '-' || s1[1] != 's' {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		flag.Usage()
		os.Exit(2)
	}
	i = atoi(s1, 2)
	if i < 1 || i > INT_MAX - 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, s1[2:])
		flag.Usage()
		os.Exit(3)
	}
	psa.start_page = i
	s1 = av[2]
	if s1[0] != '-' || s1[1] != 'e' {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname);
		flag.Usage()
		os.Exit(4)
	}
	i = atoi(s1, 2)
	if i < 1 || i > INT_MAX - 1 || i < psa.start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, s1[2:]);
		flag.Usage()
		os.Exit(5)
	}
	psa.end_page = i
	argno = 3
	for argno <= ac - 1 && av[argno][0] == '-' {
		s1 = av[argno]
		switch s1[1] {
			case 'l':
				s2 = s1[2:]
				i = atoi(s2, 0)
				if i < 1 || i > INT_MAX - 1 {
					fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", progname, s2);
					flag.Usage()
					os.Exit(6)
				}
				psa.page_len = i
				argno++
			case 'f':
				if s1[0] != '-' || s1[1] != 'f' {
					fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname);
					flag.Usage()
					os.Exit(7)
				}
				psa.page_type = true
				argno++
			case 'd':
				s2 = s1[2:]
				if len(s2) < 1 {
					fmt.Fprintf(os.Stderr,"%s: -d option requires a printer destination\n", progname);
					flag.Usage()
					os.Exit(8)
				}
				psa.print_dest = s2
				argno++
			default:
				fmt.Fprintf(os.Stderr, "%s: unknown option %s\n", progname, s1);
				flag.Usage()
				os.Exit(9)
		}
	}
	if argno <= ac - 1 {
		psa.in_filename = av[argno]
		if _, err := os.Stat(psa.in_filename); err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", progname, psa.in_filename);
			os.Exit(10)
		}
		if file, err := os.Open(psa.in_filename); err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" exists but cannot be read\n", progname, psa.in_filename);
			os.Exit(11)
		} else {
			file.Close()
		}
	}
	if psa.start_page <= 0 || psa.end_page <= 0 || psa.end_page < psa.start_page || psa.page_len < 1 {
		fmt.Fprintf(os.Stderr, "%s: serious error\n", progname)
		os.Exit(12)
	}
}

func process_input(sa *selpg_args) {
	var stdin io.WriteCloser
	var s1 string
	var s2 string
	var cmd *exec.Cmd
	var err error
	line_ptr := 0
	page_ptr := 1
	if sa.print_dest == "" {
		stdin = nil
	} else {
		cmd = exec.Command("cat", "-n")
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
	}
	if flag.NArg() > 0 {
		if sa.in_filename != flag.Arg(0) {
			fmt.Fprintf(os.Stderr, "%s: serious error\n", progname)
			os.Exit(13)
		}
		output, err := os.Open(sa.in_filename)
		reader := bufio.NewReader(output)
		if sa.page_type == true {
			for i := 1; i <= sa.end_page; i++ {
				s1, err = reader.ReadString('\f')
				if err != io.EOF && err != nil {
					fmt.Println(err)
					os.Exit(14)
				}
				if err == io.EOF {
					break
				}
				if i >= sa.start_page {
					print(sa, s1, stdin)
				}
			}
		} else {
			for i := 0; true; i++ {
				line, _, err := reader.ReadLine()
				if err != io.EOF && err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err == io.EOF {
					break
				}
				if j := i/sa.page_len + 1; j >= sa.start_page && j <= sa.end_page {
					print(sa, string(line), stdin)
				} else if j >= sa.end_page {
					break
				}
			}
		}
	} else {
		s1 = ""
		s2 = ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s1 += scanner.Text()
			s1 += string("\n")
		}
		for i := 0; i < len(s1); i++ {
			if s1[i] != '\n' {
				s2 += string(s1[i])
			} else {
				if page_ptr >= sa.start_page && page_ptr <= sa.end_page {
					print(sa, s2, stdin)
				} else if page_ptr > sa.end_page {
					break
				}
				s2 = ""
				line_ptr++
				if line_ptr == sa.page_len {
					page_ptr++
					line_ptr = 0
				}
			}
		}
	}
	if sa.print_dest != "" {
		stdin.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func print(sa *selpg_args, line string, stdin io.WriteCloser) {
	if sa.print_dest != "" {
		stdin.Write([]byte(line + "\n"))
	} else {
		fmt.Println(line)
	}
}

func atoi(s string, i int) int {
	result := 0
	for j := i; j < len(s); j++ {
		result *= 10
		result += int(s[j]) - '0'
	}
	return result
}