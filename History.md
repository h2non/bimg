
v1.1.1 / 2020-06-08
===================

  * feat(version): bump patch
  * refactor(docs): add libvips install reference
  * fix(ci): disable old libvips versions
  * fix(install): use latest libvips version
  * fix(tests): add heif exception in libvips < 8.8
  * refactor(ci): use libvips 8.7
  * fix(History): use proper version


v1.1.0 / 2020-06-07
===================

  * feat(ci): enable libvips versions
  * fix(ci)
  * fix(ci)
  * fix(ci): try exporting env vars
  * fix
  * feat: add Dockerfile / Docker-driven CI job
  * fix(co)
  * feat(version): bump minor to 1
  * fix(ci): try new install
  * fix(ci): try new install
  * fix(ci): add curl package
  * fix(ci): add curl package
  * fix(ci): add curl package
  * fix(ci): try new install
  * fix(ci): indent style
  * fix(ci): indent style
  * fix(ci): indent style
  * Merge pull request #299 from evanoberholster/master
  * refactor(ci): disable verions matrix
  * refactor(docs): use github.com package import path
  * feat: add test image
  * Merge pull request #281 from pohang/skip_smartcrop
  * Merge pull request #317 from larrabee/master
  * Merge pull request #307 from OrderMyGear/eslam/ch15924/some-product-images-have-a-border
  * refactor(travis): adjust matrix versions
  * Merge pull request #333 from simia-tech/master
  * Fix orientation in vipsFlip call (resizer rotateAndFlipImage)
  * chore(docs): delete old contributor
  * enable vipsAffine to use  `Extend` option value and send it to lipvips this will change the default from the one that lipvips use which is `background` to the ones that bimg use which is  `C.VIPS_EXTEND_BLACK` but because the lip add extra 1 or .5 pix the background is considered black anyway so this will not affect anyone but will fix the bug of having border on the right and bottom of some images
  * Merge pull request #327 from shoreward/master
  * update libvips documentation links
  * fix(vips.h): delete preprocessor HEIF version check
  * Merge pull request #320 from cgroschupp/feat/reduce-png-save-size
  * use VIPS_FOREIGN_PNG_FILTER_ALL in vips_pngsave_bridge
  * fix(resizer): add exported error comment
  * Merge branch 'master' of https://github.com/h2non/bimg
  * chore(ci): temporarily disable go/libvips versions
  * Merge pull request #291 from andrioid/patch-1
  * Merge pull request #293 from team-lab/gammaFilter
  * Merge pull request #315 from vansante/heif
  * feat(version): bump patch
  * Fix bug with images with alpha channel on embeding background
  * Fix typo
  * Dont upgrade version, add missing test file
  * Add support for other HEIF mimetype
  * Supporting auto rotate for HEIF/HEIC images.
  * Adding support for heif (i.e. heic files).
  * Merge branch 'master' into master
  * feat(travis): add libvips 8.6.0 matrix
  * GammaFilter
  * Adds support to Elementary OS Loki
  * Add min dimension logic to smartcrop
  * Merge pull request #271 from Dynom/ImprovingAreaWidthTestCoverage
  * Adding a test case that verifies #250
  * Bumping versions in preinstall script
  * Update Transform ICC Profiles with Input Profile

## v1.0.18 / 2017-12-22

  * Merge pull request #216 from Bynder/master
  * Merge pull request #208 from mikestead/feature/webp-lossless
  * Remove go-debug usage
  * refactor(docs): remove codesponsor :(
  * fix(options): use float64 type in Options.Threshold
  * Merge pull request #206 from tstm/add-trim-options
  * Add lossless option for saving webp
  * Set the test file to write its own file
  * Add the option to use background and threshold options on trim

## v1.0.17 / 2017-11-14

  * refactor(resizer): remove fmt statement
  * fix(type_test): use string formatting
  * Merge pull request #207 from traum-ferienwohnungen/nearest-neighbour
  * Add nearest-neighbour interpolation
  * Merge pull request #203 from traum-ferienwohnungen/fix_icc_memory_leak
  * Fix memory leak on icc_transform

## v1.0.16 / 2017-10-30

  * fix(travis): use install directive
  * Merge branch 'master' of https://github.com/h2non/bimg
  * feat: add Gopkg manifests, move fixtures to testdata, add vendor dependencies
  * Merge pull request #202 from openskydoor/openskydoor/fix-build-tag
  * fix build tag
  * fix(#199): presinstall.sh tarball download URL

## v1.0.15 / 2017-10-05

  * Merge pull request #198 from greut/webpload
  * Add shrink-on-load for webp.
  * Merge pull request #197 from greut/typos
  * Small typo.
  * feat(docs): add codesponsor

## v1.0.14 / 2017-09-12

  * Merge pull request #192 from greut/trim
  * Adding trim operation.
  * Merge pull request #191 from greut/alpha4
  * Update 8.6 to alpha4.

## v1.0.13 / 2017-09-11

  * Merge pull request #190 from greut/typos
  * Fix typo and small cleanup.

## v1.0.12 / 2017-09-10

  * Merge branch '99designs-vips-reduce'
  * fix(reduce): resolve conflicts with master
  * Use vips reduce when downscaling

## v1.0.11 / 2017-09-10

  * feat(#189): allow strip image metadata via bimg.Options.StripMetadata = bool
  * fix(resize): code format issue
  * refactor(resize): add Go version comment
  * refactor(tests): fix minor code formatting issues
  * fix(#162): garbage collection fix. split Resize() implementation for Go runtime specific
  * feat(travis): add go 1.9
  * Merge pull request #183 from greut/autorotate
  * Proper handling of the EXIF cases.
  * Merge pull request #184 from greut/libvips858
  * Merge branch 'master' into libvips858
  * Merge pull request #185 from greut/libvips860
  * Add libvips 8.6 pre-release
  * Update to libvips 8.5.8
  * fix(resize): runtime.KeepAlive is only Go
  * fix(#159): prevent buf to be freed by the GC before resize function exits
  * Merge pull request #171 from greut/fix-170
  * Check the length before jumping into buffer.
  * Merge pull request #168 from Traum-Ferienwohnungen/icc_transform
  * Add option to convert embedded ICC profiles
  * Merge pull request #166 from danjou-a/patch-1
  * Fix Resize verification value
  * Merge pull request #165 from greut/libvips846
  * Testing using libvips8.4.6 from Github.

## v1.0.10 / 2017-06-25

  * Merge pull request #164 from greut/length
  * Add Image.Length()
  * Merge pull request #163 from greut/libvips856
  * Run libvips 8.5.6 on Travis.
  * Merge pull request #161 from henry-blip/master
  * Expose vips cache memory management functions.
  * feat(docs): add watermark image note in features

## v1.0.9 / 2017-05-25

  * Merge pull request #156 from Dynom/SmartCropToGravity
  * Adding a test, verifying both ways of enabling SmartCrop work
  * Merge pull request #149 from waldophotos/master
  * Replacing SmartCrop with a Gravity option
  * refactor(docs): v8.4
  * Change for older LIBVIPS versions. `vips_bandjoin_const1` is added in libvips 8.2.
  * Second try, watermarking memory issue fix

## v1.0.8 / 2017-05-18

  * Merge pull request #145 from greut/smartcrop
  * Merge pull request #155 from greut/libvips8.5.5
  * Update libvips to 8.5.5.
  * Adding basic smartcrop support.
  * Merge pull request #153 from abracadaber/master
  * Added Linux Mint 17.3+ distro names
  * feat(docs): add new maintainer notice (thanks to @kirillDanshin)
  * Merge pull request #152 from greut/libvips85
  * Download latest version of libvips from github.
  * Merge pull request #147 from h2non/revert-143-master
  * Revert "Fix for memory issue when watermarking images"
  * Merge pull request #146 from greut/minor-major
  * Merge pull request #143 from waldophotos/master
  * Merge pull request #144 from greut/go18
  * Fix tests where minor/major were mixed up
  * Enabled go 1.8 builds.
  * Fix the unref of images, when image isn't transparent
  * Fix for memory issue when watermarking images
  * feat(docs): add maintainers sections
  * Merge pull request #132 from jaume-pinyol/WATERMARK_SUPPORT
  * Add support for image watermarks
  * Merge pull request #131 from greut/versions
  * Running tests on more specific versions.
  * refactor(preinstall.sh): remove deprecation notice
  * Update preinstall.sh
  * fix(requirements): required libvips 7.42
  * fix(History): typo
  * chore(History): add breaking change note

## v1.0.7 / 13-01-2017

- fix(#128): crop image calculation for missing width or height axis.
- feat: add TIFF save output format (**note**: this introduces a minor interface breaking change in `bimg.IsImageTypeSupportedByVips` auxiliary function).

## v1.0.6 / 12-11-2016

- feat(#118): handle 16-bit PNGs.
- feat(#119): adds JPEG2000 file for the type tests.
- feat(#121): test bimg against multiple libvips versions.

## v1.0.5 / 01-10-2016

- feat(#92): support Extend param with optional background.
- fix(#106): allow image area extraction without explicit x/y axis.
- feat(api): add Extend type with `libvips` enum alias.

## v1.0.4 / 29-09-2016

- fix(#111): safe check of magick image type support.

## v1.0.3 / 28-09-2016

- fix(#95): better image type inference and support check.
- fix(background): pass proper background RGB color for PNG image conversion.
- feat(types): validate supported image types by current `libvips` compilation.
- feat(types): consistent SVG image checking.
- feat(api): add public functions `VipsIsTypeSupported()`, `IsImageTypeSupportedByVips()` and `IsSVGImage()`.

## v1.0.2 / 27-09-2016

- feat(#95): support GIF, SVG and PDF formats.
- fix(#108): auto-width and height calculations now round instead of floor.

## v1.0.1 / 22-06-2016

- fix(#90): Do not not dereference the original image a second time.

## v1.0.0 / 21-04-2016

- refactor(api): breaking changes: normalize public members to follow Go naming idioms.
- feat(version): bump to major version. API contract won't be compromised in `v1`.
- feat(docs): add missing inline godoc documentation.
