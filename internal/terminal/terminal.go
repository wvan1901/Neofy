package terminal

import (
	"bufio"
	"fmt"
	"neofy/internal/consts"
	"os"
	"os/exec"
	"runtime"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
)

func ReaderBytes() []byte {
	STDINFILE := os.Stdin
	reader := bufio.NewReader(STDINFILE)
	inputBytes := make([]byte, 3)
	reader.Read(inputBytes)
	return inputBytes
}

func ReadInputKey() rune {
	inputBytes := ReaderBytes()
	inputRune, _ := utf8.DecodeRune(inputBytes)
	if unicode.IsControl(inputRune) {
		switch inputRune {
		case 3: //CTRL-C
			return consts.CONTROLCASCII
		case 27: //First byte is A CTRL byte
			returnRune, _ := utf8.DecodeRune(inputBytes[2:])
			switch returnRune {
			case 53: //PAGE UP
				return consts.PAGE_UP
			case 54: //PAGE DOWN
				return consts.PAGE_DOWN
			case 68: //LEFT ARROW
				return consts.LEFT_ARROW
			case 67: //RIGHT ARROW
				return consts.RIGHT_ARROW
			case 66: //DOWN ARROW
				return consts.DOWN_ARROW
			case 65: //UP ARROW
				return consts.UP_ARROW
			case 72: //HOME KEY
				return consts.HOME_KEY
			case 70: //END KEY
				return consts.END_KEY
			case 51: //DEL KEY
				return consts.DEL_KEY
			case 0: //ESC KEY
				return consts.ESC
			}
			return consts.NOTHINGKEY
		case 127: //BACKSPACE
			return consts.BACKSPACE
		case 13: //ENTER
			return '\r'
		case 19: //CTRL-S
			return consts.CONTROL_S
		case 6: //CTRL-F
			return consts.CONTROL_F
		case 4: //CTRL-D
			return consts.CONTROL_D
		case 21: //CTRL-U
			return consts.CONTROL_U
		default:
			return consts.NOTHINGKEY
		}
	}
	return inputRune
}

func Quit(l AppTerm) {
	l.CloseTerminal()
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	defer os.Exit(0)
}

type AppTerm interface {
	InitTerminal()
	GetTerminalSize() (int, int, error)
	CloseTerminal()
}

type LinuxTerm struct {
	OldState *term.State
}

func (l *LinuxTerm) InitTerminal() {
	l.enableRawMode()
}

func (l *LinuxTerm) GetTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(0)
	if err != nil {
		return 0, 0, fmt.Errorf("GetTerminalSize: %w", err)
	}
	return width, height, nil
}

func (l *LinuxTerm) CloseTerminal() {
	l.disableRawMode()
}

func (l *LinuxTerm) enableRawMode() {
	oldState, err := term.MakeRaw(0)
	if err != nil {
		panic(fmt.Errorf("enableRawMode: %w", err))
	}
	l.OldState = oldState
}

func (l *LinuxTerm) disableRawMode() {
	err := term.Restore(0, l.OldState)
	if err != nil {
		panic(fmt.Errorf("disableRawMode: %w", err))
	}
}

func Openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return fmt.Errorf("OpenBrowser: %w", err)
	}
	return nil
}

// TODO: Handle different environments
func InitAppTerm() AppTerm {
	newTerm := LinuxTerm{}
	newTerm.InitTerminal()
	return &newTerm
}
