package client

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func (cli *RerokuCli) getMethod(name string) (reflect.Method, bool) {
	methodName := "Cmd" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	return reflect.TypeOf(cli).MethodByName(methodName)
}

func ParseCommands(args ...string) error {
	cli := NewRerokuCli(os.Stdin, os.Stdout, os.Stderr)
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
	return nil
}

func NewRerokuCli(in io.ReadCloser, out, err io.Writer) *RerokuCli {
	return &RerokuCli{
		in:  in,
		out: out,
		err: err,
	}
}

type RerokuCli struct {
	in  io.ReadCloser
	out io.Writer
	err io.Writer
}
