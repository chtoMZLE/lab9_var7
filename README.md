# Lab9: Go/Rust/PyO3 + PyPI CI/CD

## Названия задач (Task 1–5)
1. `task1_go_async_processor/` — Go async request processor (goroutine worker)
2. `task2_go_tcp_server/` — Go TCP server + Python client
3. `task3_rust_pyo3_module/` — PyO3 Python module in Rust (`#[pymodule]` + `add`)
4. `task4_ci_cd_pypi/` — CI/CD build & publish PyO3 module to PyPI
5. `task5_go_compute_service/` — Go heavy-compute microservice invoked from Python (primes)

## Структура репозитория
- `work_plan` — исходный план лабораторной
- `go.mod` — модуль Go (пакеты лежат в подпапках `task*`)
- `.gitignore` — игнор сборочных артефактов (`target/`, `__pycache__/`, и т.п.)
- `task1_go_async_processor/` — async worker
- `task2_go_tcp_server/` — TCP-сервер и клиент
- `task3_rust_pyo3_module/` — Rust crate для PyO3 + `pyproject.toml` для `maturin`
- `task4_ci_cd_pypi/` — документация по CI/CD (а workflow находится в `.github/workflows/`)
- `task5_go_compute_service/` — Go HTTP сервис и Python клиент
- `.github/workflows/task4_publish_pypi.yml` — workflow публикации на PyPI

## Требования
1. Go установлен (для `task1`, `task2`, `task5`)
2. Rust установлен (для `task3` и `task4`)
3. Python 3.x установлен (для `pyo3`/`maturin` и `task5` клиента)

## Локальная проверка
- Unit-тесты Go (все пакеты): `go test ./...`
- Unit-тесты Python (тест клиента primes): `py -m unittest -q task5_go_compute_service.tests.test_client`
- Unit-тесты PyO3: `cargo test` внутри `task3_rust_pyo3_module/`

## CI/CD (Task 4) — публикация на PyPI
Workflow триггерится по тегам `v*.*.*` (например `v0.1.0`).
В GitHub Actions нужен секрет `PYPI_API_TOKEN`.

Настройки `maturin` лежат в:
- `task3_rust_pyo3_module/pyproject.toml`

Публикация выполняется:
- `.github/workflows/task4_publish_pypi.yml` (через `maturin publish`)

