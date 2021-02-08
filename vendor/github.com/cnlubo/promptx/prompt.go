package promptx

import (
	"github.com/cnlubo/promptx/utils"
	"github.com/mritd/readline"
	"io"
	"text/template"
)

const (
	defaultAskTpl        = "{{ . | cyan }} "
	DefaultPromptTpl     = "{{ . | green }} "
	defaultInvalidTpl    = "{{ . | red }} "
	defaultValidTpl      = "{{ . | blue }} "
	defaultErrorMsgTpl   = "{{ . | bold | bgRed }} "
	defaultSuccessMsgTpl = "{{ . | faint }}"
)

type Prompt struct {
	Config
	Ask        string
	Prompt     string
	isFirstRun bool
	FuncMap    template.FuncMap

	// Default is the initial value for the prompt
	Default string

	// Mask is an optional rune that sets which character to display instead of the entered characters. This
	// allows hiding private information like passwords.
	Mask rune

	stdin  io.ReadCloser
	stdout io.WriteCloser

	ask      *template.Template
	prompt   *template.Template
	valid    *template.Template
	invalid  *template.Template
	errorMsg *template.Template
	success  *template.Template
}

type Config struct {
	AskTpl        string
	PromptTpl     string
	ValidTpl      string
	InvalidTpl    string
	ErrorMsgTpl   string
	SuccessMsgTpl string
	CheckListener func(line []rune) error
}

func newDefaultConfig(check func(line []rune) error) Config {
	return Config{
		AskTpl:        defaultAskTpl,
		PromptTpl:     DefaultPromptTpl,
		InvalidTpl:    defaultInvalidTpl,
		ValidTpl:      defaultValidTpl,
		ErrorMsgTpl:   defaultErrorMsgTpl,
		SuccessMsgTpl: defaultSuccessMsgTpl,
		CheckListener: check,
	}
}

func NewDefaultPrompt(check func(line []rune) error, ask string) Prompt {
	return Prompt{
		Ask:     ask,
		Prompt:  DefaultPrompt,
		FuncMap: FuncMap,
		Config:  newDefaultConfig(check),
	}
}

func (p *Prompt) prepareTemplates() error {

	var err error
	if p.ask, err = template.New("").Funcs(FuncMap).Parse(p.AskTpl); err != nil {
		return err
	}

	if p.prompt, err = template.New("").Funcs(FuncMap).Parse(p.PromptTpl); err != nil {
		return err
	}

	if p.valid, err = template.New("").Funcs(FuncMap).Parse(p.ValidTpl); err != nil {
		return err
	}

	if p.invalid, err = template.New("").Funcs(FuncMap).Parse(p.InvalidTpl); err != nil {
		return err
	}

	if p.errorMsg, err = template.New("").Funcs(FuncMap).Parse(p.ErrorMsgTpl); err != nil {
		return err
	}

	if p.success, err = template.New("").Funcs(FuncMap).Parse(p.SuccessMsgTpl); err != nil {
		return err
	}
	return nil
}

func (p *Prompt) Run() (string, error) {

	var (
		err           error
		inputErr      error
		invalidPrompt []byte
		validPrompt   []byte
		displayPrompt []byte
	)
	p.isFirstRun = true
	if err = p.prepareTemplates(); err != nil {
		return "", err
	}
	displayPrompt = append(utils.Render(p.prompt, p.Prompt), utils.Render(p.ask, p.Ask)...)
	validPrompt = append(utils.Render(p.valid, p.Prompt), utils.Render(p.ask, p.Ask)...)
	invalidPrompt = append(utils.Render(p.invalid, p.Prompt), utils.Render(p.ask, p.Ask)...)

	cf := &readline.Config{
		Prompt:          string(displayPrompt),
		Stdin:           p.stdin,
		Stdout:          p.stdout,
		EnableMask:      p.Mask != 0,
		MaskRune:        p.Mask,
		HistoryLimit:    -1,
		UniqueEditLine:  true,
		InterruptPrompt: "^C",
	}

	if err = cf.Init(); err != nil {
		return "", err
	}

	l, err := readline.NewEx(cf)
	if err != nil {
		return "", err
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
		if err = p.CheckListener(line); err != nil {
			l.SetPrompt(string(invalidPrompt))
			l.Refresh()
		} else {
			l.SetPrompt(string(validPrompt))
			l.Refresh()
		}
		return nil, 0, false
	}

	l.Config.SetListener(listen)

	defer l.Close()

	var readErr error
	var s string
	if len(p.Default) != 0 {
		_, _ = l.WriteStdin([]byte(p.Default))
	}
	for {

		if !p.isFirstRun {
			_, _ = l.Write([]byte(moveUp))
		}
		if s, readErr = l.Readline(); readErr != nil {
			break
		}

		if inputErr = p.CheckListener([]rune(s)); inputErr != nil {
			_, _ = l.Write([]byte(moveDown))
			_, _ = l.Write([]byte(clearLine))
			_, _ = l.Write([]byte(string(utils.Render(p.errorMsg, DefaultErrorMsgPrefix+" "+inputErr.Error()))))
			p.isFirstRun = false
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
		return "", readErr
	}

	return s, nil
}

// func (p *Prompt) Run() (string, error) {
//
// 	var (
// 		err           error
// 		inputErr      error
// 		invalidPrompt []byte
// 		validPrompt   []byte
// 		successPrompt []byte
// 	)
// 	if err = p.prepareTemplates(); err != nil {
// 		return "", err
// 	}
//
// 	validPrompt = append(utils.Render(p.valid, p.Prompt), utils.Render(p.ask, p.Ask)...)
// 	invalidPrompt = append(utils.Render(p.invalid, p.Prompt), utils.Render(p.ask, p.Ask)...)
// 	successPrompt = append(utils.Render(p.success, DefaultGoodPrompt), utils.Render(p.ask, p.Ask)...)
//
// 	cf := &readline.Config{
// 		Stdin:           p.stdin,
// 		Stdout:          p.stdout,
// 		EnableMask:      p.Mask != 0,
// 		MaskRune:        p.Mask,
// 		HistoryLimit:    -1,
// 		UniqueEditLine:  true,
// 		InterruptPrompt: "^C",
// 	}
//
// 	if err = cf.Init(); err != nil {
// 		return "", err
// 	}
//
// 	l, err := readline.NewEx(cf)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	filterInput := func(r rune) (rune, bool) {
//
// 		switch r {
// 		// block CtrlZ feature
// 		case readline.CharCtrlZ:
// 			return r, false
// 		default:
// 			return r, true
// 		}
// 	}
//
// 	l.Config.FuncFilterInputRune = filterInput
//
// 	_, _ = l.Write([]byte(hideCursor))
// 	sb := terminal.New(l)
// 	input := p.Default
// 	cur := terminal.NewCursor(input, nil, false)
//
// 	listen := func(input []rune, pos int, key rune) ([]rune, int, bool) {
// 		_, _, keepOn := cur.Listen(input, pos, key)
// 		var prompt []byte
// 		err := p.CheckListener([]rune(cur.Get()))
// 		if err != nil {
// 			prompt = invalidPrompt
// 		} else {
// 			prompt = validPrompt
// 		}
//
// 		echo := cur.Format()
// 		if p.Mask != 0 {
// 			echo = cur.FormatMask(p.Mask)
// 		}
//
// 		prompt = append(prompt, []byte(echo)...)
// 		sb.Reset()
//
// 		_, _ = sb.Write(prompt)
//
// 		if inputErr != nil {
// 			validation := []byte(string(utils.Render(p.errorMsg, DefaultErrorMsgPrefix+" "+inputErr.Error())))
// 			_, _ = sb.Write(validation)
// 			inputErr = nil
// 		}
// 		_ = sb.Flush()
// 		return nil, 0, keepOn
// 	}
//
// 	l.Config.SetListener(listen)
//
// 	defer l.Close()
//
// 	var readErr error
//
// 	for {
//
// 		if _, readErr = l.Readline(); readErr != nil {
// 			break
// 		}
// 		if inputErr = p.CheckListener([]rune(cur.Get())); inputErr == nil {
// 			break
// 		}
// 	}
//
// 	if readErr != nil {
//
// 		switch readErr {
// 		case readline.ErrInterrupt:
// 			readErr = utils.ErrInterrupt
// 		case io.EOF:
// 			readErr = utils.ErrEOF
// 		}
// 		if readErr.Error() == "Interrupt" {
// 			readErr = utils.ErrInterrupt
// 		}
// 		sb.Reset()
// 		_, _ = sb.WriteString("")
// 		_ = sb.Flush()
// 		_, _ = l.Write([]byte(showCursor))
// 		return "", readErr
// 	}
//
// 	echo := cur.Format()
// 	if p.Mask != 0 {
// 		echo = cur.FormatMask(p.Mask)
// 	}
//
// 	prompt := successPrompt
// 	prompt = append(prompt, []byte(echo)...)
//
// 	sb.Reset()
// 	_, _ = sb.Write(prompt)
// 	_ = sb.Flush()
// 	_, _ = l.Write([]byte(showCursor))
//
// 	return cur.Get(), err
// }
