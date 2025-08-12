package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/lorentzforces/selfman/internal/git"
	"github.com/lorentzforces/selfman/internal/ops"
	"github.com/lorentzforces/selfman/internal/run"
	"github.com/spf13/cobra"
)

func validatePrereqs() error {
	if !git.ExecExists() {
		return fmt.Errorf("Cannot find a \"git\" executable on PATH")
	}
	return nil
}

type SelfmanCommand struct {
	cobraCmd *cobra.Command
	runFunc func(*cobra.Command, []string) (*SelfmanResult, error)
}

type SelfmanResult struct {
	// Text output is always printed before any other messages
	textOutput fmt.Stringer
	// Any mutating operations to be executed as a result of running this command
	operations []ops.Operation
}

func (self *SelfmanCommand) RunSelfmanCommand(cmd *cobra.Command, args []string) error {
	cmdResult, err := self.runFunc(cmd, args)
	if err != nil { return err }

	dryRun, err := cmd.Flags().GetBool(globalOptionDryRun)
	run.AssertNoErr(err)
	isVerbose, err := cmd.Flags().GetBool(globalOptionVerbose)
	run.AssertNoErr(err)

	var verbosity VerbosityLevel = NotVerbose
	if isVerbose {
		verbosity = Verbose
	}

	if cmdResult.textOutput != nil {
		fmt.Println(cmdResult.textOutput)
	}

	if dryRun {
		dryRunOperations(cmdResult.operations, verbosity)
		return nil
	} else {
		return executeOperations(cmdResult.operations, verbosity)
	}
}

// Since the messages printed herein are progress updates, print to stderr
func executeOperations(actions []ops.Operation, verbosity VerbosityLevel) error {
	for _, action := range actions {
		fmt.Fprintln(os.Stderr, printOperation(action, verbosity))
		msg, err := action.Execute()
		if err != nil { return err }

		fmt.Fprintf(os.Stderr, "✓")
		if len(msg) > 0 {
			fmt.Fprintf(os.Stderr, " %s", msg)
		}
		fmt.Fprintln(os.Stderr)
	}

	return nil
}

type VerbosityLevel int
const (
	Verbose VerbosityLevel = iota
	NotVerbose
)

// Since this is asked for as the main output, print to stdout
func dryRunOperations(actions []ops.Operation, verbosity VerbosityLevel) {
	fmt.Printf("Would perform the following operations:\n\n")
	for _, action := range actions {
		fmt.Println(printOperation(action, verbosity))
	}
}

func printOperation(op ops.Operation, verbosity VerbosityLevel) string {
	var buf strings.Builder
	writeOutOpWithIndent(&buf, 0, op, verbosity)
	return buf.String()
}

// print only this operation's information, regardless of whether it's a meta operation with nested
// ops
func printFlatOperation(op ops.Operation, verbosity VerbosityLevel, totalIndent string) string {
	if verbosity == Verbose {
		return op.Describe().LongDisplayWithIndent(totalIndent)
	}
	return op.Describe().ShortDisplayWithIndent(totalIndent)
}

func writeOutOpWithIndent(
	buf *strings.Builder,
	indentLevel int,
	op ops.Operation,
	verbosity VerbosityLevel,
) {
	totalIndent := strings.Repeat(run.IndentChars, indentLevel)
	buf.WriteString(printFlatOperation(op, verbosity, totalIndent))

	metaOp, ok := op.(ops.MetaOperation)
	// if not a meta op, we've already printed this operation and we can bail
	if !ok { return }

	for _, nestedOp := range metaOp.InnerOps() {
		buf.WriteString("\n")
		writeOutOpWithIndent(
			buf,
			indentLevel + 1,
			nestedOp,
			verbosity,
		)
	}
}
