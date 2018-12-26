FROM ubuntu:16.04 as builder

# Specify using --build-arg LIBVIPS_VERSION=1.2.3
ARG LIBVIPS_VERSION="8.7.0"
ARG BIMG_VERSION="dev"
ARG GO_VERSION="1.11.2"

# Setting up a sane build environment
RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates curl \
  automake build-essential gcc git libc6-dev make \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg-turbo8-dev libpng12-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev

# Build libvips
RUN cd /tmp && \
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
	CXXFLAGS=-D_GLIBCXX_USE_CXX11_ABI=0 \
    ./configure \
    --disable-debug \
    --disable-dependency-tracking \
    --disable-introspection \
    --disable-static \
    --enable-gtk-doc-html=no \
    --enable-gtk-doc=no \
    --enable-pyvips8=no \
    --without-orc \
    --without-python && \
  make && \
  make install && \
  ldconfig

# Installing Go
ENV GO_DOWNLOAD_URL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz

RUN curl -fsSL "${GO_DOWNLOAD_URL}" -o golang.tar.gz \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

# Install QA tools
WORKDIR /tmp
RUN curl -fsSL https://git.io/vp6lP -o instgm.sh && chmod u+x instgm.sh && ./instgm.sh -b "${GOPATH}/bin"

# Setup our Go environment
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:${PATH}


WORKDIR /go/src/github.com/h2non/bimg/

# Making sure all dependencies are up-to-date
RUN GO111MODULE=off go get -u github.com/golang/dep/cmd/dep

# Copying bimg
COPY . .

RUN dep ensure

# Run quality control
#RUN go test -test.v ./...
RUN GO111MODULE=off gometalinter github.com/h2non/bimg

# Compile the binary, to verify compile-time correctness. The build should fail if this step fails.
RUN GO111MODULE=off go build -a \
    -ldflags="-h -X bimg.Version=${BIMG_VERSION}" \
    github.com/h2non/bimg

# Clean up
RUN DEBIAN_FRONTEND=noninteractive apt-get remove -y automake build-essential && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

FROM ubuntu:16.04

ARG GO_VERSION
ARG LIBVIPS_VERSION
ARG BIMG_VERSION

# Exposing version information in parent images
ONBUILD ENV GO_VERSION ${GO_VERSION}
ONBUILD ENV LIBVIPS_VERSION ${LIBVIPS_VERSION}
ONBUILD ENV BIMG_VERSION ${BIMG_VERSION}

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:${PATH}

# @todo optimise a bit so that resulting image is smaller.
# Squashing the result into a single layer
COPY --from=builder / /

LABEL maintainer="tomas@aparicio.me" \
			org.label-schema.description="Small Go package for fast high-level image processing powered by libvips C library" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/h2non/bimg" \
      org.label-schema.vcs-url="https://github.com/h2non/bimg" \
      org.label-schema.version="${BIMG_VERSION}" \
      libvips.version="${LIBVIPS_VERSION}" \
      go.version="${GO_VERSION}"


WORKDIR /go/src/github.com/h2non/bimg/
