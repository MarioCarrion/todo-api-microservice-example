# Tools as Dependencies

This follows the pattern described in [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/g_5n0W27XcY).

Although not required I encourage you to install [`direnv`](https://direnv.net/) to keep your tools _"sandboxed"_, that way in case your use case requires it, you can install different versions of the same tool on different projects.

I use the [_Tools as Dependencies_](https://mariocarrion.com/2021/10/15/learning-golang-versioning-tools-as-dependencies.html) paradigm to install the tools I use via Go modules. This approach, when combined with the installation of both `direnv` and `make`, facilitates your work by allowing you to execute [`make tools`](../../Makefile#L3) directly to download and install the required packages.

## Go 1.24 `tool`

Go 1.24 added a [`tool` command](https://go.dev/doc/go1.24#go-command) to track tool dependencies, however Dependabot is not able to track updates of those tools using this new Go command that's the reason I still prefer using the old _"tools paradigm"_ instead.

See the [Github issue](https://github.com/dependabot/dependabot-core/issues/12050) for more information.
