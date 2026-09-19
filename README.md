# kitchen-print-agent-go

A tiny local bridge between a web-based order management panel and a network thermal printer.

## Why this exists

Browsers can't open raw TCP connections to arbitrary devices on your local network — that's a security restriction, not a bug. So a web panel running in Chrome or Edge has no way to talk directly to a receipt printer sitting on the same WiFi.

This agent runs quietly on any computer already on that network and does the one thing browsers can't:

```
Panel (browser)  →  this agent (localhost)  →  thermal printer (LAN)
    fetch()             HTTP :9123               TCP :9100
```

The panel makes a normal HTTP request to `localhost:9123`, which browsers are perfectly happy with. The agent then opens a raw TCP socket to the printer, which only a native program is allowed to do.

## Why Go

- Compiles to a single, dependency-free binary — no Node, no Python, nothing to install on the client's machine.
- Small memory footprint, runs fine on older hardware.
- Cross-compiles cleanly for Windows, macOS, and Linux from the same source.

## Printer protocol: ESC/POS

Most thermal **receipt** printers (not label printers) speak ESC/POS — a set of control codes for formatting output on a serial/network printer (centering text, changing font size, cutting paper). `buildReceipt()` in `main.go` generates these codes.

Label printers (e.g. Zebra) typically use a different protocol (ZPL) and aren't supported here.

## Running locally

```bash
go run main.go
```

Starts the server on `http://localhost:9123`.

## Manual test

```bash
curl -X POST http://localhost:9123/print \
  -H "Content-Type: application/json" \
  -d '{"printerIp":"192.168.1.99","order":{"id":1,"customerName":"Test","items":[{"quantity":2,"title":"Cheeseburger"}]}}'
```

If no real printer is on that IP, you'll get a connection error — that's expected and confirms the code is working correctly.

## Building a shippable binary

```bash
GOOS=windows GOARCH=amd64 go build -o kitchen-print-agent.exe main.go
```

Produces a standalone `.exe` that runs on any Windows machine with zero setup — just copy it over and double-click.

## Code layout

| Part | Responsibility |
|---|---|
| `OrderItem`, `Order`, `PrintRequest` | Data shapes received from the panel |
| `buildReceipt` | Turns an order into ESC/POS-formatted text |
| `sendToPrinter` | Opens a raw TCP connection to the printer and writes the receipt |
| `main` | Starts the HTTP server with two routes: `/health` and `/print` |

## Roadmap

- [ ] Wire up the panel side (trigger a print request when an order is confirmed)
- [ ] Auto-start on boot (Windows Service / launchd)
- [ ] USB printer support, not just network
- [ ] Simple UI for configuring the printer IP instead of hardcoding it
- [ ] Package as a proper installer