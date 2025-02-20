# README

## About

This is the official Wails Svelte-TS template.

## Project Setup

Follow these steps to set up and run the project. Ensure that all commands are executed in the root directory.

## Prerequisites

- [Aria2](https://aria2.github.io/) 
- [Wails](https://wails.io/docs/gettingstarted/installation) 

## Running the Project

### Step 1: Start Aria2

Open **Terminal 1** and run:

```bash
aria2c --enable-rpc --rpc-listen-all
```

If Aria2 is not installed, refer to the [Aria2 Installation Guide](https://aria2.github.io/).

### Step 2: Start Wails

Open **Terminal 2** and run:

```bash
wails dev -tags webkit2_41
```

If Wails is not installed, refer to the [Wails Installation Guide](https://wails.io/docs/gettingstarted/installation).

### Notes

- Keep both terminals running while developing to ensure smooth operation.
- Execute all commands from the **root directory** of the project.



## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
