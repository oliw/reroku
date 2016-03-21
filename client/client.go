package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oliw/reroku/server"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strings"
)

var VERSION string

func (cli *RerokuCli) getMethod(name string) (reflect.Method, bool) {
	methodName := "Cmd" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	return reflect.TypeOf(cli).MethodByName(methodName)
}

func ParseCommands(args ...string) error {
	cli := NewRerokuCli(os.Stdin, os.Stdout, os.Stderr, "tcp", server.DEFAULTHTTPHOST+":"+server.DEFAULTHTTPPORT)
	if len(args) > 0 {
		method, exists := cli.getMethod(args[0])
		if !exists {
			fmt.Fprintf(cli.err, "Error: Command not found %s\n", args[0])
		} else {
			ret := method.Func.CallSlice([]reflect.Value{reflect.ValueOf(cli), reflect.ValueOf(args[1:])})[0].Interface()
			if ret == nil {
				return nil
			}
			return ret.(error)
		}
	}
	return cli.CmdHelp(args...)
}

func (cli *RerokuCli) Subcmd(name, description string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(cli.err, "Usage: docker %s \n%s\n", name, description)
	}
	return flags
}

func (cli *RerokuCli) CmdHelp(args ...string) error {
	fmt.Fprintf(cli.err, "Usage: reroku [OPTIONS] COMMAND [arg...]\n")
	return nil
}

func (cli *RerokuCli) CmdVersion(args ...string) error {
	cmd := cli.Subcmd("version", "Show reroku version information")
	if err := cmd.Parse(args); err != nil {
		return nil
	}
	if cmd.NArg() > 0 {
		cmd.Usage()
		return nil
	}
	body, _, err := cli.call("GET", "/version")
	if err != nil {
		return err
	}
	var out server.APIVersion
	err = json.Unmarshal(body, &out)
	if err != nil {
		return err
	}
	fmt.Fprintf(cli.out, "Client Version: %s\n", VERSION)
	fmt.Fprintf(cli.out, "Server Version: %s\n", out.Version)
	return nil
}

func (cli *RerokuCli) call(method, path string) ([]byte, int, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, -1, err
	}
	req.Header.Set("User-Agent", "Docker-Client/"+VERSION)
	dial, err := net.Dial(cli.proto, cli.addr)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			fmt.Fprintf(cli.err, "Can't connect to Reroku daemon. Is it running?\n")
		}
		return nil, -1, err
	}
	clientconn := httputil.NewClientConn(dial, nil)
	resp, err := clientconn.Do(req)
	defer clientconn.Close()
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if len(body) == 0 {
			return nil, resp.StatusCode, fmt.Errorf("Error: %s", http.StatusText(resp.StatusCode))
		}
		return nil, resp.StatusCode, fmt.Errorf("Error: %s", body)
	}
	return body, resp.StatusCode, nil
}

func NewRerokuCli(in io.ReadCloser, out, err io.Writer, proto, addr string) *RerokuCli {
	return &RerokuCli{
		in:    in,
		out:   out,
		err:   err,
		proto: proto,
		addr:  addr,
	}
}

type RerokuCli struct {
	in    io.ReadCloser
	out   io.Writer
	err   io.Writer
	proto string
	addr  string
}
