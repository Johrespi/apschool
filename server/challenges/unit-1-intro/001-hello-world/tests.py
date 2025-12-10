import sys
from io import StringIO

# Capturar stdout
captured_output = StringIO()
sys.stdout = captured_output

# El codigo del usuario se ejecuta ANTES de este archivo
# Pyodide concatena: user_code + test_code

# Obtener output
output = captured_output.getvalue().strip()

# Tests
assert output == "Hello, World!", f"Se esperaba 'Hello, World!' pero se obtuvo '{output}'"

print("ALL_TESTS_PASSED")
