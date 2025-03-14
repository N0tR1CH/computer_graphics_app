# Computer graphics app - Bialystok University of technology

Application although multiplatform it was made with mac os in mind ;)

Declarative Canvas using Svelte 4!

It does not use any external library (there comes the stuborness...)

Backend is a go application that communicates via IPC with frontend that is Webview (very electronesque).

To compile the build you need prerequisites listed down below:
- Go compiler/runtime.
- pnpm (can be installed via npm) or npm (can be installed via mise).
- wails binary visit wails.io in order to know more.

Steps:
1. Go to 'frontend' dir and execute command 'pnpm run build' or 'npm run build'
2. Go to main folder and depending on OS do the compilation or run the environment with the hot reload mode (recompiles go code)
    - [Build instructions](https://wails.io/docs/gettingstarted/building)
    - For Windows: wails build -clean -o <name>.exe
    - Hot reload: wails dev


