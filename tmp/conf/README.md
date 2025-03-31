Add DEB file to repo:

```bash
reprepro includedeb stable ../dist/kubecolor_0.5.0~SNAPSHOT-c1a9da1_i386.deb
```

Sample install docs for Debian:

```bash
echo "deb [trusted=yes arch=$(dpkg --print-architecture)] http://localhost:8000/ stable main" \
  | tee /etc/apt/sources.list.d/kubecolor.list
```
