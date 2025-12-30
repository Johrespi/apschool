# USER_OUTPUT contiene el stdout capturado del codigo del usuario
# Esta variable es inyectada por PyodideService antes de ejecutar este test

output = USER_OUTPUT.strip()

assert output == "6", f"Se esperaba '6' pero se obtuvo '{output}'"

print("ALL_TESTS_PASSED")
