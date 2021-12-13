FROM golang:1.17
LABEL maintainer="tomas@aparicio.me"

# Prepare the base system
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libaom-dev && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

ENV LD_LIBRARY_PATH="/vips/lib:/usr/local/lib:$LD_LIBRARY_PATH"
ENV PKG_CONFIG_PATH="/vips/lib/pkgconfig:/usr/local/lib/pkgconfig:/usr/lib/pkgconfig:/usr/X11/lib/pkgconfig"

# Setup libheif
ARG LIBHEIF_VERSION=1.12.0
RUN  cd /tmp && \
  curl -fsSLO https://github.com/strukturag/libheif/releases/download/v${LIBHEIF_VERSION}/libheif-${LIBHEIF_VERSION}.tar.gz && \
  tar zvxf libheif-${LIBHEIF_VERSION}.tar.gz && \
  cd /tmp/libheif-${LIBHEIF_VERSION} && \
  ./configure --prefix=/vips --disable-go && \
  make && \
  make install && \
  echo '/vips/lib' > /etc/ld.so.conf.d/vips.conf && \
  ldconfig -v && \
  rm -rf /tmp/*

# Install Go lint
ARG GOLANGCILINT_VERSION=1.43.0
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCILINT_VERSION}

# Setup libvips
ARG LIBVIPS_VERSION=8.9.2
RUN cd /tmp && \
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
    CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" \
    ./configure \
    --disable-debug \
    --disable-dependency-tracking \
    --disable-introspection \
    --disable-static \
    --enable-gtk-doc-html=no \
    --enable-gtk-doc=no \
    --enable-pyvips8=no \
    --prefix=/vips && \
  make && \
  make install && \
  ldconfig && \
  rm -rf /tmp/*

WORKDIR /go/src/app

CMD [ "/bin/bash" ]
