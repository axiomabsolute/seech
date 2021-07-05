# Seech

## Build & Run

Build + enable full-text search features for Sqlite

`go build --tags "fts5" .`

Similarly, to run:

`go ru n--tags "fts5" . <command>`

## Run test

`go test ./...`

## Example

Using test pokemon data in `data/Pokemon.csv` and XSV for preprocessing:

`xsv slice -e 10 data/pokemon.csv  | xsv select 2 | cat -n | xargs -I {} ./seech trigram index pokedex "data/pokemon.csv" "{}"`

`./seech trigram search pokedex "Saur"`

Or you can split the operations to make the process more performant

`xsv select 2 data/pokemon.csv | cat -n > data/pokemon_numbered.csv`
`./seech trigram batch pokedex "data/pokemon.csv" "data/pokemon_numbered.csv"`
`./seech trigram search pokedex "aur"`