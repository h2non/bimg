#include <stdlib.h>
#include <string.h>
#include <vips/vips.h>
#include <vips/foreign.h>
#include <vips/vector.h>

#define VIPS_VERSION_MIN(major, minor) (VIPS_MAJOR_VERSION > major || (VIPS_MAJOR_VERSION == major && VIPS_MINOR_VERSION >= minor))

#if (!VIPS_VERSION_MIN(8, 10))
	#error Unsupported VIPS version
#endif

#define EXIF_IFD0_ORIENTATION "exif-ifd0-Orientation"

#define INT_TO_GBOOLEAN(bool) (bool > 0 ? TRUE : FALSE)


enum types {
	UNKNOWN = 0,
	JPEG,
	WEBP,
	PNG,
	TIFF,
	GIF,
	PDF,
	SVG,
	HEIF,
	AVIF,
	JP2K,
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

typedef struct {
	int    Left;
	int    Top;
	float    Opacity;
} WatermarkImageOptions;

int
vips_init_image (void *buf, size_t len, int imageType, VipsImage **out) {
	switch (imageType) {
		case JPEG:
			return vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case WEBP:
			return vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case PNG:
			return vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case TIFF:
			return vips_tiffload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case GIF:
			return vips_gifload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case PDF:
			return vips_pdfload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case SVG:
			return vips_svgload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
		case HEIF:
		case AVIF:
			return vips_heifload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#if (VIPS_VERSION_MIN(8, 11))
		case JP2K:
			return vips_jp2kload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#endif
		case MAGICK:
			return vips_magickload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	}

	return 1;
}

int
vips_type_find_load_bridge(int t) {
	switch (t) {
		case JPEG:
			return vips_type_find("VipsOperation", "jpegload");
		case WEBP:
			return vips_type_find("VipsOperation", "webpload");
		case PNG:
			return vips_type_find("VipsOperation", "pngload");
		case TIFF:
			return vips_type_find("VipsOperation", "tiffload");
		case GIF:
			return vips_type_find("VipsOperation", "gifload");
		case PDF:
			return vips_type_find("VipsOperation", "pdfload");
		case SVG:
			return vips_type_find("VipsOperation", "svgload");
		case HEIF:
		case AVIF:
			return vips_type_find("VipsOperation", "heifload");
#if (VIPS_VERSION_MIN(8, 11))
		case JP2K:
			return vips_type_find("VipsOperation", "jp2kload");
#endif
		case MAGICK:
			return vips_type_find("VipsOperation", "magickload");
	}
	return 0;
}

int
vips_type_find_save_bridge(int t) {
	switch (t) {
		case JPEG:
			return vips_type_find("VipsOperation", "jpegsave_buffer");
		case WEBP:
			return vips_type_find("VipsOperation", "webpsave_buffer");
		case PNG:
			return vips_type_find("VipsOperation", "pngsave_buffer");
		case TIFF:
			return vips_type_find("VipsOperation", "tiffsave_buffer");
		case GIF:
			return vips_type_find("VipsOperation", "gifsave_buffer");
		case HEIF:
		case AVIF:
			return vips_type_find("VipsOperation", "heifsave_buffer");
#if (VIPS_VERSION_MIN(8, 11))
		case JP2K:
			return vips_type_find("VipsOperation", "jp2ksave_buffer");
#endif
		case MAGICK:
			return vips_type_find("VipsOperation", "magicksave_buffer");
	}
	return 0;
}

int
vips_jpegsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace) {
	return vips_jpegsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"optimize_coding", TRUE,
		"interlace", INT_TO_GBOOLEAN(interlace),
		NULL
	);
}

int
vips_pngsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int compression, int quality, int interlace, int palette, int speed) {
	int effort = 10 - speed;
	return vips_pngsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"compression", compression,
		"interlace", INT_TO_GBOOLEAN(interlace),
		"filter", VIPS_FOREIGN_PNG_FILTER_ALL,
		"palette", INT_TO_GBOOLEAN(palette),
		"Q", quality,
#if (VIPS_VERSION_MIN(8, 12))
		"effort", effort,
#endif
		NULL
	);
}

int
vips_webpsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless) {
	return vips_webpsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		NULL
	);
}

int
vips_tiffsave_bridge(VipsImage *in, void **buf, size_t *len) {
	return vips_tiffsave_buffer(in, buf, len, NULL);
}

int
vips_gifsave_bridge(VipsImage *in, void **buf, size_t *len) {
#if (VIPS_VERSION_MIN(8, 12))
	return vips_gifsave_buffer(in, buf, len, NULL);
#else
	return 0;
#endif
}

int
vips_avifsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless, int speed) {
    return vips_heifsave_buffer(in, buf, len,
	    "strip", INT_TO_GBOOLEAN(strip),
	    "Q", quality,
	    "lossless", INT_TO_GBOOLEAN(lossless),
	    "compression", VIPS_FOREIGN_HEIF_COMPRESSION_AV1,
#if (VIPS_MAJOR_VERSION > 8 || (VIPS_MAJOR_VERSION >= 8 && VIPS_MINOR_VERSION > 10) || (VIPS_MAJOR_VERSION >= 8 && VIPS_MINOR_VERSION >= 10 && VIPS_MICRO_VERSION >= 2))
		"speed", speed,
#endif
	    NULL
    );
}

int
vips_heifsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless) {
	return vips_heifsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		NULL
	);
}

int
vips_magicksave_bridge(VipsImage *in, void **buf, size_t *len, const char *format, int quality) {
	return vips_magicksave_buffer(in, buf, len,
		"format", format,
		"quality", quality,
		NULL
	);
}

int
vips_jp2ksave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless) {
#if (VIPS_VERSION_MIN(8, 11))
	return vips_jp2ksave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		NULL
	);
#else
	return 0;
#endif
}

static unsigned long
has_profile_embed(VipsImage *image) {
	return vips_image_get_typeof(image, VIPS_META_ICC_NAME);
}

static void
remove_profile(VipsImage *image) {
	vips_image_remove(image, VIPS_META_ICC_NAME);
}

static int
has_alpha_channel(VipsImage *image) {
	return (
		(image->Bands == 2 && image->Type == VIPS_INTERPRETATION_B_W) ||
		(image->Bands == 4 && image->Type != VIPS_INTERPRETATION_CMYK) ||
		(image->Bands == 5 && image->Type == VIPS_INTERPRETATION_CMYK)
	) ? 1 : 0;
}

void
vips_enable_cache_set_trace() {
	vips_cache_set_trace(TRUE);
}

int
vips_jpegload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink) {
	return vips_jpegload_buffer(buf, len, out, "shrink", shrink, NULL);
}

int
vips_webpload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink) {
	return vips_webpload_buffer(buf, len, out, "shrink", shrink, NULL);
}

int
vips_flip_bridge(VipsImage *in, VipsImage **out, int direction) {
	return vips_flip(in, out, direction, NULL);
}

int
vips_resize_bridge(VipsImage *in, VipsImage **out, double xscale, double yscale) {
	return vips_resize(in, out, xscale, "vscale", yscale, NULL);
}

int
vips_rotate_bridge(VipsImage *in, VipsImage **out, int angle) {
	angle %= 360;

    // Check, if a special rotation is hit, that allows us to use
    // optimized functions.
    if (angle % 90 == 0) {
    	int rotate = VIPS_ANGLE_D0;
	    switch (angle) {
    		case 90:  rotate = VIPS_ANGLE_D90; break;
    		case 180: rotate = VIPS_ANGLE_D180; break;
    		case 270: rotate = VIPS_ANGLE_D270; break;
    	}
    	return vips_rot(in, out, rotate, NULL);
    }

	return vips_rotate(in, out, angle, NULL);
}

int
vips_autorot_bridge(VipsImage *in, VipsImage **out) {
	return vips_autorot(in, out, NULL);
}

const char *
vips_exif_tag(VipsImage *image, const char *tag) {
	const char *exif;
	if (
		vips_image_get_typeof(image, tag) != 0 &&
		!vips_image_get_string(image, tag, &exif)
	) {
		return &exif[0];
	}
	return "";
}

int
vips_exif_tag_to_int(VipsImage *image, const char *tag) {
	int value = 0;
	const char *exif = vips_exif_tag(image, tag);
	if (strcmp(exif, "")) {
		value = atoi(&exif[0]);
	}
	return value;
}

int
vips_exif_orientation(VipsImage *image) {
	return vips_exif_tag_to_int(image, EXIF_IFD0_ORIENTATION);
}

int
interpolator_window_size(char const *name) {
	VipsInterpolate *interpolator = vips_interpolate_new(name);
	int window_size = vips_interpolate_get_window_size(interpolator);
	g_object_unref(interpolator);
	return window_size;
}

const char *
vips_enum_nick_bridge(VipsImage *image) {
	return vips_enum_nick(VIPS_TYPE_INTERPRETATION, image->Type);
}

int
vips_zoom_bridge(VipsImage *in, VipsImage **out, int xfac, int yfac) {
	return vips_zoom(in, out, xfac, yfac, NULL);
}

int
vips_embed_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, VipsArrayDouble *background) {
	if (extend == VIPS_EXTEND_BACKGROUND) {
		return vips_embed(in, out, left, top, width, height, "extend", extend, "background", background, NULL);
	}
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
}

int
vips_extract_area_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height) {
	return vips_extract_area(in, out, left, top, width, height, NULL);
}

int
vips_colourspace_issupported_bridge(VipsImage *in) {
	return vips_colourspace_issupported(in) ? 1 : 0;
}

VipsInterpretation
vips_image_guess_interpretation_bridge(VipsImage *in) {
	return vips_image_guess_interpretation(in);
}

int
vips_colourspace_bridge(VipsImage *in, VipsImage **out, VipsInterpretation space) {
	return vips_colourspace(in, out, space, NULL);
}

int
vips_icc_transform_bridge (VipsImage *in, VipsImage **out, const char *output_icc_profile) {
	// `output_icc_profile` represents the absolute path to the output ICC profile file
	return vips_icc_transform(in, out, output_icc_profile, "embedded", TRUE, NULL);
}


int
vips_icc_transform_with_default_bridge (VipsImage *in, VipsImage **out, const char *output_icc_profile, const char *input_icc_profile) {
	// `output_icc_profile` represents the absolute path to the output ICC profile file
	return vips_icc_transform(in, out, output_icc_profile, "input_profile", input_icc_profile, "embedded", FALSE, NULL);
}

int
vips_is_16bit (VipsInterpretation interpretation) {
	return interpretation == VIPS_INTERPRETATION_RGB16 || interpretation == VIPS_INTERPRETATION_GREY16;
}

int
vips_flatten_background_bridge(VipsImage *in, VipsImage **out, double r, double g, double b) {
	if (vips_is_16bit(in->Type)) {
		r = 65535 * r / 255;
		g = 65535 * g / 255;
		b = 65535 * b / 255;
	}

	double background[3] = {r, g, b};
	VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);

	return vips_flatten(in, out,
		"background", vipsBackground,
		"max_alpha", vips_is_16bit(in->Type) ? 65535.0 : 255.0,
		NULL
	);
}

int
vips_watermark_replicate (VipsImage *orig, VipsImage *in, VipsImage **out) {
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

	g_object_unref(cache);
	return 0;
}

int
vips_watermark(VipsImage *in, VipsImage **out, WatermarkTextOptions *to, WatermarkOptions *o) {
	double ones[3] = { 1, 1, 1 };

	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);
	g_object_ref(in);
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
}

int
vips_gaussblur_bridge(VipsImage *in, VipsImage **out, double sigma, double min_ampl) {
	return vips_gaussblur(in, out, sigma, NULL, "min_ampl", min_ampl, NULL);
}

int
vips_sharpen_bridge(VipsImage *in, VipsImage **out, int radius, double x1, double y2, double y3, double m1, double m2) {
	return vips_sharpen(in, out, "radius", radius, "x1", x1, "y2", y2, "y3", y3, "m1", m1, "m2", m2, NULL);
}

int
vips_add_band(VipsImage *in, VipsImage **out, double c) {
	return vips_bandjoin_const1(in, out, c, NULL);
}

int
vips_watermark_image(VipsImage *in, VipsImage *sub, VipsImage **out, WatermarkImageOptions *o) {
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);

	g_object_ref(in);
	g_object_ref(sub);

  // add in and sub for unreffing and later use
	t[0] = in;
	t[1] = sub;

  if (has_alpha_channel(in) == 0) {
		vips_add_band(in, &t[0], 255.0);
		// in is no longer in the array and won't be unreffed, so add it at the end
		t[8] = in;
	}

	if (has_alpha_channel(sub) == 0) {
		vips_add_band(sub, &t[1], 255.0);
		// sub is no longer in the array and won't be unreffed, so add it at the end
		t[9] = sub;
	}

	// Place watermark image in the right place and size it to the size of the
	// image that should be watermarked
	if (
		vips_embed(t[1], &t[2], o->Left, o->Top, t[0]->Xsize, t[0]->Ysize, NULL)) {
			g_object_unref(base);
		return 1;
	}

	// Create a mask image based on the alpha band from the watermark image
	// and place it in the right position
	if (
		vips_extract_band(t[1], &t[3], t[1]->Bands - 1, "n", 1, NULL) ||
		vips_linear1(t[3], &t[4], o->Opacity, 0.0, NULL) ||
		vips_cast(t[4], &t[5], VIPS_FORMAT_UCHAR, NULL) ||
		vips_copy(t[5], &t[6], "interpretation", t[0]->Type, NULL) ||
		vips_embed(t[6], &t[7], o->Left, o->Top, t[0]->Xsize, t[0]->Ysize, NULL))	{
			g_object_unref(base);
		return 1;
	}

	// Blend the mask and watermark image and write to output.
	if (vips_ifthenelse(t[7], t[2], t[0], out, "blend", TRUE, NULL)) {
		g_object_unref(base);
		return 1;
	}

	g_object_unref(base);
	return 0;
}

int
vips_smartcrop_bridge(VipsImage *in, VipsImage **out, int width, int height) {
	return vips_smartcrop(in, out, width, height, NULL);
}

int vips_find_trim_bridge(VipsImage *in, int *top, int *left, int *width, int *height, double r, double g, double b, double threshold) {
	if (vips_is_16bit(in->Type)) {
		r = 65535 * r / 255;
		g = 65535 * g / 255;
		b = 65535 * b / 255;
	}

	double background[3] = {r, g, b};
	VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);
	return vips_find_trim(in, top, left, width, height, "background", vipsBackground, "threshold", threshold, NULL);
}

int vips_gamma_bridge(VipsImage *in, VipsImage **out, double exponent) {
	return vips_gamma(in, out, "exponent", 1.0 / exponent, NULL);
}

int vips_addalpha_bridge(VipsImage *in, VipsImage **out) {
	return vips_addalpha(in, out);
}
