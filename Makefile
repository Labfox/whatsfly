.PHONY: setup compile sync test build-macos test-setup format lint cleanup bump help

help:
	@echo "Available targets:"
	@echo "  setup       - Install uv (if needed), sync dependencies, and install pre-commit hooks"
	@echo "  compile     - Update uv.lock"
	@echo "  sync        - Sync dependencies with uv.lock"
	@echo "  test        - Run tests using pytest"
	@echo "  build-macos - Build Go backend for macOS (arm64 and amd64)"
	@echo "  test-setup  - Check setuptools metadata and build the package"
	@echo "  format      - Format code with black"
	@echo "  lint        - Lint code with ruff and mypy"
	@echo "  cleanup     - Remove build artifacts and virtual environment"
	@echo "  bump        - Update Go dependencies and sync Python environment"

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

test-setup:
	rm -rf build dist *.egg-info
	uv run python -m setuptools.dist check
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
