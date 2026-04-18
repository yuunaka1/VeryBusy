# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-04-18

### Added
- `hex` mode for real-time scrolling hexdump of suspicious PE files.

### Changed
- `soc` layout revised back to 5 panels (Logs, Alerts, Network, Hex, Assets), removing the telemetry graphs to keep the UI clean.

## [1.0.0] - 2026-04-18

### Added
- `verybusy soc` full 5-panel dashboard mode (Left 2, Right 3).
- `logs` stream with randomly generated, authentic-looking system events and varying tick speeds.
- Prepopulated `alerts` screen with a random mix of simulated security events like Brute-force and Ransomware detection.
- `network` mode displaying scrolling matrix-style IP connections.
- Terminal-friendly UI built with Bubble Tea and Lip Gloss.
