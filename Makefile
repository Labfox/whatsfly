.PHONY: setup compile sync test build-macos build-linux build-windows test-setup format lint cleanup bump help

help:
	@echo "Available targets:"
	@echo "  setup         - Install uv (if needed), sync dependencies, and install pre-commit hooks"
	@echo "  compile       - Update uv.lock"
	@echo "  sync          - Sync dependencies with uv.lock"
	@echo "  test          - Run tests using pytest"
	@echo "  build-macos   - Build Go backend for macOS (arm64 and amd64)"
	@echo "  build-linux   - Build Go backend for Linux (amd64)"
	@echo "  build-windows - Build Go backend for Windows (amd64)"
	@echo "  test-setup    - Build the backend for current platform and build the package"
	@echo "  format        - Format code with black"
	@echo "  lint          - Lint code with ruff and mypy"
	@echo "  cleanup       - Remove build artifacts and virtual environment"
	@echo "  bump          - Update Go dependencies and sync Python environment"

setup:
	@if ! command -v uv >/dev/null 2>&1; then pip install uv; fi
	uv sync --all-extras
	uv run pre-commit install

compile:
	uv lock

sync:
	uv sync --all-extras

test:
	uv run pytest $(ARGS)

build-macos:
	$(MAKE) -C backend build-macos

build-linux:
	$(MAKE) -C backend build-linux

build-windows:
	$(MAKE) -C backend build-windows

test-setup:
	rm -rf build dist *.egg-info
	# Build the backend for the current platform and place it in the expected location
	cd backend && go build -buildmode=c-shared -o ../whatsfly/dependencies/latest.$$(python3 -c "import platform; print({'linux': 'so', 'windows': 'dll', 'darwin': 'dylib'}.get(platform.system().lower(), 'so'))") .
	uv build && echo ok

format:
	uv run black .

lint:
	uv run ruff check .
	uv run ruff format .
	uv run mypy .

cleanup:
	rm -rf build dist *.egg-info whatsfly/last_binary_update.txt whatsfly/__pycache__ whatsfly/*/__pycache__ .venv

bump:
	$(MAKE) -C backend bump
	$(MAKE) compile
	$(MAKE) sync
	$(MAKE) setup
	$(MAKE) lint
	$(MAKE) format
	$(MAKE) test-setup
	$(MAKE) cleanup
