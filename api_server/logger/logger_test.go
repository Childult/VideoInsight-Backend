package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	InitLog()

	Info.Println("hello")
	Warning.Println("hello")
	Debug.Println("hello")
	Error.Println("hello")

	// 颜色输出, 测试时好像失效, 在main里运行是有的. 目前确定是vscode OUTPUT窗口的问题, 将命令复制到终端运行时结果正常
	// fmt.Printf("\n %c[0;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	// fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	// fmt.Printf("\n %c[4;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	// fmt.Printf("\n %c[5;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	// fmt.Printf("\n %c[7;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	// fmt.Printf("\n %c[8;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)

	// for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
	// 	for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
	// 		for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
	// 			fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
	// 		}
	// 		fmt.Println("")
	// 	}
	// 	fmt.Println("")
	// }

}
