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
python3 -m http.server
```

Install instructions:

```bash
sudo apt-get install apt-transport-https wget --yes
wget -O /tmp/kubecolor.deb localhost:8000/packages/deb/pool/main/k/kubecolor/kubecolor_0.5.0~SNAPSHOT-c14790a_$(dpkg --print-architecture).deb
sudo dpkg -i /tmp/kubecolor.deb
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
python3 -m http.server
```

Install instructions:

```bash
sudo dnf config-manager addrepo --from-repofile http://localhost:8000/packages/rpm/kubecolor.repo
sudo dnf install kubecolor
```
