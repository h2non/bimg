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

VipsImage*
vips_image_buffer_seq(void *buf, size_t len)
{
	// todo: handle postclose callback
	return vips_image_new_from_buffer(buf, len, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
};

int
vips_jpegload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink)
{
	return vips_jpegload_buffer(buf, len, out, "shrink", shrink, NULL);
};

int
vips_webpload_buffer_seq(void *buf, size_t len, VipsImage **out)
{
	return vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
};

int
vips_flip_seq(VipsImage *in, VipsImage **out)
{
	return vips_flip(in, out, VIPS_DIRECTION_HORIZONTAL, NULL);
};

int
vips_pngload_buffer_seq(void *buf, size_t len, VipsImage **out)
{
	return vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
};

int
vips_shrink_0(VipsImage *in, VipsImage **out, double xshrink, double yshrink)
{
	return vips_shrink(in, out, xshrink, yshrink, NULL);
};

int
vips_copy_0(VipsImage *in, VipsImage **out)
{
	return vips_copy(in, out, NULL);
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

  if (out != NULL) {
  	// Listen for "postclose" signal to delete input buffer
  	//g_signal_connect(out, "postclose", G_CALLBACK(vips_malloc_cb), buf);
	}

  return code;
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
vips_embed_extend(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend)
{
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
};

int
vips_colourspace_0(VipsImage *in, VipsImage **out, VipsInterpretation space)
{
	return vips_colourspace(in, out, space, NULL);
};

int
vips_extract_area_0(VipsImage *in, VipsImage **out, int left, int top, int width, int height)
{
	return vips_extract_area(in, out, left, top, width, height, NULL);
};

int
vips_jpegsave_custom(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace)
{
	return vips_jpegsave_buffer(in, buf, len, "strip", strip, "Q", quality, "optimize_coding", TRUE, "interlace", interlace, NULL);
};
