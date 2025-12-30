# USER_OUTPUT contiene el stdout capturado del codigo del usuario
# Esta variable es inyectada por PyodideService antes de ejecutar este test

output = USER_OUTPUT.strip()

assert output == "ESPOL", f"Se esperaba 'ESPOL' pero se obtuvo '{output}'"

print("ALL_TESTS_PASSED")
