Add DEB file to repo:

```bash
assets/packaging/reprepro.sh
```

Host test server for `dist/deb`:

```bash
tmp/serve.py
```

Install instructions:

```bash
sudo apt-get install apt-transport-https wget --yes
wget -O /tmp/kubecolor.deb localhost:8000/pool/main/k/kubecolor/kubecolor_0.5.0~SNAPSHOT-c14790a_$(dpkg --print-architecture).deb
dpkg -i /tmp/kubecolor.deb
```
