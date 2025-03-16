pkgname=hypr-gtk
pkgver=1.0
pkgrel=1
pkgdesc="for hyprland setting"
arch=('x86_64' 'aarch64')
url="https://github.com/falser101/hypr-gtk"
license=('Apache')
makedepends=('go' 'git')
source=("$url/archive/v$pkgver.tar.gz")
sha256sums=('b65f6735423bdc06cb617feccb2ec4087bcc0113c85e5f25089dd18b42b2e6db')

build() {
  cd "$pkgname-$pkgver"
  go build -o $pkgname ./main.go
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 $pkgname "$pkgdir/usr/bin/$pkgname"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
