package example

import (
	"fmt"
	"github.com/cindoralla/gopdf/core"
	"testing"
)

func SimpleLink() {
	r := core.CreateReport()
	font1 := core.FontMap{
		FontName: FONT_MY,
		FileName: "ttf//microsoft.ttf",
	}
	font2 := core.FontMap{
		FontName: FONT_MD,
		FileName: "ttf//mplus-1p-bold.ttf",
	}
	r.SetFonts([]*core.FontMap{&font1, &font2})
	r.SetPage("A4", "P")
	r.FisrtPageNeedHeader = true
	r.FisrtPageNeedFooter = true

	r.RegisterExecutor(core.Executor(SimpleLinkExecutor), core.Detail)

	r.Execute(fmt.Sprintf("simple_link_test.pdf"))
	r.SaveAtomicCellText("simple_link_test.txt")
}

func SimpleLinkExecutor(report *core.Report) {
	var (
		lineHight = 16.0
	)
	report.SetFont(FONT_MY, 12)
	report.SetMargin(10, 20)
	x, y := report.GetXY()
	report.ExternalLink(x, y, lineHight, "外部链接(百度)", "https://www.baidu.com")
	report.InternalLinkAnchor(x, y+4*lineHight, lineHight, "内部链接", "1")

	report.AddNewPage(false)
	report.SetMargin(10, 20)
	x, y = report.GetXY()
	report.InternalLinkLink(x, y, "Hello world", "1")

}

func TestSimpleLink(t *testing.T) {
	SimpleLink()
}
