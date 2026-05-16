package pdf

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/tdarci/go-nothings/quakes-roster/internal/model"
	"github.com/tdarci/go-nothings/quakes-roster/internal/util"

	"github.com/jung-kurt/gofpdf"
)

const (
	pageW = 215.9
	pageH = 279.4
)

func Run() {
	var players []model.Player
	util.LoadJSON("players.json", &players)

	// // Sort by jersey number (important for roster readability)
	// sort.Slice(players, func(i, j int) bool {
	// 	ni, _ := strconv.Atoi(players[i].Number)
	// 	nj, _ := strconv.Atoi(players[j].Number)
	// 	return ni < nj
	// })

	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetFont("Helvetica", "", 10)
	pdf.AliasNbPages("")
	pdf.SetCompression(true)

	cellW := (pageW - 20) / 3
	cellH := (pageH - 20) / 3

	for i, p := range players {

		// new page every 9 players
		if i%9 == 0 {
			pdf.AddPage()
		}

		col := i % 3
		row := (i / 3) % 3

		x := 10 + float64(col)*cellW
		y := 10 + float64(row)*cellH

		drawPlayer(pdf, p, x, y, cellW, cellH)
	}

	err := pdf.OutputFileAndClose("roster.pdf")
	if err != nil {
		log.Fatal(err)
	}
}

func drawPlayer(pdf *gofpdf.Fpdf, p model.Player, x, y, w, h float64) {

	// -----------------------------
	// IMAGE (safe aspect ratio + no overlap)
	// -----------------------------
	imgPath := cropAndSave(p.ImagePath)

	f, err := os.Open(imgPath)
	if err != nil {
		return
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return
	}

	imgW := float64(cfg.Width)
	imgH := float64(cfg.Height)

	targetW := w
	scale := targetW / imgW
	drawH := imgH * scale

	maxImgH := h * 0.60
	if drawH > maxImgH {
		drawH = maxImgH
		targetW = (drawH / imgH) * imgW
	}

	offsetX := x + (w-targetW)/2

	pdf.Image(imgPath, offsetX, y, targetW, drawH, false, "", 0, "")

	// -----------------------------
	// JERSEY NUMBER BADGE (overlay)
	// -----------------------------
	badgeR := 4.5

	cx := offsetX + badgeR + 1.5
	cy := y + badgeR + 1.5

	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Circle(cx, cy, badgeR, "FD")

	// BOLD NUMBER INSIDE BADGE
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(0, 0, 0)

	// slight centering tweak for bold text
	pdf.SetXY(cx-2.5, cy-2.5)
	pdf.CellFormat(5, 5, p.Number, "", 0, "C", false, 0, "")

	// -----------------------------
	// TEXT (below image, no jersey number here anymore)
	// -----------------------------
	textY := y + drawH + 2

	// NAME
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(x, textY)
	pdf.CellFormat(w, 5, p.Name, "", 0, "C", false, 0, "")

	// POSITION ONLY (no number anymore)
	pdf.SetFont("Helvetica", "", 8)
	line1 := p.Position

	pdf.SetXY(x, textY+5)
	pdf.CellFormat(w, 4, line1, "", 0, "C", false, 0, "")

	// AGE + BIRTHPLACE
	line2 := ""

	if p.Age != "" {
		line2 += "Age " + p.Age
	}

	if p.Birthplace != "" {
		if line2 != "" {
			line2 += " - "
		}
		line2 += p.Birthplace
	}

	pdf.SetXY(x, textY+9)
	pdf.CellFormat(w, 4, line2, "", 0, "C", false, 0, "")
}

func cropAndSave(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return path
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return path
	}

	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	// center vertical crop (keeps faces)
	left := w / 4
	right := left + w/2

	cropped := image.NewRGBA(image.Rect(0, 0, right-left, h))

	for y := 0; y < h; y++ {
		for x := left; x < right; x++ {
			cropped.Set(x-left, y, img.At(x, y))
		}
	}

	outPath := path + "_crop.jpg"

	out, err := os.Create(outPath)
	if err != nil {
		return path
	}
	defer out.Close()

	// IMPORTANT: ensure consistent encoding
	jpeg.Encode(out, cropped, &jpeg.Options{Quality: 92})

	return outPath
}
