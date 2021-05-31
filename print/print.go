package print

import "fmt"

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见
func ColorPrint(foreColor int, backColor int, content string) {
	// fmt.Printf("\x1b[%dmhello world 30: 黑 \x1b[0m\n", 30)
	// fmt.Printf("\x1b[%dmhello world 31: 红 \x1b[0m\n", 31)
	// fmt.Printf("\x1b[%dmhello world 32: 绿 \x1b[0m\n", 32)
	// fmt.Printf("\x1b[%dmhello world 33: 黄 \x1b[0m\n", 33)
	// fmt.Printf("\x1b[%dmhello world 34: 蓝 \x1b[0m\n", 34)
	// fmt.Printf("\x1b[%dmhello world 35: 紫 \x1b[0m\n", 35)
	// fmt.Printf("\x1b[%dmhello world 36: 青 \x1b[0m\n", 36)
	// fmt.Printf("\x1b[%dmhello world 37: 白 \x1b[0m\n", 37)

	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 47: 白   30: 黑 \n", 47, 30)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 46: 青   31: 红 \n", 46, 31)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 45: 紫   32: 绿 \n", 45, 32)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 44: 蓝   33: 黄 \n", 44, 33)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 43: 黄   34: 蓝 \n", 43, 34)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 42: 绿   35: 紫 \n", 42, 35)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 41: 红   36: 青 \n", 41, 36)
	// fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 40: 黑   37: 白 \n", 40, 37)

	if backColor == 0 {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", foreColor, content)
	} else {
		fmt.Printf("\x1b[%d;%dm%s\x1b[0m\n", foreColor, backColor, content)
	}
}
