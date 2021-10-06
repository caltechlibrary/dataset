/**
 * (depreciated code) - this code will remain for a version or two
 * as Caltech Library projects are updated.
 */
package cli

import (
	"fmt"
	"io"
)

// (depreciated in favor of Verb) Action describes an "action" that a cli might take. Actions aren't prefixed with a "-".
type Action struct {
	// Name is usually a verb like list, test, build as needed by the cli
	Name string
	// Fn is action that will be run by Cli.Run() if Name is the first non-option arg on the command line
	//NOTE: currently the signature is io based but may be changed to *os.File in
	// the future
	Fn func(io.Reader, io.Writer, io.Writer, []string) int
	// Usage is a short description of what the action does and description of any expected additoinal parameters
	Usage string
}

// (depreciated) String prints an actions' name and description
func (a *Action) String() string {
	return fmt.Sprintf("%s - %s", a.Name, a.Usage)
}

// (depreciated) AddVerb associates a verb and synopsis without assigning a function
// (e.g. if you aren't going to use cli.Run()
func (c *Cli) AddVerb(verb string, usage string) error {
	c.actions[verb] = &Action{
		Name:  verb,
		Usage: usage,
	}
	_, ok := c.actions[verb]
	if ok == false {
		return fmt.Errorf("Failed to add verb docs for %q", verb)
	}
	return nil
}

// (depreciated) AddAction associates a wrapping function with a action name, the wrapping function
// has 4 parameters in io.Reader, out io.Writer, err io.Writer, args []string. It should return
// an integer reflecting an exit code like you'd pass to os.Exit().
func (c *Cli) AddAction(verb string, fn func(io.Reader, io.Writer, io.Writer, []string) int, usage string) error {
	c.actions[verb] = &Action{
		Name:  verb,
		Fn:    fn,
		Usage: usage,
	}
	_, ok := c.actions[verb]
	if ok == false {
		return fmt.Errorf("Failed to add action %q", verb)
	}
	return nil
}

// (depreciated) Action returns a doc string for a given verb
func (c *Cli) Action(verb string) string {
	action, ok := c.actions[verb]
	if ok == false {
		return fmt.Sprintf("%q is not a defined action", verb)
	}
	return action.Usage
}

// (depreciated) Actions returns a map of actions and their doc strings
func (c *Cli) Actions() map[string]string {
	actions := map[string]string{}
	for k, action := range c.actions {
		actions[k] = action.Usage
	}
	return actions
}
