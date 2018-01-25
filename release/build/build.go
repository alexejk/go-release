package build

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alexejk/go-release-tools/config"
)

type Builder struct {
	workDir string
}

func NewBuilder(workDir string) *Builder {

	b := &Builder{
		workDir: workDir,
	}

	return b
}

func (b *Builder) Build() error {

	buildCmd := config.GetString("project.build.command")
	args := strings.Split(buildCmd, " ")

	program := args[0]
	var pArgs []string
	if len(args) > 1 {
		pArgs = args[1:]
	}

	execCmd := exec.Command(program, pArgs...)
	execCmd.Dir = b.workDir
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = os.Environ()

	// Forward SIGINT, SIGTERM, SIGKILL to the child command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func() {
		sig := <-sigChan
		if execCmd.Process != nil {
			execCmd.Process.Signal(sig) //nolint: errcheck
		}
	}()

	return execCmd.Run()
}
