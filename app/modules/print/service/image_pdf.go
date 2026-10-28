package service

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image/color"
	"image/jpeg"
	"image/png"
)

// A4 纸在 PDF 中的点尺寸（1 点 = 1/72 英寸），图片按比例居中铺满。
const (
	a4Width  = 595.28
	a4Height = 841.89
)

// imageToPDF 把单张图片包成单页 A4 PDF。
// JPEG 直接以 DCTDecode 内嵌（不重新编码，避免二次损失）；PNG 解码后按 FlateDecode 写入。
// 参数 content 为图片字节，ext 为扩展名（jpg/jpeg/png）。
func imageToPDF(content []byte, ext string) ([]byte, error) {
	var (
		width, height int
		colorSpace    string
		filter        string
		stream        []byte
	)

	switch ext {
	case "jpg", "jpeg":
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("JPEG 解析失败：%w", err)
		}
		width, height = cfg.Width, cfg.Height
		colorSpace = "/DeviceRGB"
		if cfg.ColorModel == color.GrayModel {
			colorSpace = "/DeviceGray"
		}
		filter = "/DCTDecode"
		stream = content
	case "png":
		img, err := png.Decode(bytes.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("PNG 解析失败：%w", err)
		}
		bounds := img.Bounds()
		width, height = bounds.Dx(), bounds.Dy()

		raw := make([]byte, 0, width*height*3)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				// 透明像素按白底合成，否则 PDF 里会变成黑块
				if a < 0xffff {
					raw = append(raw,
						blendWhite(r, a),
						blendWhite(g, a),
						blendWhite(b, a),
					)
					continue
				}
				raw = append(raw, byte(r>>8), byte(g>>8), byte(b>>8))
			}
		}

		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		if _, err := zw.Write(raw); err != nil {
			return nil, err
		}
		if err := zw.Close(); err != nil {
			return nil, err
		}

		stream = buf.Bytes()
		filter = "/FlateDecode"
		colorSpace = "/DeviceRGB"
	default:
		return nil, fmt.Errorf("不支持的图片类型：%s", ext)
	}

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("图片尺寸异常")
	}

	return buildSinglePagePDF(width, height, colorSpace, filter, stream), nil
}

// blendWhite 把带 alpha 的分量合成到白底。
func blendWhite(value, alpha uint32) byte {
	// value 与 alpha 均为 16 位，先归一化再合成
	blended := (value*alpha + 0xffff*(0xffff-alpha)) / 0xffff

	return byte(blended >> 8)
}

// buildSinglePagePDF 组装单页 PDF：图片按比例缩放并居中到 A4 页面。
func buildSinglePagePDF(width, height int, colorSpace, filter string, stream []byte) []byte {
	scale := a4Width / float64(width)
	if s := a4Height / float64(height); s < scale {
		scale = s
	}

	drawWidth := float64(width) * scale
	drawHeight := float64(height) * scale
	offsetX := (a4Width - drawWidth) / 2
	offsetY := (a4Height - drawHeight) / 2

	content := fmt.Sprintf("q %.2f 0 0 %.2f %.2f %.2f cm /Im0 Do Q", drawWidth, drawHeight, offsetX, offsetY)

	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Resources << /XObject << /Im0 4 0 R >> >> /Contents 5 0 R >>", a4Width, a4Height),
		fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace %s /BitsPerComponent 8 /Filter %s /Length %d >>\nstream\n%s\nendstream",
			width, height, colorSpace, filter, len(stream), stream),
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content),
	}

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")

	offsets := make([]int, len(objects)+1)
	for i, body := range objects {
		offsets[i+1] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", i+1, body)
	}

	xrefOffset := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", len(objects)+1)
	buf.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)

	return buf.Bytes()
}
