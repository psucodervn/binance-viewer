package telegram

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"sort"
	"time"

	"github.com/fogleman/gg"
	"github.com/nleeper/goment"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/colornames"
	"gopkg.in/tucnak/telebot.v2"

	"copytrader/internal/model"
	"copytrader/internal/user"
)

var (
	FontFaceLarge, _  = gg.LoadFontFace("assets/fonts/Plex.ttf", 30)
	FontFace, _       = gg.LoadFontFace("assets/fonts/Plex.ttf", 20)
	FontFaceItalic, _ = gg.LoadFontFace("assets/fonts/PlexItalic.ttf", 16)
	ColorProfit       = color.RGBA{R: 14, G: 203, B: 139, A: 255}
	ColorLoss         = color.RGBA{R: 246, G: 70, B: 93, A: 255}
	ColorBinance      = color.RGBA{R: 240, G: 185, B: 11, A: 255}
	Background, _     = gg.LoadJPG("assets/bg.jpg")
	BackgroundPattern = gg.NewSurfacePattern(Background, gg.RepeatY)
	LocalLocation, _  = time.LoadLocation("Asia/Ho_Chi_Minh")
)

func profitColor(v float64) color.Color {
	if v < 0 {
		return ColorLoss
	}
	return ColorProfit
}

func cmdImageCard(b *Bot) interface{} {
	return func(m *telebot.Message) {
		start := time.Now()
		u := b.loadUser(m.Chat)
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()

		from, _ := goment.New()
		from.UTC().StartOf("day")
		to, _ := goment.New(from)
		to = to.Add(1, "day").Subtract(1, "second")

		accounts, err := prepareAccounts(ctx, u, from, to)
		if err != nil {
			_, _ = b.bot.Send(m.Chat, "Query Binance API failed")
			return
		}

		// args := splitArgs(m.Payload)
		// if len(args) > 0 {
		//   d, _ := strconv.Atoi(args[0])
		//   if d > 0 {
		//     accounts = accounts[:d]
		//   }
		// }

		titleY := 80.0
		tableY := 150.0
		if len(accounts) > 2 {
			titleY -= 10.0
			tableY -= 20.0
		}
		tableH := 40.0 * float64(len(accounts)+2)
		width, height := Background.Bounds().Dx(), int(math.Max(360, tableY+tableH+40))

		dc := gg.NewContext(width, height)
		dc.SetFillStyle(BackgroundPattern)
		dc.DrawRectangle(0, 0, float64(width), float64(height))
		dc.Fill()

		dc.SetFontFace(FontFaceLarge)
		dc.SetColor(ColorBinance)
		dc.DrawString(u.Name, 40, titleY)
		tw, _ := dc.MeasureString(u.Name)
		dc.SetFontFace(FontFace)
		dc.SetColor(colornames.Gray)
		dc.DrawString("'s accounts", 40+tw+5, titleY)

		dc.SetFontFace(FontFace)
		tbl := buildInfoTable(accounts)
		drawTable(dc, 40, tableY, float64(dc.Width()-2*40), tableH, tbl)

		dc.SetFontFace(FontFaceItalic)
		dc.SetColor(colornames.Gray)
		dc.DrawString(fmt.Sprintf("(*) updated at %s", time.Now().In(LocalLocation).Format("02/01/2006, 15:04:05 GMT MST")), 40, float64(dc.Height())-40)

		var bf bytes.Buffer
		if err := png.Encode(&bf, dc.Image()); err != nil {
			log.Err(err).Msg("encode image failed")
			return
		}
		log.Debug().Dur("elapsed", time.Now().Sub(start)).Send()

		if _, err := b.bot.Send(m.Chat, &telebot.Photo{File: telebot.FromReader(&bf)}); err != nil {
			log.Err(err).Msg("send image failed")
		}
	}
}

func prepareAccounts(ctx context.Context, u model.User, from *goment.Goment, to *goment.Goment) ([]accountInfo, error) {
	var accounts []accountInfo
	for _, a := range u.Accounts {
		cli := user.GetUserClient(a)
		acc, err := cli.Info(ctx)
		if err != nil {
			return nil, err
		}
		incomes, err := cli.ListIncome(ctx, from.ToTime(), to.ToTime())
		accounts = append(accounts, accountInfo{a, acc, incomes})
		// accounts = append(accounts, accountInfo{Info: a, Account: &futures.Account{}})
	}
	// accounts = append(accounts, accounts[0], accounts[1], accounts[0], accounts[1])
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Info.Name < accounts[j].Info.Name
	})
	return accounts, nil
}

func drawTable(dc *gg.Context, x, y float64, width, height float64, tbl Table) {
	xx := x
	yy := y
	lineHeight := height / (float64(len(tbl.Rows)) + 2.4)
	totalSpan := tbl.Headers.TotalSpan()
	spans := tbl.Headers.Spans()
	aligns := tbl.Headers.GetAligns()
	borderColor := colornames.Gray
	borderColor.A = 128

	dc.SetColor(colornames.Gray)
	for _, col := range tbl.Headers {
		if col.Align == AlignLeft {
			dc.DrawString(col.Text, xx, yy)
		} else if col.Align == AlignRight {
			dc.DrawStringAnchored(col.Text, xx+col.GetSpan()/totalSpan*width, yy, 1, 0)
		}
		xx += col.GetSpan() / totalSpan * width
	}

	dc.SetColor(borderColor)
	dc.DrawLine(x-2, yy+lineHeight*0.3, x+width+2, yy+lineHeight*0.3)
	dc.Stroke()

	// yy += lineHeight
	for _, row := range tbl.Rows {
		yy += lineHeight
		xx = x
		for i, col := range row.Columns {
			dc.SetColor(col.GetColor())
			text := fmt.Sprintf("%v", col.Data)
			if aligns[i] == AlignLeft {
				dc.DrawString(text, xx, yy)
			} else if aligns[i] == AlignRight {
				dc.DrawStringAnchored(text, xx+spans[i]/totalSpan*width, yy, 1, 0)
			}
			xx += spans[i] / totalSpan * width
		}
	}

	dc.SetColor(borderColor)
	dc.DrawLine(x-2, yy+lineHeight*0.4, x+width+2, yy+lineHeight*0.4)
	dc.Stroke()

	xx = x
	yy += lineHeight * 1.1
	for i, col := range tbl.Footers {
		if i == 0 {
			dc.SetColor(colornames.Gray)
		} else {
			dc.SetColor(col.GetColor())
		}
		text := fmt.Sprintf("%v", col.Data)
		if aligns[i] == AlignLeft {
			dc.DrawString(text, xx, yy)
		} else if aligns[i] == AlignRight {
			dc.DrawStringAnchored(text, xx+spans[i]/totalSpan*width, yy, 1, 0)
		}
		xx += spans[i] / totalSpan * width
	}
}
