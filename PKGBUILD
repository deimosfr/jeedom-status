# Maintainer: Pierre Mavro <deimosfr@gmail.com>
pkgname=jeedom-status
pkgver=tbd
pkgrel=1
pkgdesc="Add Jeedom global status to your favorite desktop bar (i3blocks, polybar, etc...)"
arch=(x86_64)
url="https://github.com/deimosfr/jeedom-status"
license=('GPL')
makedepends=(git cargo)
source=("https://github.com/deimosfr/jeedom-status/archive/v$pkgver.tar.gz")

build() {
	cd "$pkgname-$pkgver"
    cargo build --release
}

package() {
	cd "$pkgname-$pkgver"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
    install -Dm755 "target/release/$pkgname" "$pkgdir/usr/bin/$pkgname"
}
