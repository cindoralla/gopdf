package gopdf

import (
	"testing"

	"github.com/cindoralla/gopdf/core"
)

func ComplexHLineReport() {
	r := core.CreateReport()
	r.SetPage("A4", "P")

	r.RegisterExecutor(core.Executor(ComplexHLineReportExecutor), core.Detail)

	r.Execute("hr_test.pdf")
	r.SaveAtomicCellText("hr_test.txt")
}
func ComplexHLineReportExecutor(report *core.Report) {
	unit := 2.83
	hr := NewHLine(report)
	hr.SetColor(0.7).SetWidth(5 * unit)
	hr.GenerateAtomicCell()
}

func TestComplexHLineReport(t *testing.T) {
	ComplexHLineReport()
}

type A interface {
	A()
}

type B interface {
	B()
}

type T struct {
}

func (t *T) A() {

}
func (t *T) B() {

}

func TestAB(t *testing.T) {
	var b B
	b = &T{}
	if _, ok := b.(A); ok {
		t.Log("ok")
	}
}
