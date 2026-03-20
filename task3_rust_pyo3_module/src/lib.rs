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
}
