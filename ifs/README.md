# IFS

The project contained in this `ifs` directory is an exploration of creating fractals using [iterated function systems](https://en.wikipedia.org/wiki/Iterated_function_system). These systems simply run the same function over and over again, using the output from the prior iteration as input for the next.

## Architecture

The portion of the code that runs the generation of the fractal is separated from the code that renders the fractal to the screen.

All fractal generators can be found in directories below the `generators` directory. They each implment the `engine.Generator` interface.

Each generator is wired up to its renderer in `harness.go`.

## To Do
- [x] Choose fractal to run via command-line argument.
- [x] Configure via YAML config file.
- [] Add more configuration to each type of fractal.