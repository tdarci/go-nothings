# IFS

The project contained in this `ifs` directory is an exploration of creating fractals using [iterated function systems](https://en.wikipedia.org/wiki/Iterated_function_system). These systems simply run the same function over and over again, using the output from the prior iteration as input for the next.

## Running

The code is able to create 3 different fractals, as specified as the first argument when running:
- "sierpinski": Creates a Sierpinski Triangle.
- "carpet": Generates a Sierpinski Carpet.
- "fern": Makes a Barnsley's Fern.

You can also generate this weird thing using "rectangle", which was created just to test out using negative numbers and offsets in the logical coordinate system.

The program will run using default configurations. If you like, you may specify a configuration file as the second argument. Example files `defaultCfg.yaml` and `testCfg.yaml` are included as starting places, and are based on the `Config` struct defined in `config/config.go`.

Here is a sample run command, running directly from Go: `go run ./main -- fern defaultCfg.yaml`

## Architecture

The portion of the code that runs the generation of the fractal is separated from the code that renders the fractal to the screen.

All fractal generators can be found in directories below the `generators` directory. They each implment the `engine.Generator` interface. Renderers are found in directories beneath `renderers`. There are 3 available:
- `EbitRenderer`: The best one. Opens a new window and draws there, using [ebitengine](https://ebitengine.org/).
- `TCellRenderer`: Terminal-based rendering. Super low resolution, but cool. It uses the [tcell library](https://github.com/gdamore/tcell).
- `ListPointRenderer`: Simply spits out a list of the points being generated.

Each generator is wired up to its renderer in `harness.go`.

