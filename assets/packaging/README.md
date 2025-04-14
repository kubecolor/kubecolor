# Packaging kubecolor

## Debian/Ubuntu

Build:

```bash
goreleaser release --skip=publish --clean --snapshot
assets/packaging/reprepro-pkg.sh
assets/packaging/version.sh
```

Host test server for `packages/`:

```bash
python3 -m http.server
```

### Install instructions

```bash
REPO_HOST=http://localhost:8000/packages
#REPO_HOST=https://kubecolor.github.io/packages

sudo apt-get update
sudo apt-get install apt-transport-https wget --yes
wget -O /tmp/kubecolor.deb $REPO_HOST/deb/pool/main/k/kubecolor/kubecolor_$(wget -q -O- $REPO_HOST/deb/version)_$(dpkg --print-architecture).deb
sudo dpkg -i /tmp/kubecolor.deb
sudo apt update
```

The install instructions can be used in `docker.io/library/ubuntu` or
`docker.io/library/debian` Docker images to try install `kubecolor`

## Fedora

Build:

```bash
goreleaser release --skip=publish --clean --snapshot
assets/packaging/rpmsign-all.sh
assets/packaging/createrepo-pkg.sh
assets/packaging/rpm-repomd-sign.sh
```

Host test server for `packages/`:

```bash
python3 -m http.server
```

### Install instructions

#### DNF5/Fedora 41 and above

```bash
sudo dnf install dnf5-plugins
sudo dnf config-manager addrepo --from-repofile http://localhost:8000/packages/rpm/kubecolor.repo
sudo dnf install kubecolor
```

#### DNF4/Fedora 40 and below

```bash
sudo dnf install 'dnf-command(config-manager)'
sudo dnf config-manager --addrepo http://localhost:8000/packages/rpm/kubecolor.repo
sudo dnf install kubecolor
```

#### openSUSE/SUSE Linux

```bash
sudo zypper addrepo http://localhost:8000/packages/rpm/kubecolor.repo
sudo zypper install kubecolor
```
