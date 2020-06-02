# Maintainer: Pierre Mavro <deimosfr@gmail.com>
pkgname=jeedom-status
pkgver=tbd
pkgrel=1
pkgdesc="Add Jeedom global status to your favorite desktop bar (i3blocks, polybar, etc...)"
arch=(x86_64)
url="https://github.com/deimosfr/jeedom-status"
license=('GPL')
makedepends=(git go)
source=("https://github.com/deimosfr/jeedom-status/archive/v$pkgver.tar.gz")

build() {
	cd "$pkgname-$pkgver"
    export CGO_LDFLAGS="${LDFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
    go build -o jeedom-status main.go
}

package() {
	cd "$pkgname-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
