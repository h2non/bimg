#include <stdlib.h>
#include <vips/vips.h>
#include <vips/vips7compat.h>

enum types {
	UNKNOWN = 0,
	JPEG,
	WEBP,
	PNG,
	TIFF,
	MAGICK
};

typedef struct {
	const char *Text;
	const char *Font;
} WatermarkTextOptions;

typedef struct {
	int    Width;
	int    DPI;
	int    Margin;
	int    NoReplicate;
	float  Opacity;
	double Background[3];
} WatermarkOptions;

void
vips_enable_cache_set_trace()
{
	vips_cache_set_trace(TRUE);
};

int
vips_affine_interpolator(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator)
{
	return vips_affine(in, out, a, b, c, d, "interpolate", interpolator, NULL);
};

int
vips_jpegload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink)
{
	return vips_jpegload_buffer(buf, len, out, "shrink", shrink, NULL);
};

int
vips_flip_bridge(VipsImage *in, VipsImage **out, int direction)
{
	return vips_flip(in, out, direction, NULL);
};

int
vips_shrink_bridge(VipsImage *in, VipsImage **out, double xshrink, double yshrink)
{
	return vips_shrink(in, out, xshrink, yshrink, NULL);
};

int
vips_rotate(VipsImage *in, VipsImage **out, int angle)
{
	int rotate = VIPS_ANGLE_D0;

	if (angle == 90) {
		rotate = VIPS_ANGLE_D90;
	} else if (angle == 180) {
		rotate = VIPS_ANGLE_D180;
	} else if (angle == 270) {
		rotate = VIPS_ANGLE_D270;
	}

	return vips_rot(in, out, rotate, NULL);
};

int
vips_exif_orientation(VipsImage *image) {
	int orientation = 0;
	const char *exif;
	if (
		vips_image_get_typeof(image, "exif-ifd0-Orientation") != 0 &&
		!vips_image_get_string(image, "exif-ifd0-Orientation", &exif)
	) {
		orientation = atoi(&exif[0]);
	}
	return orientation;
};

int
has_profile_embed(VipsImage *image) {
	return (vips_image_get_typeof(image, VIPS_META_ICC_NAME) > 0) ? 1 : 0;
};

int
has_alpha_channel(VipsImage *image) {
	return (
		(image->Bands == 2 && image->Type == VIPS_INTERPRETATION_B_W) ||
		(image->Bands == 4 && image->Type != VIPS_INTERPRETATION_CMYK) ||
		(image->Bands == 5 && image->Type == VIPS_INTERPRETATION_CMYK)
	) ? 1 : 0;
};

int
interpolator_window_size(char const *name) {
	VipsInterpolate *interpolator = vips_interpolate_new(name);
	int window_size = vips_interpolate_get_window_size(interpolator);
	g_object_unref(interpolator);
	return window_size;
};

const char *
vips_enum_nick_bridge(VipsImage *image) {
	return vips_enum_nick(VIPS_TYPE_INTERPRETATION, image->Type);
};

int
vips_zoom_bridge(VipsImage *in, VipsImage **out, int xfac, int yfac)
{
	return vips_zoom(in, out, xfac, yfac, NULL);
};

int
vips_embed_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend)
{
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
};

int
vips_extract_area_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height)
{
	return vips_extract_area(in, out, left, top, width, height, NULL);
};

int
vips_jpegsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace)
{
	return vips_jpegsave_buffer(in, buf, len, "strip", strip, "Q", quality, "optimize_coding", TRUE, "interlace", interlace, NULL);
};

int
vips_pngsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int compression, int quality, int interlace)
{
#if (VIPS_MAJOR_VERSION >= 8 || (VIPS_MAJOR_VERSION >= 7 && VIPS_MINOR_VERSION >= 42))
	return vips_pngsave_buffer(in, buf, len, "strip", FALSE, "compression", compression,
		"interlace", interlace, "filter", VIPS_FOREIGN_PNG_FILTER_NONE, NULL);
#else
	return vips_pngsave_buffer(in, buf, len, "strip", FALSE, "compression", compression,
		"interlace", interlace, NULL);
#endif
};

int
vips_webpsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace)
{
	return vips_webpsave_buffer(in, buf, len, "strip", strip, "Q", quality, NULL);
};

int
vips_init_image(void *buf, size_t len, int imageType, VipsImage **out) {
	int code = 1;

	if (imageType == JPEG) {
		code = vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == PNG) {
		code = vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == WEBP) {
		code = vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == TIFF) {
		code = vips_tiffload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#if (VIPS_MAJOR_VERSION >= 8)
	} else if (imageType == MAGICK) {
		code = vips_magickload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#endif
	}

	return code;
};

int
vips_watermark_replicate(VipsImage *orig, VipsImage *in, VipsImage **out)
{
	VipsImage *cache = vips_image_new();

	if (
		vips_replicate(in, &cache,
			1 + orig->Xsize / in->Xsize,
			1 + orig->Ysize / in->Ysize, NULL) ||
		vips_crop(cache, out, 0, 0, orig->Xsize, orig->Ysize, NULL)
	) {
		g_object_unref(cache);
		return 1;
	}

	return 0;
};

int
vips_watermark(VipsImage *in, VipsImage **out, WatermarkTextOptions *to, WatermarkOptions *o)
{
	double ones[3] = { 1, 1, 1 };

	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);
	t[0] = in;

	// Make the mask.
	if (
		vips_text(&t[1], to->Text,
			"width", o->Width,
			"dpi", o->DPI,
			"font", to->Font,
			NULL) ||
		vips_linear1(t[1], &t[2], o->Opacity, 0.0, NULL) ||
		vips_cast(t[2], &t[3], VIPS_FORMAT_UCHAR, NULL) ||
		vips_embed(t[3], &t[4], 100, 100, t[3]->Xsize + o->Margin, t[3]->Ysize + o->Margin, NULL)
		) {
		g_object_unref(base);
		return 1;
	}

	// Replicate if necessary
	if (o->NoReplicate != 1) {
		VipsImage *cache = vips_image_new();
		if (vips_watermark_replicate(t[0], t[4], &cache)) {
			g_object_unref(cache);
			g_object_unref(base);
			return 1;
		}
		g_object_unref(t[4]);
		t[4] = cache;
	}

	// Make the constant image to paint the text with.
	if (
		vips_black(&t[5], 1, 1, NULL) ||
		vips_linear(t[5], &t[6], ones, o->Background, 3, NULL) ||
		vips_cast(t[6], &t[7], VIPS_FORMAT_UCHAR, NULL) ||
		vips_copy(t[7], &t[8], "interpretation", t[0]->Type, NULL) ||
		vips_embed(t[8], &t[9], 0, 0, t[0]->Xsize, t[0]->Ysize, "extend", VIPS_EXTEND_COPY, NULL)
		) {
		g_object_unref(base);
		return 1;
	}

	// Blend the mask and text and write to output.
	if (vips_ifthenelse(t[4], t[9], t[0], out, "blend", TRUE, NULL)) {
		g_object_unref(base);
		return 1;
	}

	g_object_unref(base);
	return 0;
};
