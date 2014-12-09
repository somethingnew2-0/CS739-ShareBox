// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"replica"
	"strconv"
	"strings"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  void ping()")
	fmt.Fprintln(os.Stderr, "  void add(Replica r)")
	fmt.Fprintln(os.Stderr, "  void modify(Replica r)")
	fmt.Fprintln(os.Stderr, "  void remove(string shardId)")
	fmt.Fprintln(os.Stderr, "  Replica download(string shardId)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	var parsedUrl url.URL
	var trans thrift.TTransport
	_ = math.MinInt32 // will become unneeded eventually
	_ = strconv.Atoi
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Parse()

	if len(urlString) > 0 {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	client := replica.NewReplicatorClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "ping":
		if flag.NArg()-1 != 0 {
			fmt.Fprintln(os.Stderr, "Ping requires 0 args")
			flag.Usage()
		}
		fmt.Print(client.Ping())
		fmt.Print("\n")
		break
	case "add":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Add requires 1 args")
			flag.Usage()
		}
		arg22 := flag.Arg(1)
		mbTrans23 := thrift.NewTMemoryBufferLen(len(arg22))
		defer mbTrans23.Close()
		_, err24 := mbTrans23.WriteString(arg22)
		if err24 != nil {
			Usage()
			return
		}
		factory25 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt26 := factory25.GetProtocol(mbTrans23)
		argvalue0 := replica.NewReplica()
		err27 := argvalue0.Read(jsProt26)
		if err27 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Add(value0))
		fmt.Print("\n")
		break
	case "modify":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Modify requires 1 args")
			flag.Usage()
		}
		arg28 := flag.Arg(1)
		mbTrans29 := thrift.NewTMemoryBufferLen(len(arg28))
		defer mbTrans29.Close()
		_, err30 := mbTrans29.WriteString(arg28)
		if err30 != nil {
			Usage()
			return
		}
		factory31 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt32 := factory31.GetProtocol(mbTrans29)
		argvalue0 := replica.NewReplica()
		err33 := argvalue0.Read(jsProt32)
		if err33 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Modify(value0))
		fmt.Print("\n")
		break
	case "remove":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Remove requires 1 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		fmt.Print(client.Remove(value0))
		fmt.Print("\n")
		break
	case "download":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Download requires 1 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		fmt.Print(client.Download(value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
