pkgname=hypr-gtk
pkgver=1.0
pkgrel=1
pkgdesc="for hyprland setting"
arch=('x86_64' 'aarch64')
url="https://github.com/falser101/hypr-gtk"
license=('Apache')
makedepends=('go' 'git')
source=("$url/archive/v$pkgver.tar.gz")
sha256sums=('c2b19de76da86f82f7019794851577e6f6265cc7950be34aae4e209564f98b51')

build() {
  cd "$pkgname-$pkgver"
  go build -o $pkgname ./main.go
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 $pkgname "$pkgdir/usr/bin/$pkgname"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
