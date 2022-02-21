FROM alpine:3.15
LABEL org.opencontainers.image.authors="andreas@staffbase.com"

RUN mkdir /build

RUN apk add --no-cache \
    curl \
    g++ \
    meson \
    ninja \
    make \
    cfitsio-dev \
    expat-dev \
    fftw-dev \
    giflib-dev \
    glib-dev \
    gobject-introspection-dev \
    lcms2-dev \
    libexif-dev \
    libheif-dev \
    libimagequant-dev \
    libjpeg-turbo-dev \
    libpng-dev \
    libwebp-dev \
    openexr-dev \
    openjpeg-dev \
    orc-dev \
    pango-dev \
    poppler-dev \
    tiff-dev \
    imagemagick-dev \
    librsvg-dev

ARG CGIF_VERSION=0.2.0
RUN cd /build && \
    curl -fsSLO https://github.com/dloebl/cgif/archive/refs/tags/V${CGIF_VERSION}.tar.gz && \
    tar xf V${CGIF_VERSION}.tar.gz && \
    cd cgif-${CGIF_VERSION} && \
    meson . build && \
    meson compile -C build && \
    meson install --no-rebuild -C build

ARG GOLANGCILINT_VERSION=1.44.2
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCILINT_VERSION}

ARG VIPS_VERSION=8.12.2
RUN cd /build && \
    curl -fsSLO https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz && \
    tar xf vips-${VIPS_VERSION}.tar.gz && \
    cd vips-${VIPS_VERSION} && \
    ./configure --enable-debug=no --prefix=/usr --disable-static --enable-introspection && \
    make -j 8 install

RUN cd / && rm -rf /build
