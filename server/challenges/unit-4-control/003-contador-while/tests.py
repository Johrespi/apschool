# USER_OUTPUT contiene el stdout capturado del codigo del usuario
# Esta variable es inyectada por PyodideService antes de ejecutar este test

output = USER_OUTPUT.strip()
expected = "1\n2\n3\n4\n5"

assert output == expected, f"Se esperaba '1\\n2\\n3\\n4\\n5' pero se obtuvo '{output}'"

print("ALL_TESTS_PASSED")
