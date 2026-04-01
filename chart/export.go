package chart

import (
	"errors"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-glyph"
	"github.com/mike-ward/go-gui/gui"
)

// themer is satisfied by chart views that expose their theme.
type themer interface {
	chartTheme() *theme.Theme
}

// Drawer is implemented by chart views that support
// headless export via ExportPNG.
type Drawer interface {
	Draw(dc *gui.DrawContext)
}

// ExportPNG renders a chart view to a PNG file at the given
// pixel dimensions. v must be a chart view created by Line,
// Bar, Area, Scatter, or Pie.
func ExportPNG(v gui.View, width, height int, path string) error {
	d, ok := v.(Drawer)
	if !ok {
		return errors.New("chart: view does not implement Drawer")
	}
	if width <= 0 || height <= 0 {
		return errors.New("chart: dimensions must be positive")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background from theme, fall back to white.
	bg := gui.Color{R: 255, G: 255, B: 255, A: 255}
	if t, ok := v.(themer); ok && t.chartTheme() != nil {
		c := t.chartTheme().Background
		if c.IsSet() {
			bg = c
		}
	}
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i] = bg.R
		img.Pix[i+1] = bg.G
		img.Pix[i+2] = bg.B
		img.Pix[i+3] = bg.A
	}

	backend := &imageBackend{
		target:   img,
		textures: make(map[glyph.TextureID]atlasPage),
	}

	textSys, err := glyph.NewTextSystem(backend)
	if err != nil {
		return err
	}
	defer textSys.Free()

	measurer := &pngTextMeasurer{textSys: textSys}
	dc := gui.NewDrawContext(float32(width), float32(height), measurer)
	d.Draw(dc)

	rasterizeTriBatches(img, dc.Batches())

	for _, t := range dc.Texts() {
		cfg := toGlyphConfig(t.Style)
		_ = textSys.DrawText(t.X, t.Y, t.Text, cfg)
	}
	textSys.Commit()
	backend.flush()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encErr := png.Encode(f, img)
	closeErr := f.Close()
	if encErr != nil {
		return encErr
	}
	return closeErr
}

// -----------------------------------------------------------
// Triangle rasterizer
// -----------------------------------------------------------

// 16x MSAA sub-pixel offsets (4×4 grid).
var msaaOffsets = [16][2]float32{
	{0.125, 0.125}, {0.375, 0.125}, {0.625, 0.125}, {0.875, 0.125},
	{0.125, 0.375}, {0.375, 0.375}, {0.625, 0.375}, {0.875, 0.375},
	{0.125, 0.625}, {0.375, 0.625}, {0.625, 0.625}, {0.875, 0.625},
	{0.125, 0.875}, {0.375, 0.875}, {0.625, 0.875}, {0.875, 0.875},
}

const msaaSamples = uint32(len(msaaOffsets))

// edgeInside returns true when all three edge values have the
// same sign (>= 0 or <= 0), indicating the sample is inside.
func edgeInside(e0, e1, e2 float32) bool {
	return (e0 >= 0 && e1 >= 0 && e2 >= 0) ||
		(e0 <= 0 && e1 <= 0 && e2 <= 0)
}

// rasterizeTriBatches fills triangles onto img using
// edge-function rasterization with 16x MSAA and src-over
// compositing.
func rasterizeTriBatches(img *image.RGBA, batches []gui.DrawCanvasTriBatch) {
	bounds := img.Bounds()
	stride := img.Stride
	pix := img.Pix

	for _, batch := range batches {
		cr := uint32(batch.Color.R)
		cg := uint32(batch.Color.G)
		cb := uint32(batch.Color.B)
		ca := uint32(batch.Color.A)
		tris := batch.Triangles

		for i := 0; i+5 < len(tris); i += 6 {
			x0, y0 := tris[i], tris[i+1]
			x1, y1 := tris[i+2], tris[i+3]
			x2, y2 := tris[i+4], tris[i+5]

			mnX := int(math.Floor(float64(min(x0, x1, x2))))
			mnY := int(math.Floor(float64(min(y0, y1, y2))))
			mxX := int(math.Ceil(float64(max(x0, x1, x2))))
			mxY := int(math.Ceil(float64(max(y0, y1, y2))))

			mnX = max(mnX, bounds.Min.X)
			mnY = max(mnY, bounds.Min.Y)
			mxX = min(mxX, bounds.Max.X-1)
			mxY = min(mxY, bounds.Max.Y-1)

			dx01, dy01 := x1-x0, y1-y0
			dx12, dy12 := x2-x1, y2-y1
			dx20, dy20 := x0-x2, y0-y2

			for py := mnY; py <= mxY; py++ {
				rowOff := py*stride + mnX*4
				for px := mnX; px <= mxX; px++ {
					// Count sub-pixel samples inside the
					// triangle (4x MSAA).
					hits := uint32(0)
					for _, off := range msaaOffsets {
						sx := float32(px) + off[0]
						sy := float32(py) + off[1]
						e0 := dx01*(sy-y0) - dy01*(sx-x0)
						e1 := dx12*(sy-y1) - dy12*(sx-x1)
						e2 := dx20*(sy-y2) - dy20*(sx-x2)
						if edgeInside(e0, e1, e2) {
							hits++
						}
					}
					if hits > 0 {
						sa := ca * hits / msaaSamples
						blendPixel(pix, rowOff, cr, cg, cb, sa)
					}
					rowOff += 4
				}
			}
		}
	}
}

// blendPixel applies src-over compositing at pix offset off.
func blendPixel(pix []byte, off int, sr, sg, sb, sa uint32) {
	if sa == 0 {
		return
	}
	if sa == 255 {
		pix[off] = uint8(sr)
		pix[off+1] = uint8(sg)
		pix[off+2] = uint8(sb)
		pix[off+3] = 255
		return
	}
	inv := 255 - sa
	pix[off] = uint8((sr*sa + uint32(pix[off])*inv) / 255)
	pix[off+1] = uint8((sg*sa + uint32(pix[off+1])*inv) / 255)
	pix[off+2] = uint8((sb*sa + uint32(pix[off+2])*inv) / 255)
	pix[off+3] = uint8(sa + uint32(pix[off+3])*inv/255)
}

// -----------------------------------------------------------
// Image-backed glyph DrawBackend
// -----------------------------------------------------------

type atlasPage struct {
	data          []byte
	width, height int
}

type pendingDraw struct {
	id       glyph.TextureID
	src, dst glyph.Rect
	color    glyph.Color
}

// imageBackend implements glyph.DrawBackend, compositing glyph
// quads from atlas pages onto a target *image.RGBA.
type imageBackend struct {
	target   *image.RGBA
	textures map[glyph.TextureID]atlasPage
	nextID   glyph.TextureID
	pending  []pendingDraw
}

func (b *imageBackend) DPIScale() float32 { return 1.0 }

func (b *imageBackend) NewTexture(width, height int) glyph.TextureID {
	id := b.nextID
	b.nextID++
	b.textures[id] = atlasPage{
		data:   make([]byte, width*height*4),
		width:  width,
		height: height,
	}
	return id
}

func (b *imageBackend) UpdateTexture(id glyph.TextureID, data []byte) {
	if page, ok := b.textures[id]; ok {
		copy(page.data, data)
	}
}

func (b *imageBackend) DeleteTexture(id glyph.TextureID) {
	delete(b.textures, id)
}

func (b *imageBackend) DrawTexturedQuad(
	id glyph.TextureID, src, dst glyph.Rect, c glyph.Color,
) {
	b.pending = append(b.pending, pendingDraw{id, src, dst, c})
}

func (b *imageBackend) DrawFilledRect(dst glyph.Rect, c glyph.Color) {
	if c.A == 0 {
		return
	}
	img := b.target
	x0 := max(int(dst.X), 0)
	y0 := max(int(dst.Y), 0)
	x1 := min(int(dst.X+dst.Width), img.Bounds().Max.X)
	y1 := min(int(dst.Y+dst.Height), img.Bounds().Max.Y)
	cr, cg, cb, ca := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	for py := y0; py < y1; py++ {
		off := py*img.Stride + x0*4
		for px := x0; px < x1; px++ {
			blendPixel(img.Pix, off, cr, cg, cb, ca)
			off += 4
		}
	}
}

func (b *imageBackend) DrawTexturedQuadTransformed(
	id glyph.TextureID, src, dst glyph.Rect,
	c glyph.Color, _ glyph.AffineTransform,
) {
	// Chart text uses identity transforms; ignore the
	// transform and fall back to the untransformed path.
	b.DrawTexturedQuad(id, src, dst, c)
}

// flush replays recorded draw commands after atlas textures
// have been committed via UpdateTexture.
func (b *imageBackend) flush() {
	for _, d := range b.pending {
		b.blitQuad(d.id, d.src, d.dst, d.color)
	}
	b.pending = b.pending[:0]
}

func (b *imageBackend) blitQuad(
	id glyph.TextureID, src, dst glyph.Rect, c glyph.Color,
) {
	page, ok := b.textures[id]
	if !ok || dst.Width <= 0 || dst.Height <= 0 {
		return
	}

	img := b.target
	imgW := img.Bounds().Max.X
	imgH := img.Bounds().Max.Y

	dstX0 := max(int(dst.X), 0)
	dstY0 := max(int(dst.Y), 0)
	dstX1 := min(int(math.Ceil(float64(dst.X+dst.Width))), imgW)
	dstY1 := min(int(math.Ceil(float64(dst.Y+dst.Height))), imgH)

	// Scale factors: src pixels per dst pixel.
	scaleX := src.Width / dst.Width
	scaleY := src.Height / dst.Height

	cr, cg, cb, ca := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)

	for py := dstY0; py < dstY1; py++ {
		// Nearest-neighbor sample in atlas.
		sy := int(src.Y + float32(py-int(dst.Y))*scaleY)
		if sy < 0 || sy >= page.height {
			continue
		}
		srcRowOff := sy * page.width * 4
		imgOff := py*img.Stride + dstX0*4

		for px := dstX0; px < dstX1; px++ {
			sx := int(src.X + float32(px-int(dst.X))*scaleX)
			if sx < 0 || sx >= page.width {
				imgOff += 4
				continue
			}
			sOff := srcRowOff + sx*4
			texA := uint32(page.data[sOff+3])
			if texA == 0 {
				imgOff += 4
				continue
			}

			// Glyph alpha modulated by tint alpha.
			ga := texA * ca / 255
			if ga == 0 {
				imgOff += 4
				continue
			}

			// Src-over with tint color.
			inv := 255 - ga
			img.Pix[imgOff] = uint8((cr*ga + uint32(img.Pix[imgOff])*inv) / 255)
			img.Pix[imgOff+1] = uint8((cg*ga + uint32(img.Pix[imgOff+1])*inv) / 255)
			img.Pix[imgOff+2] = uint8((cb*ga + uint32(img.Pix[imgOff+2])*inv) / 255)
			img.Pix[imgOff+3] = uint8(ga + uint32(img.Pix[imgOff+3])*inv/255)
			imgOff += 4
		}
	}
}

// -----------------------------------------------------------
// TextMeasurer wrapping glyph.TextSystem
// -----------------------------------------------------------

type pngTextMeasurer struct {
	textSys *glyph.TextSystem
}

func (m *pngTextMeasurer) TextWidth(text string, style gui.TextStyle) float32 {
	w, err := m.textSys.TextWidth(text, toGlyphConfig(style))
	if err != nil {
		return 0
	}
	return w
}

func (m *pngTextMeasurer) TextHeight(text string, style gui.TextStyle) float32 {
	h, err := m.textSys.TextHeight(text, toGlyphConfig(style))
	if err != nil {
		return 0
	}
	return h
}

func (m *pngTextMeasurer) FontHeight(style gui.TextStyle) float32 {
	h, err := m.textSys.FontHeight(toGlyphConfig(style))
	if err != nil {
		return style.Size * 1.4
	}
	return h
}

func (m *pngTextMeasurer) FontAscent(style gui.TextStyle) float32 {
	met, err := m.textSys.FontMetrics(toGlyphConfig(style))
	if err != nil {
		return style.Size * 0.8
	}
	return met.Ascender
}

func (m *pngTextMeasurer) LayoutText(
	text string, style gui.TextStyle, wrapWidth float32,
) (glyph.Layout, error) {
	cfg := toGlyphConfig(style)
	if wrapWidth > 0 {
		cfg.Block.Width = wrapWidth
		cfg.Block.Wrap = glyph.WrapWord
	} else if wrapWidth < 0 {
		cfg.Block.Width = -wrapWidth
		cfg.Block.Wrap = glyph.WrapNone
	}
	return m.textSys.LayoutText(text, cfg)
}

// -----------------------------------------------------------
// Style conversion (inlined from glyphconv)
// -----------------------------------------------------------

func toGlyphConfig(s gui.TextStyle) glyph.TextConfig {
	align := glyph.AlignLeft
	switch s.Align {
	case gui.TextAlignCenter:
		align = glyph.AlignCenter
	case gui.TextAlignRight:
		align = glyph.AlignRight
	}
	return glyph.TextConfig{
		Style: glyph.TextStyle{
			FontName:      s.Family,
			Size:          s.Size,
			Color:         glyph.Color{R: s.Color.R, G: s.Color.G, B: s.Color.B, A: s.Color.A},
			BgColor:       glyph.Color{R: s.BgColor.R, G: s.BgColor.G, B: s.BgColor.B, A: s.BgColor.A},
			Typeface:      s.Typeface,
			Underline:     s.Underline,
			Strikethrough: s.Strikethrough,
			LetterSpacing: s.LetterSpacing,
			StrokeWidth:   s.StrokeWidth,
			StrokeColor:   glyph.Color{R: s.StrokeColor.R, G: s.StrokeColor.G, B: s.StrokeColor.B, A: s.StrokeColor.A},
			Features:      s.Features,
		},
		Block: glyph.BlockStyle{
			Align:       align,
			Wrap:        glyph.WrapWord,
			Width:       -1,
			LineSpacing: s.LineSpacing,
		},
		Gradient: s.Gradient,
	}
}
