use pyo3::prelude::*;

#[pyfunction]
fn add(a: i64, b: i64) -> PyResult<i64> {
	Ok(a + b)
}

// Python module initializer.
// The module name exposed to Python is controlled by the function name (`rust_lab9_var7`).
#[pymodule]
fn rust_lab9_var7(m: &Bound<'_, PyModule>) -> PyResult<()> {
	m.add_function(wrap_pyfunction!(add, m)?)?;
	Ok(())
}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn add_returns_sum() {
		let out = add(1, 2).unwrap();
		assert_eq!(out, 3);
	}

	#[test]
	fn module_initializer_registers_add() {
		// PyO3 in this crate is built without `auto-initialize`, so we must initialize Python explicitly.
		Python::initialize();

		Python::attach(|py| {
			let m = PyModule::new(py, "rust_lab9_var7").unwrap();
			rust_lab9_var7(&m).unwrap();

			let add_fn = m.getattr("add").unwrap();
			let res: i64 = add_fn.call1((1i64, 2i64)).unwrap().extract().unwrap();
			assert_eq!(res, 3);
		})
	}
}
