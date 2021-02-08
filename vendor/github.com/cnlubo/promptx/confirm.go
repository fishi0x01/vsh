package promptx

import (
	"github.com/cnlubo/promptx/utils"
	"github.com/mritd/readline"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"text/template"
)

const (
	yesNoPrompt                 = "y/N"
	defaultMessageTpl           = "{{ . | cyan }} "
	defaultConfirmPromptTpl     = "{{ . | bold | blue }} "
	defaultConfirmInvalidTpl    = "{{ . | red }} "
	defaultConfirmErrorMsgTpl   = "{{ . | bold | bgRed }} "
	defaultConfirmSuccessMsgTpl = "{{ . | bold | blue }}"
)

// the regex for confirm answers
var (
	yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")
	noRx  = regexp.MustCompile("^(?i:n(?:o)?)$")
)

// Confirm is a regular text input that accept yes/no answers. Response type is a bool.
type Confirm struct {
	ConfirmConfig
	Message    string
	Prompt     string
	isFirstRun bool

	FuncMap template.FuncMap

	// Default is the initial value for the confirm
	Default bool

	stdin  io.ReadCloser
	stdout io.WriteCloser

	message  *template.Template
	prompt   *template.Template
	invalid  *template.Template
	errorMsg *template.Template
	success  *template.Template
}

type ConfirmConfig struct {
	MessageTpl    string
	PromptTpl     string
	InvalidTpl    string
	ErrorMsgTpl   string
	SuccessMsgTpl string
}

func defaultConfig() ConfirmConfig {
	return ConfirmConfig{
		MessageTpl:    defaultMessageTpl,
		PromptTpl:     defaultConfirmPromptTpl,
		InvalidTpl:    defaultConfirmInvalidTpl,
		ErrorMsgTpl:   defaultConfirmErrorMsgTpl,
		SuccessMsgTpl: defaultConfirmSuccessMsgTpl,
	}
}

func NewDefaultConfirm(message string, defaultConfirm bool) Confirm {
	return Confirm{
		Message:       message,
		Prompt:        DefaultConfirmPrompt,
		Default:       defaultConfirm,
		FuncMap:       FuncMap,
		ConfirmConfig: defaultConfig(),
	}
}

func (c *Confirm) prepareTemplates() error {

	var err error
	if c.message, err = template.New("").Funcs(FuncMap).Parse(c.MessageTpl); err != nil {
		return err
	}

	if c.prompt, err = template.New("").Funcs(FuncMap).Parse(c.PromptTpl); err != nil {
		return err
	}

	if c.invalid, err = template.New("").Funcs(FuncMap).Parse(c.InvalidTpl); err != nil {
		return err
	}

	if c.errorMsg, err = template.New("").Funcs(FuncMap).Parse(c.ErrorMsgTpl); err != nil {
		return err
	}

	if c.success, err = template.New("").Funcs(FuncMap).Parse(c.SuccessMsgTpl); err != nil {
		return err
	}
	return nil
}

func (c *Confirm) Run() (bool, error) {

	var (
		err           error
		inputErr      error
		invalidPrompt []byte
		displayPrompt []byte
		validPrompt   []byte
	)
	c.isFirstRun = true
	if err = c.prepareTemplates(); err != nil {
		return false, err
	}

	displayPrompt = append(utils.Render(c.prompt, c.Prompt), utils.Render(c.message, c.Message+" ["+yesNoPrompt+"]")...)
	validPrompt = append(utils.Render(c.prompt, c.Prompt), utils.Render(c.message, c.Message+" ["+yesNoPrompt+"]")...)
	invalidPrompt = append(utils.Render(c.invalid, c.Prompt), utils.Render(c.message, c.Message+" ["+yesNoPrompt+"]")...)

	cf := &readline.Config{
		Prompt:          string(displayPrompt),
		Stdin:           c.stdin,
		Stdout:          c.stdout,
		HistoryLimit:    -1,
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
	}

	if err = cf.Init(); err != nil {
		return false, err
	}

	l, err := readline.NewEx(cf)
	if err != nil {
		return false, err
	}

	filterInput := func(r rune) (rune, bool) {

		switch r {
		// block CtrlZ feature
		case readline.CharCtrlZ:
			return r, false
		default:
			return r, true
		}
	}

	l.Config.FuncFilterInputRune = filterInput

	listen := func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		// Real-time verification
		_, err := c.getAnswer(string(line))
		if err != nil {
			l.SetPrompt(string(invalidPrompt))
		} else {
			l.SetPrompt(string(validPrompt))
		}
		l.Refresh()
		return nil, 0, false
	}

	l.Config.SetListener(listen)

	defer l.Close()

	var readErr error
	var s string
	// set default confirm value
	_, _ = l.WriteStdin([]byte(yesNo(c.Default)))

	for {

		if !c.isFirstRun {
			_, _ = l.Write([]byte(moveUp))
		}
		if s, readErr = l.Readline(); readErr != nil {
			break
		}

		if _, inputErr = c.getAnswer(s); inputErr != nil {
			_, _ = l.Write([]byte(moveDown))
			_, _ = l.Write([]byte(clearLine))
			_, _ = l.Write([]byte(string(utils.Render(c.errorMsg, DefaultErrorMsgPrefix+" "+inputErr.Error()))))
			c.isFirstRun = false
		} else {
			break
		}
	}

	if readErr != nil {

		switch readErr {
		case readline.ErrInterrupt:
			readErr = utils.ErrInterrupt
		case io.EOF:
			readErr = utils.ErrEOF
		}
		if readErr.Error() == "Interrupt" {
			readErr = utils.ErrInterrupt
		}

		return false, readErr
	}

	answer, _ := c.getAnswer(s)
	return answer, err
}

func yesNo(t bool) string {
	if t {
		return "y"
	}
	return "N"
}
func (c *Confirm) getAnswer(input string) (bool, error) {

	var answer bool
	switch {
	case yesRx.Match([]byte(input)):
		answer = true
	case noRx.Match([]byte(input)):
		answer = false
	default:
		err := errors.New("Not valid answer,please try again.")
		return false, err

	}
	return answer, nil
}
