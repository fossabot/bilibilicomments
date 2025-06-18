package common

import "fmt"

func RgbToAnsi(color int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", (color>>16)&0xFF, (color>>8)&0xFF, color&0xFF) // ANSI escape code for 24-bit RGB
}

func ResetColor() string {
	return "\033[0m"
}
