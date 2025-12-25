# USER_OUTPUT contiene el stdout capturado del código del usuario
# Esta variable es inyectada por PyodideService antes de ejecutar este test

output = USER_OUTPUT.strip()

assert output == "Hello, World!", (
    f"Se esperaba 'Hello, World!' pero se obtuvo '{output}'"
)

print("ALL_TESTS_PASSED")
