package imageutil

import (
	"image"
	"image/draw"
	"image/jpeg"
	"io"
)

func EncodeToJPEG(w io.Writer, m image.Image, quality int) error {
	if m == nil {
		return jpeg.Encode(w, m, &jpeg.Options{Quality: quality})
	}

	// 核心拦截：检查传入的 Image 对象底层的色彩模型或结构。
	// 诸如 WebP 或某些特殊格式解码后，虽然是 image.Image，但其底层的 Opaque/Stride 等物理布局
	// 会导致标准的 jpeg.Encode 编码出来的字节流带有特殊的标记，让 Emby 的部分客户端（如 Android/TV）无法解析。
	// 我们通过将其强行绘制到一个标准的、完全干净的 RGBA 像素矩阵中，彻底抹除原始格式痕迹。
	bounds := m.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	
	// 使用 Src 模式将像素纯净地复制到新画布上
	draw.Draw(rgbaImg, bounds, m, bounds.Min, draw.Src)

	// 使用清洗后的通用 RGBA 画布进行 JPEG 编码输出
	return jpeg.Encode(w, rgbaImg, &jpeg.Options{Quality: quality})
}
