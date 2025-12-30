# USER_OUTPUT contiene el stdout capturado del codigo del usuario
# Esta variable es inyectada por PyodideService antes de ejecutar este test

output = USER_OUTPUT.strip()

assert output == "15", f"Se esperaba '15' pero se obtuvo '{output}'"

print("ALL_TESTS_PASSED")
