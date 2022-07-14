package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

// Config represents the JSON configuration structure.
type Config struct {
	Distro  string `json:"distro"`
	SSHPath string `json:"ssh_path"`
	WSLPath string `json:"wsl_path"`
}

// DefaultConfig is the default value for BypaSSH config.
var DefaultConfig = Config{
	Distro:  "Ubuntu",
	SSHPath: "/usr/bin/ssh",
	WSLPath: "C:\\Windows\\system32\\wsl.exe",
}

// parseConfig returns parsed config from file and any error encountered during
// parse.
func parseConfig() (conf Config, err error) {
	conf = DefaultConfig

	dir, err := os.Executable()
	if err != nil {
		return
	}
	dir = filepath.Dir(dir)

	fpath := filepath.Join(dir, "bypassh.json")
	f, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer f.Close()

	var fconf Config
	if err = json.NewDecoder(f).Decode(&fconf); err != nil {
		return
	}

	if fconf.Distro != "" {
		conf.Distro = fconf.Distro
	}
	if fconf.SSHPath != "" {
		conf.Distro = fconf.SSHPath
	}
	if fconf.WSLPath != "" {
		conf.WSLPath = fconf.WSLPath
	}
	return
}

// replaceWindowsPaths replaces in with WSL-compatible path. If there's any
// changes between in and out, ok will return true.
func replaceWindowsPaths(in string) (out string, ok bool) {
	var builder strings.Builder

	var pos int
	for {
		rpos := strings.Index(in[pos:], `:\`)
		if rpos <= 0 {
			break
		}

		npos := pos + rpos
		builder.WriteString(in[pos : npos-1])
		builder.WriteString("/mnt/")
		builder.WriteString(strings.ToLower(in[npos-1 : npos]))
		builder.WriteString("/")
		pos = npos + 2 // Skip the :\ part
	}
	builder.WriteString(in[pos:])

	return builder.String(), pos != 0
}

// translatePaths translates any Windows-styled paths and return them as
// WSL-styled paths.
func translatePaths(input []string, distro string) (output []string) {
	for _, in := range input {
		var isPath bool
		if strings.Contains(in, `\\wsl$\`) {
			in = strings.ReplaceAll(in, `\\wsl$\`+distro, "")
			isPath = true
		}
		var ok bool
		if in, ok = replaceWindowsPaths(in); ok {
			isPath = true
		}
		if isPath {
			in = strings.ReplaceAll(in, `\`, "/")
		}

		output = append(output, in)
	}
	return output
}

const helpMessage = `Usage: %s [-options...] destination [command]
Refer to ` + "`man ssh`" + ` for more information about the ssh arguments.

For BypaSSH configuration, create "bypassh.json" file in the same location as
the binary with the following fields and its default value:

{
  // The target WSL2 distro
  "distro":   "Ubuntu",

  // Path to the SSH binary inside WSL
  "ssh_path": "/usr/bin/ssh",

  // Path to the WSL binary in Windows
  "wsl_path": "C:\\Windows\\system32\\wsl.exe"
}

For most cases, if you are using the default Ubuntu distro you can leave the
configuration as-is. Use ` + "`-P`" + ` to print the JSON config and parse any
error message to debug your configuration.
`

func main() {
	conf, err := parseConfig()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		// Fast path for OpenSSH binary detection by Visual Studio Code.
		//
		// Using naive `ssh -V` from WSL will resulting in slower start up time
		// causing VSCode to give up and fallback to the default SSH binary.
		case "-V":
			fmt.Printf("wsl-ssh-helper 1.0.0 WSL2 OpenSSH-compatible proxy binary\n")
			return
		case "-h", "--help":
			fmt.Printf(helpMessage, os.Args[0])
			return
		case "-P":
			fmt.Printf("config: ")
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			enc.Encode(conf)
			if err != nil {
				fmt.Printf("parse error: %s\n", err.Error())
			}
			return
		}
	}

	args := translatePaths(os.Args[1:], conf.Distro)
	execArgs := append([]string{"-d", conf.Distro, conf.SSHPath}, args...)

	cmd := exec.Command(conf.WSLPath, execArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("error exec start: %s\n", err.Error())
		os.Exit(1)
	}

	// Interrupt signal forwarder
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			<-sig
			interruptCmd(cmd)
		}
	}()

	err = cmd.Wait()
	if exErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exErr.ExitCode())
	} else if err != nil {
		fmt.Printf("error exec wait: %s\n", err.Error())
		os.Exit(1)
	}
}
