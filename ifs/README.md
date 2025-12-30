# IFS

The project contained in this `ifs` directory is an exploration of creating fractals using [iterated function systems](https://en.wikipedia.org/wiki/Iterated_function_system). These systems simply run the same function over and over again, using the output from the prior iteration as input for the next.

## Architecture

The portion of the code that runs the generation of the fractal is separated from the code that renders the fractal to the screen.

All fractal generators can be found in directories below the `generators` directory. They each implment the `engine.Generator` interface.

Each generator is wired up to its renderer in `harness.go`.

The system is driven by a supplied config file. See the default one in `config/defaultCfg.yaml`.

## To Do
- [x] Choose fractal to run via command-line argument.
- [x] Configure via YAML config file.
- [x] Add fern configuration for thresholds.
- [ ] Add fern configuration for coefficients.
- [x] Fix tcell renderer (has vertical issue?)
- [ ] Add one more fractal: Sierpinski Carpet
- [ ] Modify 4-vertex sierpinski triangle to be multiple triangles, with random point inside rectangle as shared vertex.
- [ ] Next: Move on to L-Systems.

## Notes

### To use for Sierpinski Carpet
```
8 transforms, with equal probability

1: (x/3,           y/3)
2: (x/3 + 1/3,     y/3)
3: (x/3 + 2/3,     y/3)

4: (x/3,           y/3 + 1/3)
5: (x/3 + 2/3,     y/3 + 1/3)

6: (x/3,           y/3 + 2/3)
7: (x/3 + 1/3,     y/3 + 2/3)
8: (x/3 + 2/3,     y/3 + 2/3)

Note: throw away first 100 iterations or so.

```