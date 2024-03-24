# KubeColor Themes

If you want to share your custom theme, please open a Pull Request and:

- explain in the PR what is the goal of the theme (Ex: simple esthetic, disability related, better efficiency to display informations,...)
- add a folder for your theme into the `themes` folder, and name it with name of your theme and either the `-dark` or `-light` suffix. Ex: `my-test-theme-dark`.
- a screen capture of the end result as `.jpg` or `.png`
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