#include <stdlib.h>
#include <string.h>
#include <vips/vips.h>
#include <vips/foreign.h>
#include <vips/vips7compat.h>

/**
 * Starting libvips 7.41, VIPS_ANGLE_x has been renamed to VIPS_ANGLE_Dx
 * "to help python". So we provide the macro to correctly build for versions
 * before 7.41.x.
 * https://github.com/jcupitt/libvips/blob/master/ChangeLog#L128
 */

#if (VIPS_MAJOR_VERSION == 7 && VIPS_MINOR_VERSION < 41)
#define VIPS_ANGLE_D0 VIPS_ANGLE_0
#define VIPS_ANGLE_D90 VIPS_ANGLE_90
#define VIPS_ANGLE_D180 VIPS_ANGLE_180
#define VIPS_ANGLE_D270 VIPS_ANGLE_270
#endif

#define EXIF_IFD0_ORIENTATION "exif-ifd0-Orientation"

#define INT_TO_GBOOLEAN(bool) (bool > 0 ? TRUE : FALSE)

#if (VIPS_MAJOR_VERSION <= 8 && VIPS_MINOR_VERSION < 6)
typedef enum {
	VIPS_BLEND_MODE_CLEAR,
	VIPS_BLEND_MODE_SOURCE,
	VIPS_BLEND_MODE_OVER,
	VIPS_BLEND_MODE_IN,
	VIPS_BLEND_MODE_OUT,
	VIPS_BLEND_MODE_ATOP,
	VIPS_BLEND_MODE_DEST,
	VIPS_BLEND_MODE_DEST_OVER,
	VIPS_BLEND_MODE_DEST_IN,
	VIPS_BLEND_MODE_DEST_OUT,
	VIPS_BLEND_MODE_DEST_ATOP,
	VIPS_BLEND_MODE_XOR,
	VIPS_BLEND_MODE_ADD,
	VIPS_BLEND_MODE_SATURATE,
	VIPS_BLEND_MODE_MULTIPLY,
	VIPS_BLEND_MODE_SCREEN,
	VIPS_BLEND_MODE_OVERLAY,
	VIPS_BLEND_MODE_DARKEN,
	VIPS_BLEND_MODE_LIGHTEN,
	VIPS_BLEND_MODE_COLOUR_DODGE,
	VIPS_BLEND_MODE_COLOUR_BURN,
	VIPS_BLEND_MODE_HARD_LIGHT,
	VIPS_BLEND_MODE_SOFT_LIGHT,
	VIPS_BLEND_MODE_DIFFERENCE,
	VIPS_BLEND_MODE_EXCLUSION,
	VIPS_BLEND_MODE_LAST
} VipsBlendMode;
#endif

enum types {
	UNKNOWN = 0,
	JPEG,
	WEBP,
	PNG,
	TIFF,
	GIF,
	PDF,
	SVG,
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

#if (VIPS_MAJOR_VERSION == 7 && VIPS_MINOR_VERSION < 41)
#undef vips_init
int
vips_init(const char *argv0);
#endif

void
vips_enable_cache_set_trace();

int
vips_affine_interpolator(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator);

int
vips_jpegload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink);

int
vips_webpload_buffer_shrink(void *buf, size_t len, VipsImage **out, int shrink);

int
vips_flip_bridge(VipsImage *in, VipsImage **out, int direction);

int
vips_shrink_bridge(VipsImage *in, VipsImage **out, double xshrink, double yshrink);

int
vips_reduce_bridge(VipsImage *in, VipsImage **out, double xshrink, double yshrink);

int
vips_type_find_bridge(int t);

int
vips_type_find_save_bridge(int t);

int
vips_rotate(VipsImage *in, VipsImage **out, int angle);

int
vips_exif_orientation(VipsImage *image);

int
interpolator_window_size(char const *name);

const char *
vips_enum_nick_bridge(VipsImage *image);

int
vips_zoom_bridge(VipsImage *in, VipsImage **out, int xfac, int yfac);

int
vips_composite_bridge(VipsImage **in, VipsImage **out, int n, int mode);

int
vips_embed_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, double r, double g, double b);

int
vips_extract_area_bridge(VipsImage *in, VipsImage **out, int left, int top, int width, int height);

int
vips_colourspace_issupported_bridge(VipsImage *in);

VipsInterpretation
vips_image_guess_interpretation_bridge(VipsImage *in);

int
vips_colourspace_bridge(VipsImage *in, VipsImage **out, VipsInterpretation space);

int
vips_icc_transform_bridge (VipsImage *in, VipsImage **out, const char *output_icc_profile);

int
vips_jpegsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace);

int
vips_pngsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int compression, int quality, int interlace);

int
vips_webpsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless);

int
vips_tiffsave_bridge(VipsImage *in, void **buf, size_t *len);

int
vips_is_16bit (VipsInterpretation interpretation);

int
vips_flatten_background_brigde(VipsImage *in, VipsImage **out, double r, double g, double b);

int
vips_init_image (void *buf, size_t len, int imageType, VipsImage **out);

int
vips_watermark_replicate (VipsImage *orig, VipsImage *in, VipsImage **out);

int
vips_watermark(VipsImage *in, VipsImage **out, WatermarkTextOptions *to, WatermarkOptions *o);

int
vips_gaussblur_bridge(VipsImage *in, VipsImage **out, double sigma, double min_ampl) ;

int
vips_sharpen_bridge(VipsImage *in, VipsImage **out, int radius, double x1, double y2, double y3, double m1, double m2);

int
vips_add_band(VipsImage *in, VipsImage **out, double c);

int
vips_watermark_image(VipsImage *in, VipsImage *sub, VipsImage **out, WatermarkImageOptions *o);

int
vips_smartcrop_bridge(VipsImage *in, VipsImage **out, int width, int height);

int
vips_find_trim_bridge(VipsImage *in, int *top, int *left, int *width, int *height, double r, double g, double b, double threshold);
