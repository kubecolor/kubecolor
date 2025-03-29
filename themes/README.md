# KubeColor Themes

If you want to share your custom theme, please open a Pull Request and:

- explain in the PR what is the goal of the theme (Ex: simple esthetic, disability related, better efficiency to display information,...)
- add a folder for your theme into the `themes` folder, and name it with name of your theme and either the `-dark` or `-light` suffix. Ex: `my-test-theme-dark`.
- a screenshot of the end result as `.jpg` or `.png`
- the `color.yaml` file that starts with comments with:
  - your name
  - your email
  - the date of the last update
  then all the variables you changed. You can remove all the comments and un-changed values.

Final structure should look like:

```
themes
├── README.md
├── my-test-theme-dark
│   ├── color.yaml
│   └── image.png
└── new-theme-dark
    ├── color.yaml
    └── image.png
```

## Creating screenshot

Requires that you have kubectl and kubecolor installed, as well as access to a Kubernetes cluster (e.g via Docker Desktop, Kind, K3s, Minikube).

Apply the files found in [../test/cluster](../test/cluster), like so:

```bash
# run from the root of the project
kubectl create -f ./test/cluster
```

Type the `kubecolor` commands that can show how your theme changes from the default theme. Ex:

```bash
kubecolor get pods -o wide -n kubecolor
kubecolor get pod -l app=working-pods  -o yaml -n kubecolor
kubecolor describe pod -l app=working-pods -n kubecolor
kubecolor get pod -l app=ailing-pods -o yaml -n kubecolor
```

Then take a screenshot of your terminal.
When you're done, you can delete the resources again:

```bash
kubectl delete -f ./test/cluster
```
