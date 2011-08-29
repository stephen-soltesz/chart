package txtg

import (
	"fmt"
	"math"
	"github.com/vdobler/chart"
)


// TextGraphics
type TextGraphics struct {
	tb   *TextBuf // the underlying text buffer
	w, h int      // width and height
	xoff int      // the initial radius for pie charts
}

// New creates a TextGraphic of dimensions w x h.
func New(w, h int) *TextGraphics {
	tg := TextGraphics{}
	tg.tb = NewTextBuf(w, h)
	tg.w, tg.h = w, h
	tg.xoff = -1
	return &tg
}


func (g *TextGraphics) Begin() {
	g.tb = NewTextBuf(g.w, g.h)
}

func (g *TextGraphics) End() {}
func (g *TextGraphics) Dimensions() (int, int) {
	return g.w, g.h
}
func (g *TextGraphics) FontMetrics(font chart.Font) (fw float32, fh int, mono bool) {
	return 1, 1, true
}

func (g *TextGraphics) TextLen(t string, font chart.Font) int {
	return len(t)
}


func (g *TextGraphics) Line(x0, y0, x1, y1 int, style chart.Style) {
	symbol := style.Symbol
	if symbol < ' ' || symbol > '~' {
		symbol = 'x'
	}
	g.tb.Line(x0, y0, x1, y1, symbol)
}

func (g *TextGraphics) Text(x, y int, t string, align string, rot int, font chart.Font) {
	// align: -1: left; 0: centered; 1: right; 2: top, 3: center, 4: bottom
	if len(align) == 2 {
		align = align[1:]
	}
	a := 0
	if rot == 0 {
		if align == "l" {
			a = -1
		}
		if align == "c" {
			a = 0
		}
		if align == "r" {
			a = 1
		}
	} else {
		if align == "l" {
			a = 2
		}
		if align == "c" {
			a = 3
		}
		if align == "r" {
			a = 4
		}
	}
	g.tb.Text(x, y, t, a)
}

func (g *TextGraphics) Rect(x, y, w, h int, style chart.Style) {
	// Normalize coordinates
	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}

	// Border
	if style.LineWidth > 0 {
		for i := 0; i < w; i++ {
			g.tb.Put(x+i, y, style.Symbol)
			g.tb.Put(x+i, y+h-1, style.Symbol)
		}
		for i := 1; i < h-1; i++ {
			g.tb.Put(x, y+i, style.Symbol)
			g.tb.Put(x+w-1, y+i, style.Symbol)
		}
	}

	// Filling
	if style.FillColor != "" {
		// TODO: fancier logic
		var s int
		if style.FillColor == "#000000" {
			s = '#' // black
		} else if style.FillColor == "#ffffff" {
			s = ' ' // white
		} else {
			s = style.Symbol
		}
		for i := 1; i < h-1; i++ {
			for j := 1; j < w-1; j++ {
				g.tb.Put(x+j, y+i, s)
			}
		}
	}
}

func (g *TextGraphics) Style(element string) chart.Style {
	b := "#000000"
	return chart.Style{Symbol: 'o', SymbolColor: b, LineColor: b, LineWidth: 1, LineStyle: chart.SolidLine}
}

func (g *TextGraphics) Font(element string) chart.Font {
	return chart.Font{}
}

func (g *TextGraphics) String() string {
	return g.tb.String()
}

func (g *TextGraphics) Symbol(x, y int, style chart.Style) {
	g.tb.Put(x, y, style.Symbol)
}
func (g *TextGraphics) Title(text string) {
	x, y := g.w/2, 1
	g.Text(x, y, text, "tc", 0, chart.Font{})
}

func (g *TextGraphics) XAxis(xrange chart.Range, y, y1 int) {
	mirror := xrange.TicSetting.Mirror
	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	for sx := xa; sx <= xe; sx++ {
		g.tb.Put(sx, y, '-')
		if mirror >= 1 {
			g.tb.Put(sx, y1, '-')
		}
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		for yy := y - 1; yy > y1+1; yy-- {
			g.tb.Put(z, yy, ':')
		}
	}

	if xrange.Label != "" {
		yy := y + 1
		if !xrange.TicSetting.Hide {
			yy++
		}
		g.tb.Text((xa+xe)/2, yy, xrange.Label, 0)
	}

	for _, tic := range xrange.Tics {
		var x int
		if !math.IsNaN(tic.Pos) {
			x = xrange.Data2Screen(tic.Pos)
		} else {
			x = -1
		}
		lx := xrange.Data2Screen(tic.LabelPos)
		if xrange.Time {
			if x != -1 {
				g.tb.Put(x, y, '|')
				if mirror >= 2 {
					g.tb.Put(x, y1, '|')
				}
				g.tb.Put(x, y+1, '|')
			}
			if tic.Align == -1 {
				g.tb.Text(lx+1, y+1, tic.Label, -1)
			} else {
				g.tb.Text(lx, y+1, tic.Label, 0)
			}
		} else {
			if x != -1 {
				g.tb.Put(x, y, '+')
				if mirror >= 2 {
					g.tb.Put(x, y1, '+')
				}
			}
			g.tb.Text(lx, y+1, tic.Label, 0)
		}
		if xrange.ShowLimits {
			if xrange.Time {
				g.tb.Text(xa, y+2, xrange.TMin.Format("2006-01-02 15:04:05"), -1)
				g.tb.Text(xe, y+2, xrange.TMax.Format("2006-01-02 15:04:05"), 1)
			} else {
				g.tb.Text(xa, y+2, fmt.Sprintf("%g", xrange.Min), -1)
				g.tb.Text(xe, y+2, fmt.Sprintf("%g", xrange.Max), 1)
			}
		}
	}
}

func (g *TextGraphics) YAxis(yrange chart.Range, x, x1 int) {
	label := yrange.Label
	mirror := yrange.TicSetting.Mirror
	ya, ye := yrange.Data2Screen(yrange.Min), yrange.Data2Screen(yrange.Max)
	for sy := min(ya, ye); sy <= max(ya, ye); sy++ {
		g.tb.Put(x, sy, '|')
		if mirror >= 1 {
			g.tb.Put(x1, sy, '|')
		}
	}
	if yrange.ShowZero && yrange.Min < 0 && yrange.Max > 0 {
		z := yrange.Data2Screen(0)
		for xx := x + 1; xx < x1; xx += 2 {
			g.tb.Put(xx, z, '-')
		}
	}

	if label != "" {
		g.tb.Text(1, (ya+ye)/2, label, 3)
	}

	for _, tic := range yrange.Tics {
		y := yrange.Data2Screen(tic.Pos)
		ly := yrange.Data2Screen(tic.LabelPos)
		if yrange.Time {
			g.tb.Put(x, y, '+')
			if mirror >= 2 {
				g.tb.Put(x1, y, '+')
			}
			if tic.Align == 0 { // centered tic
				g.tb.Put(x-1, y, '-')
				g.tb.Put(x-2, y, '-')
			}
			g.tb.Text(x, ly, tic.Label+" ", 1)
		} else {
			g.tb.Put(x, y, '+')
			if mirror >= 2 {
				g.tb.Put(x1, y, '+')
			}
			g.tb.Text(x-2, ly, tic.Label, 1)
		}
	}
}

func (g *TextGraphics) Scatter(points []chart.EPoint, plotstyle chart.PlotStyle, style chart.Style) {
	// First pass: Error bars
	for _, p := range points {
		xl, yl, xh, yh := p.BoundingBox()
		if !math.IsNaN(p.DeltaX) {
			g.tb.Line(int(xl), int(p.Y), int(xh), int(p.Y), '-')
		}
		if !math.IsNaN(p.DeltaY) {
			g.tb.Line(int(p.X), int(yl), int(p.X), int(yh), '|')
		}
	}

	// Second pass: Line
	if (plotstyle&chart.PlotStyleLines) != 0 && len(points) > 0 {
		lastx, lasty := int(points[0].X), int(points[0].Y)
		for i := 1; i < len(points); i++ {
			x, y := int(points[i].X), int(points[i].Y)
			// fmt.Printf("LineSegment %d (%d,%d) -> (%d,%d)\n", i, lastx,lasty,x,y)
			g.tb.Line(lastx, lasty, x, y, style.Symbol)
			lastx, lasty = x, y
		}
	}

	// Third pass: symbols
	if (plotstyle&chart.PlotStylePoints) != 0 && len(points) != 0 {
		for _, p := range points {
			g.tb.Put(int(p.X), int(p.Y), style.Symbol)
		}
	}
	// chart.GenericScatter(g, points, plotstyle, style)
}

func (g *TextGraphics) Boxes(boxes []chart.Box, width int, style chart.Style) {
	if width%2 == 0 {
		width += 1
	}
	hbw := (width - 1) / 2
	if style.Symbol == 0 {
		style.Symbol = '*'
	}

	for _, box := range boxes {
		x := int(box.X)
		q1, q3 := int(box.Q1), int(box.Q3)
		g.tb.Rect(x-hbw, q1, 2*hbw, q3-q1, 0, ' ')
		if !math.IsNaN(box.Med) {
			med := int(box.Med)
			g.tb.Put(x-hbw, med, '+')
			for i := 0; i < hbw; i++ {
				g.tb.Put(x-i, med, '-')
				g.tb.Put(x+i, med, '-')
			}
			g.tb.Put(x+hbw, med, '+')
		}

		if !math.IsNaN(box.Avg) && style.Symbol != 0 {
			g.tb.Put(x, int(box.Avg), style.Symbol)
		}

		if !math.IsNaN(box.High) {
			for y := int(box.High); y < q3; y++ {
				g.tb.Put(x, y, '|')
			}
		}

		if !math.IsNaN(box.Low) {
			for y := int(box.Low); y > q1; y-- {
				g.tb.Put(x, y, '|')
			}
		}

		for _, ol := range box.Outliers {
			y := int(ol)
			g.tb.Put(x, y, style.Symbol)
		}
	}
}


func (g *TextGraphics) Key(x, y int, key chart.Key) {
	m := key.Place()
	if len(m) == 0 {
		return
	}
	tw, th, cw, rh := key.Layout(g, m)
	style := g.Style("key")
	if style.LineWidth > 0 || style.FillColor != "" {
		g.tb.Rect(x, y, tw, th, 1, ' ')
	}
	x += int(chart.KeyHorSep)
	vsep := chart.KeyVertSep
	if vsep < 1 {
		vsep = 1
	}
	y += int(vsep)
	for ci, col := range m {
		yy := y

		for ri, e := range col {
			if e == nil || e.Text == "" {
				continue
			}
			plotStyle := e.PlotStyle
			// fmt.Printf("KeyEntry %s: PlotStyle = %d\n", e.Text, e.PlotStyle)
			if plotStyle == -1 {
				// heading only...
				g.tb.Text(x, yy, e.Text, -1)
			} else {
				// normal entry
				if (plotStyle & chart.PlotStyleLines) != 0 {
					g.Line(x, yy, x+int(chart.KeySymbolWidth), yy, e.Style)
				}
				if (plotStyle & chart.PlotStylePoints) != 0 {
					g.Symbol(x+int(chart.KeySymbolWidth/2), yy, e.Style)
				}
				if (plotStyle & chart.PlotStyleBox) != 0 {
					g.tb.Put(x+int(chart.KeySymbolWidth/2), yy, e.Style.Symbol)
				}
				g.tb.Text(x+int((chart.KeySymbolWidth+chart.KeySymbolSep)), yy, e.Text, -1)
			}
			yy += rh[ri] + int(chart.KeyRowSep)
		}

		x += int((chart.KeySymbolWidth + chart.KeySymbolSep + chart.KeyColSep + float32(cw[ci])))
	}

}

func (g *TextGraphics) Bars(bars []chart.Barinfo, style chart.Style) {
	chart.GenericBars(g, bars, style)
}

var CircleStretchFactor float64 = 1.85

func (g *TextGraphics) Rings(wedges []chart.Wedgeinfo, x, y, ro, ri int) {
	if g.xoff == -1 {
		g.xoff = int(float64(ro) * (CircleStretchFactor - 1))
		debug.Printf("Shifting center about %d (ro=%d, f=%.2f)", g.xoff, ro, CircleStretchFactor)
	}
	for i := range wedges {
		wedges[i].Style.LineWidth = 1
	}
	chart.GenericRings(g, wedges, x+g.xoff, y, ro, ri, 1.8)
}
