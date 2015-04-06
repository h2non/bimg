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

void
vips_malloc_cb(VipsObject *object, char *buf)
{
	g_free(buf);
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
vips_shrink_0(VipsImage *in, VipsImage **out, double xshrink, double yshrink)
{
	return vips_shrink(in, out, xshrink, yshrink, NULL);
};

int
vips_rotate(VipsImage *in, VipsImage **buf, int angle)
{
	int rotate = VIPS_ANGLE_D0;

	if (angle == 90) {
		rotate = VIPS_ANGLE_D90;
	} else if (angle == 180) {
		rotate = VIPS_ANGLE_D180;
	} else if (angle == 270) {
		rotate = VIPS_ANGLE_D270;
	}

	return vips_rot(in, buf, rotate, NULL);
};

int
vips_exif_orientation(VipsImage *image) {
	int orientation = 0;
	const char **exif;
	if (
		vips_image_get_typeof(image, "exif-ifd0-Orientation") != 0 &&
		!vips_image_get_string(image, "exif-ifd0-Orientation", exif)
	) {
		orientation = atoi(exif[0]);
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
vips_embed_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend)
{
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
};

int
vips_colourspace_bridge(VipsImage *in, VipsImage **out, VipsInterpretation space)
{
	return vips_colourspace(in, out, space, NULL);
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
	return vips_pngsave_buffer(image, buf, len, "strip", FALSE, "compression", compression,
		"interlace", interlace, NULL);
#endif
};

int
vips_webpsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace)
{
	return vips_webpsave_buffer(in, buf, len, "strip", strip, "Q", quality, "optimize_coding", TRUE, "interlace", interlace, NULL);
};

int
vips_init_image(void *buf, size_t len, int imageType, VipsImage **out) {
	int code = 1;

	if (imageType == JPEG) {
		code = vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imageType == PNG) {
		code = vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imageType == WEBP) {
		code = vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imageType == TIFF) {
		code = vips_tiffload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
#if (VIPS_MAJOR_VERSION >= 8)
	} else if (imageType == MAGICK) {
		code = vips_magickload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
#endif
	}

	// Listen for "postclose" signal to delete input buffer
	//if (out != NULL) {
		//g_signal_connect(out, "postclose", G_CALLBACK(vips_malloc_cb), buf);
	//}

	return code;
};
