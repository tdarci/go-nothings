# Lindenmayer Systems

[L-Systems are nifty things](https://en.wikipedia.org/wiki/L-system) that take an input string, process it character by character, using rules to generate a new string.

Some of the symbols in these strings are drawing instructions, allowing the string to be rendered to the screen.

Initially, I had thought generating and rendering L-Systems would be a separate thing from the Iterated Function Systems project, but it appears that they may be a subset of IFS. Let's see.

## Terms

The starting string is called the AXIOM.

The system contains an ALPHABET, or list of acceptable SYMBOLS.

Each PRODUCTION RULE pertains to a single symbol, and may be qualified by preceding or succeeding symbols.

There is often an ANGLE specified for all rotations.

There is a concept of a STACK that holds position + angle.

## Drawing Directives

- `D` or `🖊️`: Draw line 1 unit in current direction from current position.
- `W` or `🚶`: Walk 1 unit in current direction from current position (do not draw).
- `R` or `↩️`: Rotate clockwise ANGLE degrees.
- `L` or `↪️`: Rotate counterclockwise ANGLE degrees.
- `[` or `⬇️`: Push current position and angle onto the STACK.
- `]` or `⬆️`: Pop current position and angle off the STACK.
