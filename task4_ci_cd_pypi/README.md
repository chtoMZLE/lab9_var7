# Task 4: CI/CD build & publish PyO3 module to PyPI

This task adds:

- `task3_rust_pyo3_module/pyproject.toml` configured for `maturin`
- GitHub Actions workflow to build wheels for multiple OSes and publish to PyPI on version tags

The workflow expects a secret named `PYPI_API_TOKEN`.

