A simple experiment to build a convenient runtime IoC container for dependency injection.

Current approach:
- On register, obtain the function's first return value as the dependency type.
- If it is a pointer, get the type by continuously calling reflect.Elem() to dereference it.
- Create a FQN (Fully Qualified Name) of the type, using `packagePath/structOrInterfaceName` as a format.
- Store it with the factory function in the map

- On injection, find the struct field with tag `inject` added onto it.
- Get the type, and find the corresponding type stored in the dependency map.
- Call the factory function, recursively evaluate its dependencies if there is any, and
finally call the function with resolved dependencies.
- Return the result of the function call, type cast it with the provided generic type.

Pros:
- Works for some basic things (not tested.)
- Pretty convenient to use, just register and forget.

Cons:
- Cold start overhead due to reflection (not benchmarked yet)
- Might not work with high amount of dependencies that needs to be build, because this uses recursion.
