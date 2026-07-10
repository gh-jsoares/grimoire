package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Copy(text string) error {
	if text == "" {
		return nil
	}

	// Try OSC 52 if in tmux
	if os.Getenv("TMUX") != "" {
		return osc52Copy(text)
	}

	cmd, err := clipboardCmd()
	if err != nil {
		return err
	}

	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Stdin = strings.NewReader(text)
	return proc.Run()
}

func clipboardCmd() ([]string, error) {
	switch runtime.GOOS {
	case "darwin":
		return []string{"pbcopy"}, nil
	case "linux":
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			return []string{"wl-copy"}, nil
		}
		if _, err := exec.LookPath("xclip"); err == nil {
			return []string{"xclip", "-selection", "clipboard"}, nil
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return []string{"xsel", "--clipboard", "--input"}, nil
		}
		return nil, fmt.Errorf("no clipboard command found (install xclip or xsel)")
	case "windows":
		return []string{"clip.exe"}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func osc52Copy(text string) error {
	// OSC 52 clipboard sequence
	seq := fmt.Sprintf("\x1b]52;c;%s\x1b\\", encode(text))
	_, err := os.Stdout.WriteString(seq)
	return err
}

func encode(s string) string {
	// Base64 encode for OSC 52
	const encoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result []byte
	for i := 0; i < len(s); i += 3 {
		var b [3]byte
		n := 0
		for j := 0; j < 3 && i+j < len(s); j++ {
			b[j] = s[i+j]
			n++
		}
		result = append(result, encoding[b[0]>>2])
		result = append(result, encoding[((b[0]&0x03)<<4)|(b[1]>>4)])
		if n > 1 {
			result = append(result, encoding[((b[1]&0x0f)<<2)|(b[2]>>6)])
		} else {
			result = append(result, '=')
		}
		if n > 2 {
			result = append(result, encoding[b[2]&0x3f])
		} else {
			result = append(result, '=')
		}
	}
	return string(result)
}
