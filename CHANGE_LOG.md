# Change Log

## [v1.8](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.8) (2019-11-29)

- fixing the code style;
- small refactoring;
- adding usage examples;
- improving:
  - repository decor;
  - CI configuration;
  - some unit tests.

## [v1.7](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.7) (2019-11-08)

- calculate a move quality via searching terminators: move quality is directly proportional to a time of its evaluation;
- in a [transposition table](https://www.chessprogramming.org/Transposition_Table):
  - store moves qualities together with corresponding transpositions;
  - replacing same transpositions on storing only if new one has a greater move quality.

## [v1.6](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.6) (2019-10-29)

- fix setting an inner move searcher and a searching terminator to move searchers;
- refactoring.

## [v1.5](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.5) (2019-10-01)

- optimize move searching via a parallel search ([Lazy SMP](https://www.chessprogramming.org/Iterative_Deepening)):
  - launch concurrent searches with same depths;
- make a [transposition table](https://www.chessprogramming.org/Transposition_Table) safe for concurrent use (via a mutual exclusion lock over a whole storage);
- support searching termination by calling a special method (it's safe for concurrent use);
- refactoring.

## [v1.4](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.4) (2019-09-17)

- optimize move searching via the [iterative deepening](https://www.chessprogramming.org/Iterative_Deepening).

## [v1.3](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.3) (2019-08-24)

- [transposition table](https://www.chessprogramming.org/Transposition_Table):
  - store transpositions in an LRU cache instead of a hash table;
  - share a [transposition table](https://www.chessprogramming.org/Transposition_Table) between searches;
  - disable storing zero moves in a [transposition table](https://www.chessprogramming.org/Transposition_Table);
- refactoring.

## [v1.2](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.2) (2019-07-12)

- optimize move searching via the [transposition table](https://www.chessprogramming.org/Transposition_Table):
  - store transpositions in a hash table;
  - hash a transposition by its representation in [Forsyth–Edwards Notation](https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation);
  - always replace same transpositions on storing;
  - reset a [transposition table](https://www.chessprogramming.org/Transposition_Table) on every search;
- refactoring.

## [v1.1](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.1) (2019-07-03)

- optimize move searching via the [alpha-beta pruning](https://www.chessprogramming.org/Alpha-Beta);
- implement a tournament between the [negamax algorithm](https://www.chessprogramming.org/Negamax) and it with the [alpha-beta pruning](https://www.chessprogramming.org/Alpha-Beta).

## [v1.0](https://github.com/thewizardplusplus/go-chess-minimax/tree/v1.0) (2019-06-29)
