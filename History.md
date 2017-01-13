## 1.0.7 / 13-01-2017

- fix(#128): crop image calculation for missing width or height axis.
- feat: add TIFF save output format (**note**: this introduces a minor interface breaking change in `bimg.IsImageTypeSupportedByVips` auxiliary function).

## 1.0.6 / 12-11-2016

- feat(#118): handle 16-bit PNGs.
- feat(#119): adds JPEG2000 file for the type tests.
- feat(#121): test bimg against multiple libvips versions.

## 1.0.5 / 01-10-2016

- feat(#92): support Extend param with optional background.
- fix(#106): allow image area extraction without explicit x/y axis.
- feat(api): add Extend type with `libvips` enum alias.

## 1.0.4 / 29-09-2016

- fix(#111): safe check of magick image type support.

## 1.0.3 / 28-09-2016

- fix(#95): better image type inference and support check.
- fix(background): pass proper background RGB color for PNG image conversion.
- feat(types): validate supported image types by current `libvips` compilation.
- feat(types): consistent SVG image checking.
- feat(api): add public functions `VipsIsTypeSupported()`, `IsImageTypeSupportedByVips()` and `IsSVGImage()`.

## 1.0.2 / 27-09-2016

- feat(#95): support GIF, SVG and PDF formats.
- fix(#108): auto-width and height calculations now round instead of floor.

## 1.0.1 / 22-06-2016

- fix(#90): Do not not dereference the original image a second time.

## 1.0.0 / 21-04-2016

- refactor(api): breaking changes: normalize public members to follow Go naming idioms.
- feat(version): bump to major version. API contract won't be compromised in `v1`.
- feat(docs): add missing inline godoc documentation.
