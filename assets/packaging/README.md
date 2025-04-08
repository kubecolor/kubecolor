# Packaging kubecolor

## Debian/Ubuntu

Build:

```bash
goreleaser release --skip=publish --clean --snapshot
assets/packaging/reprepro-pkg.sh
```

Host test server for `site/`:

```bash
cd site
python3 -m http.server -d site
```

Install instructions:

```bash
REPO_HOST=http://localhost:8000/packages
#REPO_HOST=https://kubecolor.github.io/kubecolor-packages

sudo apt-get update
sudo apt-get install apt-transport-https wget --yes
wget -O /tmp/kubecolor.deb https://kubecolor.github.io/kubecolor-packages/deb/pool/main/k/kubecolor/kubecolor_$(wget -q -O- https://kubecolor.github.io/kubecolor-packages/deb/version)_$(dpkg --print-architecture).deb
sudo dpkg -i /tmp/kubecolor.deb
sudo apt update
```

The install instructions can be used in `docker.io/library/ubuntu` or
`docker.io/library/debian` Docker images to try install `kubecolor`

## Fedora

Build:

```bash
goreleaser release --skip=publish --clean --snapshot
assets/packaging/createrepo-pkg.sh
```

Host test server for `site/`:

```bash
cd site/
python3 -m http.server -d site
```

Install instructions:

### DNF5/Fedora 41 and above

```bash
sudo dnf install dnf5-plugins
sudo dnf config-manager addrepo --from-repofile https://localhost:8000/packages/rpm/kubecolor.repo
sudo dnf install kubecolor
```

### DNF4/Fedora 40 and below

```bash
sudo dnf install 'dnf-command(config-manager)'
sudo dnf config-manager --addrepo http://localhost:8000/packages/rpm/kubecolor.repo
sudo dnf install kubecolor
```

### openSUSE/SUSE Linux

```bash
sudo zypper addrepo http://localhost:8000/packages/rpm/kubecolor.repo
sudo zypper install kubecolor
```
