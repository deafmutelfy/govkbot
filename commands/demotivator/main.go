package demotivator

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	WIDTH = 700

	RATIO_LINE = 0.00655737704918032786
	RATIO_SIDE = 0.11147540983606557377

	RATIO_H_SIDE         = 0.09382716049382716049
	RATIO_H_BEFORE_TEXT  = 0.05432098765432098765
	RATIO_H_BETWEEN_TEXT = 0.03
	RATIO_H_AFTER_TEXT   = 0.07407407407407407407

	RATIO_FONT_TITLE = 0.052
	RATIO_FONT_BODY  = 0.025

	FONT_TITLE = ""
	FONT       = "sans"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"демотиватор"},
		Description: "сгенерировать демотиватор",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	atts := core.ExtractAttachments(obj, "photo")
	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

		return
	}
	args := core.ExtractArguments(obj)
	if len(args) == 0 {
		core.ReplySimple(obj, "ошибка: необходимо указать текст")

		return
	}

	toks := strings.Split(strings.Join(args, " "), "\n")
	title := toks[0]
	body := ""
	if len(toks) > 1 {
		body = strings.Join(toks[1:], "\n")
	}

	attachment := atts[0]

	response, err := http.Get(attachment.Photo.MaxSize().URL)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	bt, err := io.ReadAll(response.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	src := imagick.NewMagickWand()
	src.ReadImageBlob(bt)
	src.ScaleImage(WIDTH, uint(WIDTH/float64(src.GetImageWidth())*float64(src.GetImageHeight())))

	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	cRatioLine := float64(src.GetImageWidth()) * RATIO_LINE
	withBorderWidth := cRatioLine*4 + float64(src.GetImageWidth())
	withBorderHeight := cRatioLine*4 + float64(src.GetImageHeight())

	pw.SetColor("white")
	dw.SetStrokeWidth(cRatioLine)
	dw.SetStrokeColor(pw)
	dw.Rectangle(0, 0, withBorderWidth, withBorderHeight)

	parts := imagick.NewMagickWand()

	parts.NewImage(uint(withBorderWidth), uint(withBorderHeight), pw)
	parts.DrawImage(dw)

	cRatioSide := float64(WIDTH) * RATIO_SIDE
	cWidth := withBorderWidth + cRatioSide*2

	pw.SetColor("black")
	parts.SetFont(FONT_TITLE)
	parts.SetPointsize(cWidth * RATIO_FONT_TITLE)
	parts.SetGravity(imagick.GRAVITY_CENTER)
	parts.SetBackgroundColor(pw)
	parts.SetSize(uint(float64(WIDTH)-cRatioSide*2), 0)
	parts.SetOption("fill", "white")
	parts.SetOption("pango:wrap", "word-char")
	parts.ReadImage("pango:" + title)

	parts.SetFont(FONT)
	parts.SetPointsize(cWidth * RATIO_FONT_BODY)
	parts.ReadImage("pango:" + body)

	dst := imagick.NewMagickWand()

	height := src.GetImageHeight()
	cRatioHSide := float64(height) * RATIO_H_SIDE
	cRatioBeforeText := float64(height) * RATIO_H_BEFORE_TEXT
	cRatioAfterText := float64(height) * RATIO_H_AFTER_TEXT
	cRatioBetweenText := float64(height) * RATIO_H_BETWEEN_TEXT
	parts.SetIteratorIndex(1)
	wt, ht, _, _, _ := parts.GetImagePage()
	parts.SetIteratorIndex(2)
	wb, hb, _, _, _ := parts.GetImagePage()

	dst.NewImage(
		uint(cWidth),
		uint(cRatioHSide+withBorderHeight+cRatioBeforeText+float64(ht)+cRatioBetweenText+float64(hb)+cRatioAfterText),
		pw,
	)

	parts.SetIteratorIndex(0)
	dst.CompositeImage(parts, imagick.COMPOSITE_OP_SRC_OVER, int(cRatioSide), int(cRatioHSide))
	dst.CompositeImage(src, imagick.COMPOSITE_OP_SRC_OVER, int(cRatioSide+cRatioLine*2), int(cRatioHSide+cRatioLine*2))
	parts.SetIteratorIndex(1)
	dst.CompositeImage(parts, imagick.COMPOSITE_OP_SRC_OVER, int((cWidth-float64(wt))/2), int(cRatioHSide+withBorderHeight+cRatioBeforeText))
	parts.SetIteratorIndex(2)
	dst.CompositeImage(parts, imagick.COMPOSITE_OP_SRC_OVER, int((cWidth-float64(wb))/2), int(cRatioHSide+withBorderHeight+cRatioBeforeText+float64(ht)+cRatioBetweenText))

	dst.SetImageFormat("png")
	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(0, bytes.NewReader(dst.GetImageBlob()))

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)

	return
}
