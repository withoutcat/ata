package gui

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/withoutcat/ata/internal/converter"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

// StartGUI 启动GUI界面
func StartGUI() {
	// 创建应用程序
	vcl.Application.Initialize()
	vcl.Application.SetTitle("ATA - AVIF图像转换工具")

	// 创建主窗口
	mainForm := vcl.Application.CreateForm()
	mainForm.SetCaption("ATA - AVIF图像转换工具")
	mainForm.SetPosition(types.PoScreenCenter)
	mainForm.SetWidth(600)
	mainForm.SetHeight(500)

	// 创建选项卡控件
	pageControl := vcl.NewPageControl(mainForm)
	pageControl.SetParent(mainForm)
	pageControl.SetAlign(types.AlClient)

	// 创建批量转换选项卡
	convertTab := vcl.NewTabSheet(mainForm)
	convertTab.SetPageControl(pageControl)
	convertTab.SetCaption("批量转换")

	// 创建动画合成选项卡
	aniTab := vcl.NewTabSheet(mainForm)
	aniTab.SetPageControl(pageControl)
	aniTab.SetCaption("动画合成")

	// 创建幻灯片选项卡
	pptTab := vcl.NewTabSheet(mainForm)
	pptTab.SetPageControl(pageControl)
	pptTab.SetCaption("幻灯片")

	// 设置批量转换选项卡内容
	setupConvertTab(convertTab)

	// 设置动画合成选项卡内容
	setupAniTab(aniTab)

	// 设置幻灯片选项卡内容
	setupPptTab(pptTab)

	// 运行应用程序
	vcl.Application.Run()
}

// 设置批量转换选项卡内容
func setupConvertTab(tab *vcl.TTabSheet) {
	// 创建面板
	panel := vcl.NewPanel(tab)
	panel.SetParent(tab)
	panel.SetAlign(types.AlClient)
	panel.SetBevelOuter(types.BvNone)

	// 创建输入路径标签和编辑框
	pathLabel := vcl.NewLabel(panel)
	pathLabel.SetParent(panel)
	pathLabel.SetCaption("输入路径:")
	pathLabel.SetLeft(20)
	pathLabel.SetTop(20)

	pathEdit := vcl.NewEdit(panel)
	pathEdit.SetParent(panel)
	pathEdit.SetLeft(100)
	pathEdit.SetTop(20)
	pathEdit.SetWidth(350)

	// 创建浏览按钮
	browseBtn := vcl.NewButton(panel)
	browseBtn.SetParent(panel)
	browseBtn.SetCaption("浏览...")
	browseBtn.SetLeft(460)
	browseBtn.SetTop(20)
	browseBtn.SetWidth(80)
	browseBtn.SetOnClick(func(sender vcl.IObject) {
		selDialog := vcl.NewSelectDirectoryDialog(panel)
		if selDialog.Execute() {
			pathEdit.SetText(selDialog.FileName())
		}
	})

	// 创建选项
	recursiveCheck := vcl.NewCheckBox(panel)
	recursiveCheck.SetParent(panel)
	recursiveCheck.SetCaption("递归处理子目录")
	recursiveCheck.SetLeft(20)
	recursiveCheck.SetTop(60)

	deleteOriginalCheck := vcl.NewCheckBox(panel)
	deleteOriginalCheck.SetParent(panel)
	deleteOriginalCheck.SetCaption("删除原始文件")
	deleteOriginalCheck.SetLeft(20)
	deleteOriginalCheck.SetTop(90)

	forceCheck := vcl.NewCheckBox(panel)
	forceCheck.SetParent(panel)
	forceCheck.SetCaption("强制覆盖已存在的文件")
	forceCheck.SetLeft(20)
	forceCheck.SetTop(120)

	debugCheck := vcl.NewCheckBox(panel)
	debugCheck.SetParent(panel)
	debugCheck.SetCaption("调试模式")
	debugCheck.SetLeft(20)
	debugCheck.SetTop(150)

	// 创建状态标签
	statusLabel := vcl.NewLabel(panel)
	statusLabel.SetParent(panel)
	statusLabel.SetCaption("")
	statusLabel.SetLeft(20)
	statusLabel.SetTop(220)
	statusLabel.SetWidth(500)

	// 创建转换按钮
	convertBtn := vcl.NewButton(panel)
	convertBtn.SetParent(panel)
	convertBtn.SetCaption("开始转换")
	convertBtn.SetLeft(20)
	convertBtn.SetTop(180)
	convertBtn.SetWidth(120)
	convertBtn.SetOnClick(func(sender vcl.IObject) {
		path := pathEdit.Text()
		if path == "" {
			statusLabel.SetCaption("错误: 请指定输入路径")
			return
		}

		// 获取选项
		recursive := recursiveCheck.Checked()
		deleteOriginal := deleteOriginalCheck.Checked()
		force := forceCheck.Checked()
		debugMode := debugCheck.Checked()

		// 禁用按钮，防止重复点击
		convertBtn.SetEnabled(false)
		statusLabel.SetCaption("正在转换...")

		// 在后台线程中执行转换
		go func() {
			// 捕获标准输出和标准错误
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// 在后台读取输出
			outputChan := make(chan string, 10)
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if n > 0 {
						outputChan <- string(buf[:n])
					}
					if err != nil {
						break
					}
				}
				close(outputChan)
			}()

			// 执行转换
			converter.ConvertImages(path, debugMode, deleteOriginal, recursive, force)

			// 恢复标准输出和标准错误
			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// 收集输出
			output := ""
			for line := range outputChan {
				output += line
			}

			// 更新UI（必须在主线程中执行）
			vcl.ThreadSync(func() {
				convertBtn.SetEnabled(true)
				statusLabel.SetCaption("转换完成")

				// 显示结果对话框
				resultForm := vcl.NewForm(nil)
				resultForm.SetCaption("转换结果")
				resultForm.SetWidth(500)
				resultForm.SetHeight(400)
				resultForm.SetPosition(types.PoScreenCenter)

				memo := vcl.NewMemo(resultForm)
				memo.SetParent(resultForm)
				memo.SetAlign(types.AlClient)
				memo.SetScrollBars(types.SsBoth)
				memo.SetText(output)

				resultForm.Show()
			})
		}()
	})
}

// 设置动画合成选项卡内容
func setupAniTab(tab *vcl.TTabSheet) {
	// 创建面板
	panel := vcl.NewPanel(tab)
	panel.SetParent(tab)
	panel.SetAlign(types.AlClient)
	panel.SetBevelOuter(types.BvNone)

	// 创建输入路径标签和编辑框
	inputLabel := vcl.NewLabel(panel)
	inputLabel.SetParent(panel)
	inputLabel.SetCaption("输入路径:")
	inputLabel.SetLeft(20)
	inputLabel.SetTop(20)

	inputEdit := vcl.NewEdit(panel)
	inputEdit.SetParent(panel)
	inputEdit.SetLeft(100)
	inputEdit.SetTop(20)
	inputEdit.SetWidth(350)

	// 创建输入浏览按钮
	inputBrowseBtn := vcl.NewButton(panel)
	inputBrowseBtn.SetParent(panel)
	inputBrowseBtn.SetCaption("浏览...")
	inputBrowseBtn.SetLeft(460)
	inputBrowseBtn.SetTop(20)
	inputBrowseBtn.SetWidth(80)
	inputBrowseBtn.SetOnClick(func(sender vcl.IObject) {
		selDialog := vcl.NewSelectDirectoryDialog(panel)
		if selDialog.Execute() {
			inputEdit.SetText(selDialog.FileName())
		}
	})

	// 创建输出路径标签和编辑框
	outputLabel := vcl.NewLabel(panel)
	outputLabel.SetParent(panel)
	outputLabel.SetCaption("输出文件:")
	outputLabel.SetLeft(20)
	outputLabel.SetTop(50)

	outputEdit := vcl.NewEdit(panel)
	outputEdit.SetParent(panel)
	outputEdit.SetLeft(100)
	outputEdit.SetTop(50)
	outputEdit.SetWidth(350)

	// 创建输出浏览按钮
	outputBrowseBtn := vcl.NewButton(panel)
	outputBrowseBtn.SetParent(panel)
	outputBrowseBtn.SetCaption("浏览...")
	outputBrowseBtn.SetLeft(460)
	outputBrowseBtn.SetTop(50)
	outputBrowseBtn.SetWidth(80)
	outputBrowseBtn.SetOnClick(func(sender vcl.IObject) {
		saveDialog := vcl.NewSaveDialog(panel)
		saveDialog.SetTitle("保存AVIF动画")
		saveDialog.SetFilter("AVIF文件|*.avif")
		saveDialog.SetDefaultExt("avif")
		if saveDialog.Execute() {
			outputEdit.SetText(saveDialog.FileName())
		}
	})

	// 创建参数组
	paramGroup := vcl.NewGroupBox(panel)
	paramGroup.SetParent(panel)
	paramGroup.SetCaption("参数设置")
	paramGroup.SetLeft(20)
	paramGroup.SetTop(90)
	paramGroup.SetWidth(520)
	paramGroup.SetHeight(200)

	// 创建FPS标签和编辑框
	fpsLabel := vcl.NewLabel(paramGroup)
	fpsLabel.SetParent(paramGroup)
	fpsLabel.SetCaption("帧率:")
	fpsLabel.SetLeft(20)
	fpsLabel.SetTop(30)

	fpsEdit := vcl.NewEdit(paramGroup)
	fpsEdit.SetParent(paramGroup)
	fpsEdit.SetLeft(120)
	fpsEdit.SetTop(30)
	fpsEdit.SetWidth(80)
	fpsEdit.SetText("10")

	// 创建质量标签和编辑框
	crfLabel := vcl.NewLabel(paramGroup)
	crfLabel.SetParent(paramGroup)
	crfLabel.SetCaption("质量 (0-63):")
	crfLabel.SetLeft(20)
	crfLabel.SetTop(60)

	crfEdit := vcl.NewEdit(paramGroup)
	crfEdit.SetParent(paramGroup)
	crfEdit.SetLeft(120)
	crfEdit.SetTop(60)
	crfEdit.SetWidth(80)
	crfEdit.SetText("30")

	// 创建速度标签和编辑框
	speedLabel := vcl.NewLabel(paramGroup)
	speedLabel.SetParent(paramGroup)
	speedLabel.SetCaption("速度 (0-10):")
	speedLabel.SetLeft(20)
	speedLabel.SetTop(90)

	speedEdit := vcl.NewEdit(paramGroup)
	speedEdit.SetParent(paramGroup)
	speedEdit.SetLeft(120)
	speedEdit.SetTop(90)
	speedEdit.SetWidth(80)
	speedEdit.SetText("8")

	// 创建线程标签和编辑框
	threadsLabel := vcl.NewLabel(paramGroup)
	threadsLabel.SetParent(paramGroup)
	threadsLabel.SetCaption("线程数:")
	threadsLabel.SetLeft(20)
	threadsLabel.SetTop(120)

	threadsEdit := vcl.NewEdit(paramGroup)
	threadsEdit.SetParent(paramGroup)
	threadsEdit.SetLeft(120)
	threadsEdit.SetTop(120)
	threadsEdit.SetWidth(80)
	threadsEdit.SetText("0")

	// 创建宽度标签和编辑框
	widthLabel := vcl.NewLabel(paramGroup)
	widthLabel.SetParent(paramGroup)
	widthLabel.SetCaption("宽度:")
	widthLabel.SetLeft(220)
	widthLabel.SetTop(30)

	widthEdit := vcl.NewEdit(paramGroup)
	widthEdit.SetParent(paramGroup)
	widthEdit.SetLeft(320)
	widthEdit.SetTop(30)
	widthEdit.SetWidth(80)
	widthEdit.SetText("0")

	// 创建高度标签和编辑框
	heightLabel := vcl.NewLabel(paramGroup)
	heightLabel.SetParent(paramGroup)
	heightLabel.SetCaption("高度:")
	heightLabel.SetLeft(220)
	heightLabel.SetTop(60)

	heightEdit := vcl.NewEdit(paramGroup)
	heightEdit.SetParent(paramGroup)
	heightEdit.SetLeft(320)
	heightEdit.SetTop(60)
	heightEdit.SetWidth(80)
	heightEdit.SetText("0")

	// 创建缩放标签和编辑框
	scaleLabel := vcl.NewLabel(paramGroup)
	scaleLabel.SetParent(paramGroup)
	scaleLabel.SetCaption("缩放比例:")
	scaleLabel.SetLeft(220)
	scaleLabel.SetTop(90)

	scaleEdit := vcl.NewEdit(paramGroup)
	scaleEdit.SetParent(paramGroup)
	scaleEdit.SetLeft(320)
	scaleEdit.SetTop(90)
	scaleEdit.SetWidth(80)
	scaleEdit.SetText("1.0")

	// 创建背景标签和编辑框
	bgLabel := vcl.NewLabel(paramGroup)
	bgLabel.SetParent(paramGroup)
	bgLabel.SetCaption("背景颜色:")
	bgLabel.SetLeft(220)
	bgLabel.SetTop(120)

	bgEdit := vcl.NewEdit(paramGroup)
	bgEdit.SetParent(paramGroup)
	bgEdit.SetLeft(320)
	bgEdit.SetTop(120)
	bgEdit.SetWidth(80)
	bgEdit.SetText("white")

	// 创建选项
	alphaCheck := vcl.NewCheckBox(panel)
	alphaCheck.SetParent(panel)
	alphaCheck.SetCaption("保留透明通道")
	alphaCheck.SetLeft(20)
	alphaCheck.SetTop(300)

	deleteOriginalCheck := vcl.NewCheckBox(panel)
	deleteOriginalCheck.SetParent(panel)
	deleteOriginalCheck.SetCaption("删除原始文件")
	deleteOriginalCheck.SetLeft(150)
	deleteOriginalCheck.SetTop(300)

	forceCheck := vcl.NewCheckBox(panel)
	forceCheck.SetParent(panel)
	forceCheck.SetCaption("强制覆盖已存在的文件")
	forceCheck.SetLeft(280)
	forceCheck.SetTop(300)

	debugCheck := vcl.NewCheckBox(panel)
	debugCheck.SetParent(panel)
	debugCheck.SetCaption("调试模式")
	debugCheck.SetLeft(450)
	debugCheck.SetTop(300)

	// 创建状态标签
	statusLabel := vcl.NewLabel(panel)
	statusLabel.SetParent(panel)
	statusLabel.SetCaption("")
	statusLabel.SetLeft(20)
	statusLabel.SetTop(370)
	statusLabel.SetWidth(500)

	// 创建创建按钮
	createBtn := vcl.NewButton(panel)
	createBtn.SetParent(panel)
	createBtn.SetCaption("创建动画")
	createBtn.SetLeft(20)
	createBtn.SetTop(330)
	createBtn.SetWidth(120)
	createBtn.SetOnClick(func(sender vcl.IObject) {
		inputPath := inputEdit.Text()
		if inputPath == "" {
			statusLabel.SetCaption("错误: 请指定输入路径")
			return
		}

		outputPath := outputEdit.Text()
		if outputPath == "" {
			// 默认输出路径为输入目录下的output.avif
			outputPath = filepath.Join(inputPath, "output.avif")
			outputEdit.SetText(outputPath)
		}

		// 获取参数
		fps, err := strconv.Atoi(fpsEdit.Text())
		if err != nil || fps <= 0 {
			statusLabel.SetCaption("错误: 帧率必须是正整数")
			return
		}

		crf, err := strconv.Atoi(crfEdit.Text())
		if err != nil || crf < 0 || crf > 63 {
			statusLabel.SetCaption("错误: 质量必须在0-63之间")
			return
		}

		speed, err := strconv.Atoi(speedEdit.Text())
		if err != nil || speed < 0 || speed > 10 {
			statusLabel.SetCaption("错误: 速度必须在0-10之间")
			return
		}

		threads, err := strconv.Atoi(threadsEdit.Text())
		if err != nil || threads < 0 {
			statusLabel.SetCaption("错误: 线程数必须是非负整数")
			return
		}

		width, err := strconv.Atoi(widthEdit.Text())
		if err != nil || width < 0 {
			statusLabel.SetCaption("错误: 宽度必须是非负整数")
			return
		}

		height, err := strconv.Atoi(heightEdit.Text())
		if err != nil || height < 0 {
			statusLabel.SetCaption("错误: 高度必须是非负整数")
			return
		}

		scale, err := strconv.ParseFloat(scaleEdit.Text(), 64)
		if err != nil || scale <= 0 {
			statusLabel.SetCaption("错误: 缩放比例必须是正数")
			return
		}

		background := bgEdit.Text()
		if background == "" {
			background = "white"
		}

		alpha := alphaCheck.Checked()
		deleteOriginal := deleteOriginalCheck.Checked()
		force := forceCheck.Checked()
		debugMode := debugCheck.Checked()

		// 禁用按钮，防止重复点击
		createBtn.SetEnabled(false)
		statusLabel.SetCaption("正在创建动画...")

		// 在后台线程中执行创建
		go func() {
			// 捕获标准输出和标准错误
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// 在后台读取输出
			outputChan := make(chan string, 10)
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if n > 0 {
						outputChan <- string(buf[:n])
					}
					if err != nil {
						break
					}
				}
				close(outputChan)
			}()

			// 执行创建
			converter.CreateAnimation(inputPath, outputPath, fps, crf, speed, threads, alpha, width, height, scale, background, debugMode, deleteOriginal, force)

			// 恢复标准输出和标准错误
			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// 收集输出
			output := ""
			for line := range outputChan {
				output += line
			}

			// 更新UI（必须在主线程中执行）
			vcl.ThreadSync(func() {
				createBtn.SetEnabled(true)
				statusLabel.SetCaption("动画创建完成")

				// 显示结果对话框
				resultForm := vcl.NewForm(nil)
				resultForm.SetCaption("创建结果")
				resultForm.SetWidth(500)
				resultForm.SetHeight(400)
				resultForm.SetPosition(types.PoScreenCenter)

				memo := vcl.NewMemo(resultForm)
				memo.SetParent(resultForm)
				memo.SetAlign(types.AlClient)
				memo.SetScrollBars(types.SsBoth)
				memo.SetText(output)

				resultForm.Show()
			})
		}()
	})
}

// 设置幻灯片选项卡内容
func setupPptTab(tab *vcl.TTabSheet) {
	// 创建面板
	panel := vcl.NewPanel(tab)
	panel.SetParent(tab)
	panel.SetAlign(types.AlClient)
	panel.SetBevelOuter(types.BvNone)

	// 创建输入路径标签和编辑框
	inputLabel := vcl.NewLabel(panel)
	inputLabel.SetParent(panel)
	inputLabel.SetCaption("输入路径:")
	inputLabel.SetLeft(20)
	inputLabel.SetTop(20)

	inputEdit := vcl.NewEdit(panel)
	inputEdit.SetParent(panel)
	inputEdit.SetLeft(100)
	inputEdit.SetTop(20)
	inputEdit.SetWidth(350)

	// 创建输入浏览按钮
	inputBrowseBtn := vcl.NewButton(panel)
	inputBrowseBtn.SetParent(panel)
	inputBrowseBtn.SetCaption("浏览...")
	inputBrowseBtn.SetLeft(460)
	inputBrowseBtn.SetTop(20)
	inputBrowseBtn.SetWidth(80)
	inputBrowseBtn.SetOnClick(func(sender vcl.IObject) {
		selDialog := vcl.NewSelectDirectoryDialog(panel)
		if selDialog.Execute() {
			inputEdit.SetText(selDialog.FileName())
		}
	})

	// 创建输出路径标签和编辑框
	outputLabel := vcl.NewLabel(panel)
	outputLabel.SetParent(panel)
	outputLabel.SetCaption("输出文件:")
	outputLabel.SetLeft(20)
	outputLabel.SetTop(50)

	outputEdit := vcl.NewEdit(panel)
	outputEdit.SetParent(panel)
	outputEdit.SetLeft(100)
	outputEdit.SetTop(50)
	outputEdit.SetWidth(350)

	// 创建输出浏览按钮
	outputBrowseBtn := vcl.NewButton(panel)
	outputBrowseBtn.SetParent(panel)
	outputBrowseBtn.SetCaption("浏览...")
	outputBrowseBtn.SetLeft(460)
	outputBrowseBtn.SetTop(50)
	outputBrowseBtn.SetWidth(80)
	outputBrowseBtn.SetOnClick(func(sender vcl.IObject) {
		saveDialog := vcl.NewSaveDialog(panel)
		saveDialog.SetTitle("保存AVIF幻灯片")
		saveDialog.SetFilter("AVIF文件|*.avif")
		saveDialog.SetDefaultExt("avif")
		if saveDialog.Execute() {
			outputEdit.SetText(saveDialog.FileName())
		}
	})

	// 创建参数组
	paramGroup := vcl.NewGroupBox(panel)
	paramGroup.SetParent(panel)
	paramGroup.SetCaption("参数设置")
	paramGroup.SetLeft(20)
	paramGroup.SetTop(90)
	paramGroup.SetWidth(520)
	paramGroup.SetHeight(200)

	// 创建FPS标签和编辑框
	fpsLabel := vcl.NewLabel(paramGroup)
	fpsLabel.SetParent(paramGroup)
	fpsLabel.SetCaption("帧率:")
	fpsLabel.SetLeft(20)
	fpsLabel.SetTop(30)

	fpsEdit := vcl.NewEdit(paramGroup)
	fpsEdit.SetParent(paramGroup)
	fpsEdit.SetLeft(120)
	fpsEdit.SetTop(30)
	fpsEdit.SetWidth(80)
	fpsEdit.SetText("1")

	// 创建质量标签和编辑框
	crfLabel := vcl.NewLabel(paramGroup)
	crfLabel.SetParent(paramGroup)
	crfLabel.SetCaption("质量 (0-63):")
	crfLabel.SetLeft(20)
	crfLabel.SetTop(60)

	crfEdit := vcl.NewEdit(paramGroup)
	crfEdit.SetParent(paramGroup)
	crfEdit.SetLeft(120)
	crfEdit.SetTop(60)
	crfEdit.SetWidth(80)
	crfEdit.SetText("30")

	// 创建速度标签和编辑框
	speedLabel := vcl.NewLabel(paramGroup)
	speedLabel.SetParent(paramGroup)
	speedLabel.SetCaption("速度 (0-10):")
	speedLabel.SetLeft(20)
	speedLabel.SetTop(90)

	speedEdit := vcl.NewEdit(paramGroup)
	speedEdit.SetParent(paramGroup)
	speedEdit.SetLeft(120)
	speedEdit.SetTop(90)
	speedEdit.SetWidth(80)
	speedEdit.SetText("8")

	// 创建线程标签和编辑框
	threadsLabel := vcl.NewLabel(paramGroup)
	threadsLabel.SetParent(paramGroup)
	threadsLabel.SetCaption("线程数:")
	threadsLabel.SetLeft(20)
	threadsLabel.SetTop(120)

	threadsEdit := vcl.NewEdit(paramGroup)
	threadsEdit.SetParent(paramGroup)
	threadsEdit.SetLeft(120)
	threadsEdit.SetTop(120)
	threadsEdit.SetWidth(80)
	threadsEdit.SetText("0")

	// 创建宽度标签和编辑框
	widthLabel := vcl.NewLabel(paramGroup)
	widthLabel.SetParent(paramGroup)
	widthLabel.SetCaption("宽度:")
	widthLabel.SetLeft(220)
	widthLabel.SetTop(30)

	widthEdit := vcl.NewEdit(paramGroup)
	widthEdit.SetParent(paramGroup)
	widthEdit.SetLeft(320)
	widthEdit.SetTop(30)
	widthEdit.SetWidth(80)
	widthEdit.SetText("0")

	// 创建高度标签和编辑框
	heightLabel := vcl.NewLabel(paramGroup)
	heightLabel.SetParent(paramGroup)
	heightLabel.SetCaption("高度:")
	heightLabel.SetLeft(220)
	heightLabel.SetTop(60)

	heightEdit := vcl.NewEdit(paramGroup)
	heightEdit.SetParent(paramGroup)
	heightEdit.SetLeft(320)
	heightEdit.SetTop(60)
	heightEdit.SetWidth(80)
	heightEdit.SetText("0")

	// 创建缩放标签和编辑框
	scaleLabel := vcl.NewLabel(paramGroup)
	scaleLabel.SetParent(paramGroup)
	scaleLabel.SetCaption("缩放比例:")
	scaleLabel.SetLeft(220)
	scaleLabel.SetTop(90)

	scaleEdit := vcl.NewEdit(paramGroup)
	scaleEdit.SetParent(paramGroup)
	scaleEdit.SetLeft(320)
	scaleEdit.SetTop(90)
	scaleEdit.SetWidth(80)
	scaleEdit.SetText("1.0")

	// 创建背景标签和编辑框
	bgLabel := vcl.NewLabel(paramGroup)
	bgLabel.SetParent(paramGroup)
	bgLabel.SetCaption("背景颜色:")
	bgLabel.SetLeft(220)
	bgLabel.SetTop(120)

	bgEdit := vcl.NewEdit(paramGroup)
	bgEdit.SetParent(paramGroup)
	bgEdit.SetLeft(320)
	bgEdit.SetTop(120)
	bgEdit.SetWidth(80)
	bgEdit.SetText("white")

	// 创建选项
	alphaCheck := vcl.NewCheckBox(panel)
	alphaCheck.SetParent(panel)
	alphaCheck.SetCaption("保留透明通道")
	alphaCheck.SetLeft(20)
	alphaCheck.SetTop(300)

	deleteOriginalCheck := vcl.NewCheckBox(panel)
	deleteOriginalCheck.SetParent(panel)
	deleteOriginalCheck.SetCaption("删除原始文件")
	deleteOriginalCheck.SetLeft(150)
	deleteOriginalCheck.SetTop(300)

	forceCheck := vcl.NewCheckBox(panel)
	forceCheck.SetParent(panel)
	forceCheck.SetCaption("强制覆盖已存在的文件")
	forceCheck.SetLeft(280)
	forceCheck.SetTop(300)

	debugCheck := vcl.NewCheckBox(panel)
	debugCheck.SetParent(panel)
	debugCheck.SetCaption("调试模式")
	debugCheck.SetLeft(450)
	debugCheck.SetTop(300)

	// 创建状态标签
	statusLabel := vcl.NewLabel(panel)
	statusLabel.SetParent(panel)
	statusLabel.SetCaption("")
	statusLabel.SetLeft(20)
	statusLabel.SetTop(370)
	statusLabel.SetWidth(500)

	// 创建创建按钮
	createBtn := vcl.NewButton(panel)
	createBtn.SetParent(panel)
	createBtn.SetCaption("创建幻灯片")
	createBtn.SetLeft(20)
	createBtn.SetTop(330)
	createBtn.SetWidth(120)
	createBtn.SetOnClick(func(sender vcl.IObject) {
		inputPath := inputEdit.Text()
		if inputPath == "" {
			statusLabel.SetCaption("错误: 请指定输入路径")
			return
		}

		outputPath := outputEdit.Text()
		if outputPath == "" {
			// 默认输出路径为输入目录下的output.avif
			outputPath = filepath.Join(inputPath, "output.avif")
			outputEdit.SetText(outputPath)
		}

		// 获取参数
		fps, err := strconv.Atoi(fpsEdit.Text())
		if err != nil || fps <= 0 {
			statusLabel.SetCaption("错误: 帧率必须是正整数")
			return
		}

		crf, err := strconv.Atoi(crfEdit.Text())
		if err != nil || crf < 0 || crf > 63 {
			statusLabel.SetCaption("错误: 质量必须在0-63之间")
			return
		}

		speed, err := strconv.Atoi(speedEdit.Text())
		if err != nil || speed < 0 || speed > 10 {
			statusLabel.SetCaption("错误: 速度必须在0-10之间")
			return
		}

		threads, err := strconv.Atoi(threadsEdit.Text())
		if err != nil || threads < 0 {
			statusLabel.SetCaption("错误: 线程数必须是非负整数")
			return
		}

		width, err := strconv.Atoi(widthEdit.Text())
		if err != nil || width < 0 {
			statusLabel.SetCaption("错误: 宽度必须是非负整数")
			return
		}

		height, err := strconv.Atoi(heightEdit.Text())
		if err != nil || height < 0 {
			statusLabel.SetCaption("错误: 高度必须是非负整数")
			return
		}

		scale, err := strconv.ParseFloat(scaleEdit.Text(), 64)
		if err != nil || scale <= 0 {
			statusLabel.SetCaption("错误: 缩放比例必须是正数")
			return
		}

		background := bgEdit.Text()
		if background == "" {
			background = "white"
		}

		alpha := alphaCheck.Checked()
		deleteOriginal := deleteOriginalCheck.Checked()
		force := forceCheck.Checked()
		debugMode := debugCheck.Checked()

		// 禁用按钮，防止重复点击
		createBtn.SetEnabled(false)
		statusLabel.SetCaption("正在创建幻灯片...")

		// 在后台线程中执行创建
		go func() {
			// 捕获标准输出和标准错误
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// 在后台读取输出
			outputChan := make(chan string, 10)
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if n > 0 {
						outputChan <- string(buf[:n])
					}
					if err != nil {
						break
					}
				}
				close(outputChan)
			}()

			// 执行创建
			converter.CreateAnimation(inputPath, outputPath, fps, crf, speed, threads, alpha, width, height, scale, background, debugMode, deleteOriginal, force)

			// 恢复标准输出和标准错误
			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// 收集输出
			output := ""
			for line := range outputChan {
				output += line
			}

			// 更新UI（必须在主线程中执行）
			vcl.ThreadSync(func() {
				createBtn.SetEnabled(true)
				statusLabel.SetCaption("幻灯片创建完成")

				// 显示结果对话框
				resultForm := vcl.NewForm(nil)
				resultForm.SetCaption("创建结果")
				resultForm.SetWidth(500)
				resultForm.SetHeight(400)
				resultForm.SetPosition(types.PoScreenCenter)

				memo := vcl.NewMemo(resultForm)
				memo.SetParent(resultForm)
				memo.SetAlign(types.AlClient)
				memo.SetScrollBars(types.SsBoth)
				memo.SetText(output)

				resultForm.Show()
			})
		}()
	})
}