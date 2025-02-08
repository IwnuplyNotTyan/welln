pkgname=welln
pkgver=1.0.0
pkgrel=1
pkgdesc="Simpel todo manager ~"
arch=('x86_64')
url="https://github.com/iwnuplynottyan/welln"
license=('MIT')
depends=()
makedepends=('go' 'git')
source=("git+https://github.com/iwnuplynottyan/welln.git")
sha256sums=('SKIP')

build() {
  cd "$srcdir/$pkgname"
  go mod tidy
  go build -o welln
}

package() {
  cd "$srcdir/$pkgname"
  install -Dm755 welln "$pkgdir/usr/bin/welln"
}
